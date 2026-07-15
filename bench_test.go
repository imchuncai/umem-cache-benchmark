// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"math"
	"math/bits"
	"math/rand/v2"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/twmb/murmur3"
)

const (
	TIMEOUT = 30 * time.Second
	SEED    = 47
	HEX     = "0123456789abcdef"
)

type Client interface {
	Init(remoteIPV6 string, threadNR int, config *tls.Config) error
	GetOrSet(key []byte, fallbackVal func() []byte) ([]byte, error)
}

type TestCase struct {
	HexIndex [16]byte
	KeySize  uint32
	KvSize   uint32
}

func __parallel(b *testing.B, client Client, zipf *rand.Zipf, kvSizeMax uint32, randSize bool) (
	throughput uint64, output uint64, miss uint64) {
	// key size is rand from 16~47 bytes
	if kvSizeMax < 47 {
		panic("bad kvSizeMax")
	}

	cases := make([]TestCase, b.N)
	for i := 0; i < b.N; i++ {
		index := zipf.Uint64()
		var hex [16]byte
		for i := 0; i < 16; i++ {
			hex[i] = HEX[(index>>(i<<2))&0xf]
		}
		h1, h2 := murmur3.SeedSum128(SEED, SEED, hex[:])

		keySize := uint32(h1)&31 + 16
		if randSize {
			hi, _ := bits.Mul32(uint32(h2), kvSizeMax-keySize+1)
			cases[i] = TestCase{hex, keySize, hi + keySize}
		} else {
			cases[i] = TestCase{hex, keySize, kvSizeMax}
		}
		throughput += uint64(cases[i].KvSize)
	}

	// make the benchmark result as stable as possible
	var atomicI atomic.Uint64
	atomicI.Store(math.MaxUint64)

	keyTemplate := make([]byte, 47)
	for i := 0; i < len(keyTemplate); i += len(HEX) {
		copy(keyTemplate[i:], HEX)
	}

	valTemplate := make([]byte, kvSizeMax-16)
	for i := 0; i < len(valTemplate); i += len(HEX) {
		copy(valTemplate[i:], HEX)
	}

	b.StartTimer()
	b.ResetTimer()
	defer b.StopTimer()
	b.RunParallel(func(p *testing.PB) {
		keyTemp := bytes.Clone(keyTemplate)
		var __output, __miss uint64
		for p.Next() {
			tc := cases[atomicI.Add(1)]
			copy(keyTemp, tc.HexIndex[:])

			fallbackVal := func() []byte {
				__miss++
				__output += uint64(tc.KvSize)
				return valTemplate[:tc.KvSize-tc.KeySize]
			}

			_, err := client.GetOrSet(keyTemp[:tc.KeySize], fallbackVal)
			if err != nil {
				b.Fatalf("got error: %v", err)
			}
		}
		atomic.AddUint64(&output, __output)
		atomic.AddUint64(&miss, __miss)
	})
	return
}

func parallel[T any, PT interface {
	*T
	Client
}](b *testing.B) {
	if b.N == 1 {
		// benchmark is called twice, drop the first
		return
	}

	args := flag.Args()
	if len(args) < 7 {
		b.Fatal("bad args length")
	}
	randSize := args[0] == "true"
	serverMemory, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		b.Fatalf("bad arg serverMemory: %s", args[1])
	}
	__kvSizeMax, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil || __kvSizeMax > math.MaxUint32 {
		b.Fatalf("bad arg kvSizeMax")
	}
	kvSizeMax := uint32(__kvSizeMax)
	parallelism, err := strconv.Atoi(args[3])
	if err != nil {
		b.Fatalf("bad arg parallelism")
	}
	b.SetParallelism(parallelism)
	threadNR, err := strconv.Atoi(args[4])
	if err != nil {
		b.Fatalf("bad arg threadNR")
	}
	tlsEnable, err := strconv.Atoi(args[5])
	if err != nil {
		b.Fatalf("bad arg tls")
	}
	remoteIP := args[6]

	var config *tls.Config
	if tlsEnable == 1 {
		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		if err != nil {
			b.Fatalf("load tls key pair failed: %v", err)
		}
		caCert, err := os.ReadFile("ca-cert.pem")
		if err != nil {
			b.Fatalf("read tls file: ca-cert.pem failed: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			b.Fatal("tls append certs failed")
		}
		config = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
			ServerName:   "nas.local",
		}
	}

	client := PT(new(T))
	err = client.Init(remoteIP, threadNR, config)
	if err != nil {
		b.Fatalf("client init failed: %v", err)
	}

	cap := serverMemory / uint64(kvSizeMax)
	if randSize {
		cap *= 2
	}
	r := rand.New(rand.NewPCG(SEED, SEED))
	zipf := rand.NewZipf(r, 1.0001, 1.0, cap*1000)

	// warmup
	__parallel(b, client, zipf, kvSizeMax, randSize)

	throughput, output, miss := __parallel(b, client, zipf, kvSizeMax, randSize)
	hit := b.N - int(miss)
	hitRate := float64(hit) / float64(b.N) * 100
	fmt.Printf("\n======================================================================\n"+
		"server: %8d    warmup: %8d    get: %8d    hit: %8d\n"+
		"VmHWM: %7d kB   hit_rate: %.2f%%    per_memory_hit_rate: %.2f%%\n"+
		"%.3fs\t    output: %4.0f Mb/s   input: %4.0f Mb/s\n"+
		"======================================================================\n",
		cap, b.N, b.N, hit, 0, hitRate, hitRate, b.Elapsed().Seconds(),
		float64(output*8)/1024/1024/b.Elapsed().Seconds(),
		float64((throughput-output)*8)/1024/1024/b.Elapsed().Seconds(),
	)
	b.ReportMetric(float64(hit)/b.Elapsed().Seconds(), "hit/s/mem")
}

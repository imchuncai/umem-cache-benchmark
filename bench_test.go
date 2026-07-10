// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"hash/fnv"
	"math"
	"math/bits"
	"math/rand/v2"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

const (
	TIMEOUT = 30 * time.Second
	SEED    = 47
)

type Client interface {
	Init(remoteIPV6 string, threadNR int, config *tls.Config) error
	GetOrSet(key []byte, i uint64, fallbackVal func() []byte) ([]byte, error)
}

func __parallel(b *testing.B, client Client, zipf *rand.Zipf, kvSizeLimit uint32, randSize bool) (
	throughput uint64, output uint64, miss uint64) {
	// key size is rand from 16~47 bytes
	if kvSizeLimit < 47 {
		panic("bad kvSizeLimit")
	}

	// make the benchmark result as stable as possible
	var atomicI atomic.Uint64
	atomicI.Store(math.MaxUint64)
	indexes := make([]uint64, b.N)
	for i := 0; i < b.N; i++ {
		indexes[i] = zipf.Uint64()
	}

	const hex = "0123456789abcdef"
	kvTemplate := make([]byte, kvSizeLimit)
	for i := 0; i < len(kvTemplate); i += len(hex) {
		copy(kvTemplate[i:], hex)
	}

	KvSize := func(h uint64, keyLen uint32) uint32 {
		return kvSizeLimit
	}
	if randSize {
		KvSize = func(h uint64, keyLen uint32) uint32 {
			hi, _ := bits.Mul32(uint32(h>>32), kvSizeLimit-keyLen)
			return hi + keyLen
		}
	}

	b.StartTimer()
	b.ResetTimer()
	defer b.StopTimer()
	b.RunParallel(func(p *testing.PB) {
		s := bytes.Clone(kvTemplate)
		r := fnv.New64a()
		var __throughput, __output, __miss uint64
		for p.Next() {
			index := indexes[atomicI.Add(1)]
			for i := 0; i < 16; i++ {
				s[i] = hex[(index>>(i<<2))&0xf]
			}
			r.Reset()
			r.Write(s[:16])
			h := r.Sum64()

			keyLen := uint32(h)&31 + 16
			fallbackVal := func() []byte {
				__miss++
				kvSize := KvSize(h, keyLen)
				__output += uint64(kvSize)
				return s[keyLen:kvSize]
			}
			val, err := client.GetOrSet(s[:keyLen], index, fallbackVal)
			if err != nil {
				b.Fatalf("got error: %v", err)
			}
			__throughput += uint64(keyLen) + uint64(len(val))
		}
		atomic.AddUint64(&throughput, __throughput)
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
	__kvSizeLimit, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil || __kvSizeLimit > math.MaxUint32 {
		b.Fatalf("bad arg kvSizeLimit")
	}
	kvSizeLimit := uint32(__kvSizeLimit)
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

	cap := serverMemory / uint64(kvSizeLimit)
	if randSize {
		cap *= 2
	}
	r := rand.New(rand.NewPCG(SEED, SEED))
	zipf := rand.NewZipf(r, 1.0001, 1.0, cap*1000)

	// warmup
	__parallel(b, client, zipf, kvSizeLimit, randSize)

	throughput, output, miss := __parallel(b, client, zipf, kvSizeLimit, randSize)
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

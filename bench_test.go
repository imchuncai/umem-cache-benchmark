// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/imchuncai/umem-cache-benchmark/test"
	"github.com/prometheus/procfs"
)

const (
	THREAD_NR = 4
	TIMEOUT   = 10 * time.Second
)

func percent(i, n int) float64 {
	return float64(i) / float64(n) * 100
}

type GetOrSetFunc func(key []byte, i int, f func() []byte) ([]byte, error)

type RunServer func(serverMemory int, kvSizeLimit int, remoteIP string) ([]*exec.Cmd, GetOrSetFunc, error)

func __parallel(b *testing.B, pool *test.Pool, getOrSet GetOrSetFunc) (hot, hotMiss, miss uint64) {
	b.RunParallel(func(p *testing.PB) {
		var __hot, __hotMiss, __miss uint64

		for p.Next() {
			tc := pool.RandCase()
			if tc.Hot {
				__hot++
			}
			fallbackGet := func() []byte {
				if tc.Hot {
					__hotMiss++
				}
				__miss++
				return tc.Val
			}
			_, err := getOrSet(tc.Key, tc.I, fallbackGet)
			if err != nil {
				b.Fatalf("got error: %v", err)
			}
		}

		atomic.AddUint64(&hot, __hot)
		atomic.AddUint64(&hotMiss, __hotMiss)
		atomic.AddUint64(&miss, __miss)
	})
	return
}

func getVmHWM(b testing.TB, cmds []*exec.Cmd) uint64 {
	vmHWM := uint64(0)
	for _, cmd := range cmds {
		p, err := procfs.NewProc(cmd.Process.Pid)
		if err != nil {
			b.Fatalf("get process failed: %v", err)
		}

		status, err := p.NewStatus()
		if err != nil {
			b.Fatalf("get process status failed: %v", err)
		}

		vmHWM += status.VmHWM
	}
	return vmHWM
}

func stop(cmd *exec.Cmd) error {
	err := cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("signal SIGTERM failed: %w", err)
	}
	_, err = cmd.Process.Wait()
	if err != nil {
		return fmt.Errorf("wait failed: %w", err)
	}
	return nil
}

func parallel(b *testing.B, run RunServer) {
	b.StopTimer()

	if b.N == 1 {
		// benchmark is called twice, drop the first
		return
	}

	args := flag.Args()
	if len(args) < 4 {
		b.Fatal("bad args")
	}
	randSize := args[0] == "true"
	serverMemory, err := strconv.Atoi(args[1])
	if err != nil {
		b.Fatalf("bad arg serverMemory: %s", args[1])
	}
	kvSizeLimit, err := strconv.Atoi(args[2])
	if err != nil {
		b.Fatalf("bad arg kvSizeLimit")
	}
	hotCasePercent, err := strconv.Atoi(args[3])
	if err != nil {
		b.Fatalf("bad arg hotCasePercent")
	}
	hotCaseAccessPercent, err := strconv.Atoi(args[4])
	if err != nil {
		b.Fatalf("bad arg hotCasePercent")
	}
	hotCaseServerMemoryPercent, err := strconv.Atoi(args[5])
	if err != nil {
		b.Fatalf("bad arg hotCaseServerMemoryPercent")
	}
	parallelism, err := strconv.Atoi(args[6])
	if err != nil {
		b.Fatalf("bad arg parallelism")
	}
	b.SetParallelism(parallelism)

	remoteIP := ""
	if len(args) > 7 {
		remoteIP = args[7]
	}

	cmds, getOrSet, err := run(serverMemory, kvSizeLimit, remoteIP)
	if err != nil {
		b.Fatalf("run server failed: %v", err)
	}
	defer func() {
		var err error
		for i, cmd := range cmds {
			e := stop(cmd)
			if e != nil {
				err = errors.Join(err, fmt.Errorf("stop machine: %d failed: %w", i, e))
			}
		}
		if err != nil {
			b.Fatalf("stop failed: %v", err)
		}
	}()
	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			b.Fatalf("start server failed: %v", err)
		}
	}
	time.Sleep(100 * time.Millisecond)
	for _, cmd := range cmds {
		err = cmd.Process.Signal(syscall.Signal(0))
		if err != nil {
			b.Fatalf("server exited early: %v", err)
		}
	}

	pool := test.NewPool(kvSizeLimit, randSize, serverMemory, hotCasePercent, hotCaseAccessPercent, hotCaseServerMemoryPercent)
	// warmup
	__parallel(b, pool, getOrSet)

	b.StartTimer()
	hot, hotMiss, miss := __parallel(b, pool, getOrSet)
	b.StopTimer()

	hit := b.N - int(miss)
	hotHit := hot - hotMiss
	vmHWM := getVmHWM(b, cmds)
	hitRate := percent(hit, b.N)
	fmt.Printf("\n   =======================================================\n"+
		"    case:%8d"+"    hot:%8d(%d%%)"+"    hot_access: %d%% \n"+
		"     get:%8d"+"    hit:%8d"+"    hit_rate: %.2f%% \n"+
		"     hot:%8d"+"    hit:%8d"+"    hit_rate: %.2f%% \n"+
		"     VmHWM: %7d kB    per_memory_hit_rate: %.2f%%\n"+
		"     %.3fs\n"+
		"   =======================================================\n",
		pool.CaseN(), pool.HotN(), hotCasePercent, hotCaseAccessPercent,
		b.N, hit, percent(hit, b.N),
		hot, hotHit, percent(int(hotHit), int(hot)),
		vmHWM>>10, hitRate/float64(vmHWM)*float64(serverMemory),
		b.Elapsed().Seconds(),
	)
	hitF := float64(hit) / b.Elapsed().Seconds()
	b.ReportMetric(hitF/float64(vmHWM)*float64(serverMemory), "hit/s/mem")
}

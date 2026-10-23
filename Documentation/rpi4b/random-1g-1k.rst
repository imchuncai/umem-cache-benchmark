.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-1g-1k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 10% higher than Memcached and 14% higher than Redis.
	Umem-cache's hit throughput is 54% higher than Memcached and 69% higher than Redis.

	Umem-cache's P90   latency is 55% lower than Memcached and 50% lower than Redis.
	Umem-cache's P99.9 latency is 48% lower than Memcached and 55% lower than Redis.

Memcached
=========
::

	commit e44dd0b01234bc0faf970e9225e3423e98022129

Build Command
-------------
::

	./autogen.sh
	./configure --enable-tls
	make -j

Run Command
-----------
::

	taskset -c 1 ./memcached --memory-limit=1024 \
	--max-item-size=1048576 -t 1

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20513448
	VmHWM: 1076832 kB   hit_rate: 61.13%    per_memory_hit_rate: 59.53%
	P1: 272 us  P50: 1021 us  P90: 2028 us  P99: 2499 us  P99.9: 4055 us
	1522.509s	    output:   35 Mb/s   input:   56 Mb/s
	======================================================================
	33554432	     45374 ns/op	     13120 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3100.684s

Umem-cache
==========
::

	commit 079c92dc41d66b84753ccc3d5ae78974460009c0

Build Command
-------------
::

	make MEM_LIMIT=1073741824 THREAD_NR=1 MAX_CONN=24

Run Command
-----------
::

	taskset -c 1 ./umem-cache 10047

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22017474
	VmHWM: 1049776 kB   hit_rate: 65.62%    per_memory_hit_rate: 65.54%
	P1: 251 us  P50: 758 us  P90: 911 us  P99: 1532 us  P99.9: 2127 us
	1087.104s	    output:   43 Mb/s   input:   84 Mb/s
	======================================================================
	33554432	     32398 ns/op	     20230 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2223.153s

Redis
=====
::

	commit 6bf6224c3dad518329ddc893ef9c5d58dcbabdeb

Build Command
-------------
::

	make -j BUILD_TLS=yes

Run Command
-----------
::

	taskset -c 1 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 24 --maxmemory-policy allkeys-lfu \
	--port 6379

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=33554432x \
	-args true 1073741824 1024 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20026192
	VmHWM: 1086536 kB   hit_rate: 59.68%    per_memory_hit_rate: 57.60%
	P1: 312 us  P50: 1016 us  P90: 1836 us  P99: 2831 us  P99.9: 4730 us
	1611.751s	    output:   34 Mb/s   input:   51 Mb/s
	======================================================================
	33554432	     48034 ns/op	     11991 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3261.912s

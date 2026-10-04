.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-1g-1k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 10% higher than Memcached and 14% higher than Redis.
	Umem-cache's hit throughput is 55% higher than Memcached and 71% higher than Redis.

	Umem-cache's P90   latency is 56% lower than Memcached and 52% lower than Redis.
	Umem-cache's P99.9 latency is 42% lower than Memcached and 54% lower than Redis.

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
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20520871
	VmHWM: 1076828 kB   hit_rate: 61.16%    per_memory_hit_rate: 59.55%
	P1: 289 us  P50: 1046 us  P90: 2082 us  P99: 2525 us  P99.9: 4201 us
	1546.794s	    output:   34 Mb/s   input:   55 Mb/s
	======================================================================
	33554432	     46098 ns/op	     12919 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3143.283s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

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
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22017471
	VmHWM: 1049696 kB   hit_rate: 65.62%    per_memory_hit_rate: 65.55%
	P1: 252 us  P50: 763 us  P90: 920 us  P99: 1592 us  P99.9: 2441 us
	1096.946s	    output:   42 Mb/s   input:   83 Mb/s
	======================================================================
	33554432	     32692 ns/op	     20050 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2246.070s

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
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20017106
	VmHWM: 1084568 kB   hit_rate: 59.66%    per_memory_hit_rate: 57.68%
	P1: 282 us  P50: 1038 us  P90: 1912 us  P99: 3482 us  P99.9: 5342 us
	1650.950s	    output:   33 Mb/s   input:   50 Mb/s
	======================================================================
	33554432	     49202 ns/op	     11722 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3340.746s

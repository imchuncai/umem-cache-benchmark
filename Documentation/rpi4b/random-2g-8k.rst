.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-2g-8k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 9% higher than Memcached and 12% higher than Redis.
	Umem-cache's hit throughput is 13% higher than Memcached and 29% higher than Redis.

	Umem-cache's P90   latency is 35% lower than Memcached and 35% lower than Redis.
	Umem-cache's P99.9 latency is 41% lower than Memcached and 48% lower than Redis.

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

	taskset -c 1 ./memcached --memory-limit=2048 \
	--max-item-size=1048576 -t 1

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4972019
	VmHWM: 2118756 kB   hit_rate: 59.27%    per_memory_hit_rate: 58.67%
	P1: 418 us  P50: 1402 us  P90: 2689 us  P99: 3310 us  P99.9: 4641 us
	539.081s	    output:  200 Mb/s   input:  298 Mb/s
	======================================================================
	 8388608	     64263 ns/op	      9129 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1099.359s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

Build Command
-------------
::

	make MEM_LIMIT=2147483648 THREAD_NR=1 MAX_CONN=24

Run Command
-----------
::

	taskset -c 1 ./umem-cache 10047

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5356170
	VmHWM: 2098292 kB   hit_rate: 63.85%    per_memory_hit_rate: 63.82%
	P1: 632 us  P50: 1475 us  P90: 1744 us  P99: 2033 us  P99.9: 2720 us
	518.814s	    output:  183 Mb/s   input:  334 Mb/s
	======================================================================
	 8388608	     61847 ns/op	     10318 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1036.874s

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
	--maxmemory 2147483648 --maxclients 24 --maxmemory-policy allkeys-lfu \
	--port 6379

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=8388608x \
	-args true 2147483648 8192 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4923126
	VmHWM: 2154040 kB   hit_rate: 58.69%    per_memory_hit_rate: 57.14%
	P1: 486 us  P50: 1547 us  P90: 2665 us  P99: 3550 us  P99.9: 5259 us
	598.650s	    output:  181 Mb/s   input:  267 Mb/s
	======================================================================
	 8388608	     71365 ns/op	      8007 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1210.304s

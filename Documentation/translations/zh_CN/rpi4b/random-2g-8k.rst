.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-2g-8k
====================

结论
====
::

	Umem-cache的命中率比Memcached高9%，比Redis高12%。
	Umem-cache的命中吞吐量比Memcached高13%，比Redis高32%。

	Umem-cache的P90  延迟比Memcached低36%，比Redis低36%。
	Umem-cache的P99.9延迟比Memcached低21%，比Redis低27%。

Memcached
=========
::

	commit e44dd0b01234bc0faf970e9225e3423e98022129

编译命令
-------
::

	./autogen.sh
	./configure --enable-tls
	make -j

运行命令
-------
::

	taskset -c 1 ./memcached --memory-limit=2048 \
	--max-item-size=1048576 -t 1

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4971643
	VmHWM: 2118828 kB   hit_rate: 59.27%    per_memory_hit_rate: 58.66%
	P1: 421 us  P50: 1405 us  P90: 2704 us  P99: 3287 us  P99.9: 4495 us
	536.943s	    output:  200 Mb/s   input:  299 Mb/s
	======================================================================
	 8388608	     64009 ns/op	      9164 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1095.178s

Umem-cache
==========
::

	commit 079c92dc41d66b84753ccc3d5ae78974460009c0

编译命令
-------
::

	make MEM_LIMIT=2147483648 THREAD_NR=1 MAX_CONN=24

运行命令
-------
::

	taskset -c 1 ./umem-cache 10047

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5356168
	VmHWM: 2098340 kB   hit_rate: 63.85%    per_memory_hit_rate: 63.81%
	P1: 649 us  P50: 1467 us  P90: 1734 us  P99: 2026 us  P99.9: 3535 us
	517.242s	    output:  184 Mb/s   input:  335 Mb/s
	======================================================================
	 8388608	     61660 ns/op	     10349 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1031.488s

Redis
=====
::

	commit 6bf6224c3dad518329ddc893ef9c5d58dcbabdeb

编译命令
-------
::

	make -j BUILD_TLS=yes

运行命令
-------
::

	taskset -c 1 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 2147483648 --maxclients 24 --maxmemory-policy allkeys-lfu \
	--port 6379

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=8388608x \
	-args true 2147483648 8192 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4920285
	VmHWM: 2152320 kB   hit_rate: 58.65%    per_memory_hit_rate: 57.15%
	P1: 496 us  P50: 1578 us  P90: 2690 us  P99: 3511 us  P99.9: 4872 us
	611.086s	    output:  178 Mb/s   input:  261 Mb/s
	======================================================================
	 8388608	     72847 ns/op	      7845 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1238.917s

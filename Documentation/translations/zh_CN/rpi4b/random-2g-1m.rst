.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-2g-1m
====================

结论
====
::

	Umem-cache的命中率比Memcached高15%，比Redis高13%。
	Umem-cache的命中吞吐量比Memcached高17%，比Redis高11%。

	Umem-cache的P90  延迟比Memcached低0%，比Redis低0%。
	Umem-cache的P99.9延迟比Memcached低0%，比Redis低0%。

	注意：网络吞吐量已经接近千兆网络的极限。

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
	--max-item-size=2097152 -t 1

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30666
	VmHWM: 2120968 kB   hit_rate: 46.79%    per_memory_hit_rate: 46.27%
	P1: 1227 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	175.530s	    output:  800 Mb/s   input:  736 Mb/s
	======================================================================
	   65536	   2678370 ns/op	       173 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	350.894s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34855
	VmHWM: 2098300 kB   hit_rate: 53.18%    per_memory_hit_rate: 53.16%
	P1: 902 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	172.410s	    output:  714 Mb/s   input:  849 Mb/s
	======================================================================
	   65536	   2630768 ns/op	       202 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	344.534s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=65536x \
	-args true 2147483648 1048576 8 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31490
	VmHWM: 2139892 kB   hit_rate: 48.05%    per_memory_hit_rate: 47.09%
	P1: 1490 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	169.830s	    output:  805 Mb/s   input:  782 Mb/s
	======================================================================
	   65536	   2591395 ns/op	       182 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	340.584s

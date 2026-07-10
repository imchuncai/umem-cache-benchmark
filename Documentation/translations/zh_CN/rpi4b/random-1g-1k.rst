.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-1g-1k
====================

结论
====
::

	Umem-cache的命中率比Memcached高10%，比Redis高16%。
	Umem-cache的命中吞吐量比Memcached高15%，比Redis高61%。

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

	taskset -c 1,2,3 ./memcached --memory-limit=1024 \
	--max-item-size=1048576 -t 3

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20522752
	VmHWM: 1077388 kB   hit_rate: 61.16%    per_memory_hit_rate: 59.53%
	935.141s	    output:   56 Mb/s   input:   93 Mb/s
	======================================================================
	33554432	     27869 ns/op	     21359 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1902.440s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

编译命令
-------
::

	make MEM_LIMIT=1073741824 THREAD_NR=3 MAX_CONN=144

运行命令
-------
::

	taskset -c 1,2,3 ./umem-cache 10047

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22008953
	VmHWM: 1049792 kB   hit_rate: 65.59%    per_memory_hit_rate: 65.52%
	893.215s	    output:   52 Mb/s   input:  104 Mb/s
	======================================================================
	33554432	     26620 ns/op	     24612 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1836.055s

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

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6380

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6381

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=33554432x \
	-args true 1073741824 1024 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20120083
	VmHWM:  372452 kB   hit_rate: 59.96%    per_memory_hit_rate: 56.28%
	VmHWM:  372344 kB
	VmHWM:  372360 kB
	1233.329s	    output:   44 Mb/s   input:   69 Mb/s
	======================================================================
	33554432	     36756 ns/op	     15312 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2494.176s

.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-1g-1k
====================

结论
====
::

	Umem-cache的命中率比Memcached高9%，比Redis高14%。
	Umem-cache的命中吞吐量比Memcached高6%，比Redis高38%。

Memcached
=========
::

	commit f1674f0231e5d291db474c4ad297f5f069d32521

编译命令
-------
::

	./autogen.sh
	./configure --enable-tls
	make -j

运行命令
-------
::

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=1024 \
	--max-item-size=1048576 -t 3 -u root

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20520387
	VmHWM: 1070756 kB   hit_rate: 61.16%    per_memory_hit_rate: 59.89%
	809.483s	    output:   65 Mb/s   input:  107 Mb/s
	======================================================================
	33554432	     24124 ns/op	     24825 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1649.656s

Umem-cache
==========
::

	commit eeeac62ec7a2ec135b4f2e419a533ba8b3282ccc

编译命令
-------
::

	make MEM_LIMIT=1073741824 THREAD_NR=3

运行命令
-------
::

	taskset -c 1,2,3 ./umem-cache 10047

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22008683
	VmHWM: 1050204 kB   hit_rate: 65.59%    per_memory_hit_rate: 65.49%
	833.143s	    output:   56 Mb/s   input:  111 Mb/s
	======================================================================
	33554432	     24830 ns/op	     26376 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1711.379s

Redis
=====
::

	commit 138263a1b480fcd2e756be27f369203a46481d06

编译命令
-------
::

	make -j BUILD_TLS=yes

运行命令
-------
::

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6380

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20202812
	VmHWM:  547596 kB   hit_rate: 60.21%    per_memory_hit_rate: 57.67%
	VmHWM:  547212 kB
	1010.400s	    output:   53 Mb/s   input:   85 Mb/s
	======================================================================
	33554432	     30112 ns/op	     19151 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2043.197s

.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-2g-8k
====================

结论
====
::

	Umem-cache的命中率比Memcached高8%，比Redis高12%。
	Umem-cache的命中吞吐量比Memcached高24%，比Redis高20%。

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

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=1048576 -t 3 -u root

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4971757
	VmHWM: 2109752 kB   hit_rate: 59.27%    per_memory_hit_rate: 58.91%
	448.235s	    output:  240 Mb/s   input:  367 Mb/s
	======================================================================
	 8388608	     53434 ns/op	     11026 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	907.248s

Umem-cache
==========
::

	commit e4a931ea3ee8fc82b3693333fa12a264e36f3dd0

编译命令
-------
::

	make MEM_LIMIT=2147483648 THREAD_NR=3

运行命令
-------
::

	taskset -c 1,2,3 ./umem-cache 10047

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5364257
	VmHWM: 2098796 kB   hit_rate: 63.95%    per_memory_hit_rate: 63.90%
	392.144s	    output:  242 Mb/s   input:  453 Mb/s
	======================================================================
	 8388608	     46747 ns/op	     13669 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	797.183s

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
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6380

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4957479
	VmHWM: 1085036 kB   hit_rate: 59.10%    per_memory_hit_rate: 57.08%
	VmHWM: 1086236 kB
	421.843s	    output:  255 Mb/s   input:  391 Mb/s
	======================================================================
	 8388608	     50288 ns/op	     11351 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	857.427s

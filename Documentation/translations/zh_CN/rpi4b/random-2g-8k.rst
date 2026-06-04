.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-2g-8k
====================

结论
====
::

	Umem-cache的命中率比Memcached高8%，比Redis高12%。
	Umem-cache的命中吞吐量比Memcached高23%，比Redis高32%。

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

	./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=1048576 -t 4 -u root

测试结果
-------
::

	go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-4   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4971498
	VmHWM: 2110052 kB   hit_rate: 59.26%    per_memory_hit_rate: 58.90%
	462.127s	    output:  233 Mb/s   input:  356 Mb/s
	======================================================================
	 8388608	     55090 ns/op	     10692 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	936.017s

Umem-cache
==========
::

	commit eeeac62ec7a2ec135b4f2e419a533ba8b3282ccc

编译命令
-------
::

	make MEM_LIMIT=2147483648

运行命令
-------
::

	./umem-cache 10047

测试结果
-------
::

	go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-4   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5355767
	VmHWM: 2098884 kB   hit_rate: 63.85%    per_memory_hit_rate: 63.79%
	407.248s	    output:  234 Mb/s   input:  435 Mb/s
	======================================================================
	 8388608	     48548 ns/op	     13140 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	829.725s

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

	./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6379

	./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6380

测试结果
-------
::

	go test -bench=^BenchmarkRedis2$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-4   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4945884
	VmHWM: 1086360 kB   hit_rate: 58.96%    per_memory_hit_rate: 56.92%
	VmHWM: 1085976 kB
	478.144s	    output:  226 Mb/s   input:  344 Mb/s
	======================================================================
	 8388608	     56999 ns/op	      9986 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	969.150s

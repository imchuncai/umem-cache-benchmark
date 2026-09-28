.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-2g-8k
====================

结论
====
::

	Umem-cache的命中率比Memcached高9%，比Redis高12%。
	Umem-cache的命中吞吐量比Memcached高12%，比Redis高29%。

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
	-args true 2147483648 8192 16 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4971600
	VmHWM: 2120896 kB   hit_rate: 59.27%    per_memory_hit_rate: 58.60%
	546.445s	    output:  197 Mb/s   input:  294 Mb/s
	======================================================================
	 8388608	     65141 ns/op	      8996 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1116.719s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

编译命令
-------
::

	make MEM_LIMIT=2147483648 THREAD_NR=1 MAX_CONN=48

运行命令
-------
::

	taskset -c 1 ./umem-cache 10047

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5356172
	VmHWM: 2098304 kB   hit_rate: 63.85%    per_memory_hit_rate: 63.82%
	532.885s	    output:  178 Mb/s   input:  325 Mb/s
	======================================================================
	 8388608	     63525 ns/op	     10046 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1063.805s

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
	--maxmemory 2147483648 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6379

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4919563
	VmHWM: 2153256 kB   hit_rate: 58.65%    per_memory_hit_rate: 57.12%
	615.716s	    output:  177 Mb/s   input:  259 Mb/s
	======================================================================
	 8388608	     73399 ns/op	      7782 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1245.649s

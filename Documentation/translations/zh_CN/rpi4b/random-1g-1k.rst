.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-1g-1k
====================

结论
====
::

	Umem-cache的命中率比Memcached高10%，比Redis高14%。
	Umem-cache的命中吞吐量比Memcached高51%，比Redis高73%。

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

	taskset -c 1 ./memcached --memory-limit=1024 \
	--max-item-size=1048576 -t 1

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20519319
	VmHWM: 1076840 kB   hit_rate: 61.15%    per_memory_hit_rate: 59.55%
	1552.207s	    output:   34 Mb/s   input:   55 Mb/s
	======================================================================
	33554432	     46259 ns/op	     12872 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3161.324s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

编译命令
-------
::

	make MEM_LIMIT=1073741824 THREAD_NR=1 MAX_CONN=48

运行命令
-------
::

	taskset -c 1 ./umem-cache 10047

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22017478
	VmHWM: 1049720 kB   hit_rate: 65.62%    per_memory_hit_rate: 65.55%
	1128.043s	    output:   41 Mb/s   input:   81 Mb/s
	======================================================================
	33554432	     33618 ns/op	     19497 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2306.516s

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
	--maxmemory 1073741824 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6379

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20003221
	VmHWM: 1084664 kB   hit_rate: 59.61%    per_memory_hit_rate: 57.63%
	1718.700s	    output:   32 Mb/s   input:   48 Mb/s
	======================================================================
	33554432	     51221 ns/op	     11251 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3465.425s

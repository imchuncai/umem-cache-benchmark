.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-2g-1m
====================

结论
====
::

	Umem-cache的命中率比Memcached高14%，比Redis高17%。
	Umem-cache的命中吞吐量比Memcached高8%，比Redis高13%。

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

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=2097152 -t 3

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 16 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30606
	VmHWM: 2123100 kB   hit_rate: 46.70%    per_memory_hit_rate: 46.13%
	172.923s	    output:  809 Mb/s   input:  778 Mb/s
	======================================================================
	   65536	   2638597 ns/op	       175 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	347.590s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

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

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34561
	VmHWM: 2098244 kB   hit_rate: 52.74%    per_memory_hit_rate: 52.71%
	183.148s	    output:  675 Mb/s   input:  823 Mb/s
	======================================================================
	   65536	   2794611 ns/op	       189 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	363.982s

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
	--maxmemory 715827882 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6380

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6381

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis3$ -benchtime=65536x \
	-args true 2147483648 1048576 16 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis3-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31441
	VmHWM:  737396 kB   hit_rate: 47.98%    per_memory_hit_rate: 45.23%
	VmHWM:  741776 kB
	VmHWM:  745168 kB
	176.856s	    output:  768 Mb/s   input:  783 Mb/s
	======================================================================
	   65536	   2698607 ns/op	       168 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	359.229s

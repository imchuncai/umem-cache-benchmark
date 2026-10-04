.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
基准测试-random-2g-1m
====================

结论
====
::

	Umem-cache的命中率比Memcached高15%，比Redis高13%。
	Umem-cache的命中吞吐量比Memcached高15%，比Redis高12%。

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
	server:     4096    warmup:    65536    get:    65536    hit:    30641
	VmHWM: 2120972 kB   hit_rate: 46.75%    per_memory_hit_rate: 46.23%
	P1: 1224 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	171.942s	    output:  816 Mb/s   input:  751 Mb/s
	======================================================================
	   65536	   2623624 ns/op	       176 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	348.334s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

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
	server:     4096    warmup:    65536    get:    65536    hit:    34859
	VmHWM: 2098296 kB   hit_rate: 53.19%    per_memory_hit_rate: 53.16%
	P1: 852 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	171.863s	    output:  717 Mb/s   input:  852 Mb/s
	======================================================================
	   65536	   2622417 ns/op	       203 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	343.530s

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
	server:     4096    warmup:    65536    get:    65536    hit:    31555
	VmHWM: 2141336 kB   hit_rate: 48.15%    per_memory_hit_rate: 47.16%
	P1: 1641 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	169.389s	    output:  806 Mb/s   input:  785 Mb/s
	======================================================================
	   65536	   2584671 ns/op	       182 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	344.143s

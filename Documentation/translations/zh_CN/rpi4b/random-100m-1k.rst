.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

=======================
基准测试-随机大小-100M-1K
=======================

结论
====

性能测试结果显示MEMCACHED和UMEM-CACHE的缓存性能表现十分接近，有较高的命中率和较快的速度。
UMEM-CACHE比REDIS和POGOCACHE在速度上快了60%以上，同时命中率还要高15%。

MEMCACHED
=========
::

	commit 4b9e6198fc44c9eb3ae80802a1b0dcbaf9602969

编译命令
-------
::

	./autogen.sh
	./configure
	make -j

运行命令
-------
::

	./memcached --conn-limit=512 --memory-limit=100 --max-item-size=1048576 -t 4 -u root

测试结果
-------
::

	go test -bench=^BenchmarkMemcached$ -benchtime=3276800x		       \
	-args true 104857600 1024 20 80 80 16 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-4   	
	=======================================================
	case:  819200    hot:  163840(20%)    hot_access: 80% 
	get: 3276800    hit: 1996788    hit_rate: 60.94% 
	hot: 2621031    hit: 1943341    hit_rate: 74.14% 
	VmHWM:  106664 kB    per_memory_hit_rate: 58.50%
	80.812s
	=======================================================
	3276800	     24662 ns/op	      23721 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	163.600s

UMEM-CACHE
==========
::

	commit 53f97eb219364fb18e15431e069b2ceef877b5d9

编译命令
-------
::

	make

运行命令
-------
::

	./umem-cache 10047

测试结果
-------
::

	go test -bench=^BenchmarkUmemCache$ -benchtime=3276800x		       \
	-args true 104857600 1024 20 80 80 16 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-4   	
	=======================================================
	case:  819200    hot:  163840(20%)    hot_access: 80% 
	get: 3276800    hit: 1977711    hit_rate: 60.35% 
	hot: 2621031    hit: 1925530    hit_rate: 73.46% 
	VmHWM:  103828 kB    per_memory_hit_rate: 59.52%
	79.524s
	=======================================================
	3276800	     24269 ns/op	      24527 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	161.030s

REDIS
=====
::

	commit e6e0cf5764c99fc1414e46197126e84360536df6

编译命令
-------
::

	make -j

运行命令
-------
::

	./src/redis-server --protected-mode no --appendonly no --save ""       \
	--maxmemory 52428800 --maxclients 512 --maxmemory-policy allkeys-lru --port 6379

	./src/redis-server --protected-mode no --appendonly no --save ""       \
	--maxmemory 52428800 --maxclients 512 --maxmemory-policy allkeys-lru --port 6380

测试结果
-------
::

	go test -bench=^BenchmarkRedis2$ -benchtime=3276800x		       \
	-args true 104857600 1024 20 80 80 16 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-4   	
	=======================================================
	case:  819200    hot:  163840(20%)    hot_access: 80% 
	get: 3276800    hit: 1903939    hit_rate: 58.10% 
	hot: 2621031    hit: 1852229    hit_rate: 70.67% 
	VmHWM:   60464 kB    per_memory_hit_rate: 49.20%
	VmHWM:   60464 kB
	108.184s
	=======================================================
	3276800	     33015 ns/op	      14903 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	218.985s

POGOCACHE
=========
::

	commit 71972a9f161d96d91b0f67bfe28897d00bfbd49b

编译命令
-------
::

	make -j NOMIMALLOC=1

运行命令
-------
::

	./pogocache --threads=4 --maxmemory=104857600 --maxconns=512 --port=9401 -h 192.168.101.10

测试结果
-------
::

	go test -bench=^BenchmarkPogocache$ -benchtime=3276800x		       \
	-args true 104857600 1024 20 80 80 16 192.168.101.10
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkPogocache-4   	
	=======================================================
	case:  819200    hot:  163840(20%)    hot_access: 80% 
	get: 3276800    hit: 1860590    hit_rate: 56.78% 
	hot: 2621031    hit: 1798916    hit_rate: 68.63% 
	VmHWM:  112324 kB    per_memory_hit_rate: 51.76%
	129.028s
	=======================================================
	3276800	     39376 ns/op	      13146 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	258.677s

.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
BENCHMARK-RANDOM-2G-1M
======================

CONCLUSION
==========

The test results showed that the performance of MEMCACHED, UMEM-CACHE, and REDIS
are very close, and UMEM-CACHE's hit rate is about 12% higher than REDIS. And
POGOCACHE really not respect the memory limit, it used extra 50% of the memory.

The reason that the performance is close among these apps, is that the cache
value is relatively large, the performance bottleneck is at the server size
network output. The server side output speed is more than 850Mb/s and we are
under Gigabit network.

MEMCACHED
=========
::

	commit 4b9e6198fc44c9eb3ae80802a1b0dcbaf9602969

BUILD COMMAND
-------------
::

	./autogen.sh
	./configure
	make -j

RUN COMMAND
-----------
::

	./memcached --conn-limit=512 --memory-limit=2048 --max-item-size=2097152 -t 4 -u root

TEST RESULT
-----------
::

	go test -bench=^BenchmarkMemcached$ -benchtime=65536x		       \
	-args true 2147483648 1048576 20 80 80 16 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-4   	
	=======================================================
	case:   16384    hot:    3276(20%)    hot_access: 80% 
	get:   65536    hit:   43398    hit_rate: 66.22% 
	hot:   52582    hit:   42203    hit_rate: 80.26% 
	VmHWM: 2111664 kB    per_memory_hit_rate: 65.76%
	201.886s
	=======================================================
	65536	   3080534 ns/op	      213 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	401.022s

	out IO speed: 860Mb/s

UMEM-CACHE
==========
::

	commit 53f97eb219364fb18e15431e069b2ceef877b5d9

BUILD COMMAND
-------------
::

	make MEM_LIMIT=2147483648

RUN COMMAND
-----------
::

	./umem-cache 10047

TEST RESULT
-----------
::

	go test -bench=^BenchmarkUmemCache$ -benchtime=65536x		       \
	-args true 2147483648 1048576 20 80 80 16 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-4   	
	=======================================================
	case:   16384    hot:    3276(20%)    hot_access: 80% 
	get:   65536    hit:   45057    hit_rate: 68.75% 
	hot:   52582    hit:   43690    hit_rate: 83.09% 
	VmHWM: 2098372 kB    per_memory_hit_rate: 68.71%
	207.934s
	=======================================================
	65536	   3172826 ns/op	      217 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	418.542s

	out IO speed: 867Mb/s

REDIS
=====
::

	commit e6e0cf5764c99fc1414e46197126e84360536df6

BUILD COMMAND
-------------
::

	make -j

RUN COMMAND
-----------
::

	./src/redis-server --protected-mode no --appendonly no --save ""       \
	--maxmemory 2147483648 --maxclients 512 --maxmemory-policy allkeys-lru --port 6379

TEST RESULT
-----------
::

	go test -bench=^BenchmarkRedis1$ -benchtime=65536x		       \
	-args true 2147483648 1048576 20 80 80 16 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis1-4   	
	=======================================================
	case:   16384    hot:    3276(20%)    hot_access: 80% 
	get:   65536    hit:   41561    hit_rate: 63.42% 
	hot:   52582    hit:   40311    hit_rate: 76.66% 
	VmHWM: 2159772 kB    per_memory_hit_rate: 61.58%
	195.393s
	=======================================================
	65536	   2981461 ns/op	      207 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	389.165s

	out IO speed: 851Mb/s

POGOCACHE
=========
::

	commit 71972a9f161d96d91b0f67bfe28897d00bfbd49b

BUILD COMMAND
-------------
::

	make -j NOMIMALLOC=1

RUN COMMAND
-----------
::

	./pogocache --threads=4 --maxmemory=2147483648 --maxconns=512 --port=9401 -h 192.168.101.10

TEST RESULT
-----------
::

	go test -bench=^BenchmarkPogocache$ -benchtime=65536x		       \
	-args true 2147483648 1048576 20 80 80 16 192.168.101.10
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkPogocache-4   	
	=======================================================
	case:   16384    hot:    3276(20%)    hot_access: 80% 
	get:   65536    hit:   46512    hit_rate: 70.97% 
	hot:   52582    hit:   43827    hit_rate: 83.35% 
	VmHWM: 3312024 kB    per_memory_hit_rate: 44.94%
	214.533s
	=======================================================
	65536	   3273515 ns/op	      137 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	426.389s

	out IO speed: 867Mb/s

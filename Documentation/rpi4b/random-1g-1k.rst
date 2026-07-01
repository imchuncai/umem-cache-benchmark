.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-1g-1k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 9% higher than Memcached and 14% higher than Redis.
	Umem-cache's hit throughput is 4% higher than Memcached and 59% higher than Redis.

Memcached
=========
::

	commit f1674f0231e5d291db474c4ad297f5f069d32521

Build Command
-------------
::

	./autogen.sh
	./configure --enable-tls
	make -j

Run Command
-----------
::

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=1024 \
	--max-item-size=1048576 -t 3 -u root

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20517888
	VmHWM: 1069876 kB   hit_rate: 61.15%    per_memory_hit_rate: 59.93%
	804.568s	    output:   65 Mb/s   input:  108 Mb/s
	======================================================================
	33554432	     23978 ns/op	     24994 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1640.953s

Umem-cache
==========
::

	commit e4a931ea3ee8fc82b3693333fa12a264e36f3dd0

Build Command
-------------
::

	make MEM_LIMIT=1073741824 THREAD_NR=3

Run Command
-----------
::

	taskset -c 1,2,3 ./umem-cache 10047

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22008681
	VmHWM: 1050180 kB   hit_rate: 65.59%    per_memory_hit_rate: 65.49%
	848.809s	    output:   55 Mb/s   input:  109 Mb/s
	======================================================================
	33554432	     25296 ns/op	     25889 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1747.198s

Redis
=====
::

	commit 138263a1b480fcd2e756be27f369203a46481d06

Build Command
-------------
::

	make -j BUILD_TLS=yes

Run Command
-----------
::

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6380

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20142915
	VmHWM:  547412 kB   hit_rate: 60.03%    per_memory_hit_rate: 57.53%
	VmHWM:  546748 kB
	1183.075s	    output:   46 Mb/s   input:   72 Mb/s
	======================================================================
	33554432	     35258 ns/op	     16317 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2393.508s

.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-1g-1k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 9% higher than Memcached and 14% higher than Redis.
	Umem-cache's hit throughput is 14% higher than Memcached and 58% higher than Redis.

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

	./memcached --conn-limit=512 --memory-limit=1024 \
	--max-item-size=1048576 -t 4 -u root

Test Result
-----------
::

	go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-4   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20520204
	VmHWM: 1070020 kB   hit_rate: 61.15%    per_memory_hit_rate: 59.93%
	823.355s	    output:   64 Mb/s   input:  105 Mb/s
	======================================================================
	33554432	     24538 ns/op	     24423 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1682.248s

Umem-cache
==========
::

	commit eeeac62ec7a2ec135b4f2e419a533ba8b3282ccc

Build Command
-------------
::

	make MEM_LIMIT=1073741824

Run Command
-----------
::

	./umem-cache 10047

Test Result
-----------
::

	go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-4   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22018167
	VmHWM: 1050208 kB   hit_rate: 65.62%    per_memory_hit_rate: 65.52%
	792.383s	    output:   59 Mb/s   input:  117 Mb/s
	======================================================================
	33554432	     23615 ns/op	     27744 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1630.474s

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

	./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6379

	./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6380

Test Result
-----------
::

	go test -bench=^BenchmarkRedis2$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-4   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20166736
	VmHWM:  547228 kB   hit_rate: 60.10%    per_memory_hit_rate: 57.59%
	VmHWM:  547052 kB
	1102.509s	    output:   49 Mb/s   input:   77 Mb/s
	======================================================================
	33554432	     32857 ns/op	     17528 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2232.995s

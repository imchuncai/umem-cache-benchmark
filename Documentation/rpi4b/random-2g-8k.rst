.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-2g-8k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 8% higher than Memcached and 12% higher than Redis.
	Umem-cache's hit throughput is 24% higher than Memcached and 31% higher than Redis.

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

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=1048576 -t 3 -u root

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4972504
	VmHWM: 2109712 kB   hit_rate: 59.28%    per_memory_hit_rate: 58.92%
	444.308s	    output:  242 Mb/s   input:  371 Mb/s
	======================================================================
	 8388608	     52966 ns/op	     11125 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	899.984s

Umem-cache
==========
::

	commit 855aee6d8e727184d9a23806597403ded1d941b4

Build Command
-------------
::

	make MEM_LIMIT=2147483648 THREAD_NR=3

Run Command
-----------
::

	taskset -c 1,2,3 ./umem-cache 10047

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5364253
	VmHWM: 2098776 kB   hit_rate: 63.95%    per_memory_hit_rate: 63.90%
	389.924s	    output:  243 Mb/s   input:  455 Mb/s
	======================================================================
	 8388608	     46483 ns/op	     13747 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	793.432s

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
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6380

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4950236
	VmHWM: 1085300 kB   hit_rate: 59.01%    per_memory_hit_rate: 57.00%
	VmHWM: 1085836 kB
	456.566s	    output:  236 Mb/s   input:  360 Mb/s
	======================================================================
	 8388608	     54427 ns/op	     10473 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	923.758s

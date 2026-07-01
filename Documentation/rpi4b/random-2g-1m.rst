.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-2g-1m
======================

Conclusion
==========
::

	Umem-cache's hit rate is 14% higher than Memcached and 14% higher than Redis.
	Umem-cache's hit throughput is 12% higher than Memcached and 7% higher than Redis.

	Note: network throughput is nearing the limit of gigabit networks.

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
	--max-item-size=2097152 -t 3 -u root

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30626
	VmHWM: 2117068 kB   hit_rate: 46.73%    per_memory_hit_rate: 46.29%
	177.015s	    output:  790 Mb/s   input:  760 Mb/s
	======================================================================
	   65536	   2701029 ns/op	       171 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	355.644s

Umem-cache
==========
::

	commit e4a931ea3ee8fc82b3693333fa12a264e36f3dd0

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

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34581
	VmHWM: 2098536 kB   hit_rate: 52.77%    per_memory_hit_rate: 52.73%
	179.903s	    output:  687 Mb/s   input:  838 Mb/s
	======================================================================
	   65536	   2745095 ns/op	       192 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	359.770s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=65536x \
	-args true 2147483648 1048576 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31510
	VmHWM: 1079720 kB   hit_rate: 48.08%    per_memory_hit_rate: 46.44%
	VmHWM: 1091596 kB
	170.005s	    output:  800 Mb/s   input:  814 Mb/s
	======================================================================
	   65536	   2594066 ns/op	       179 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	343.199s

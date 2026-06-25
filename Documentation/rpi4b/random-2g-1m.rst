.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-2g-1m
======================

Conclusion
==========
::

	Umem-cache's hit rate is 14% higher than Memcached and 14% higher than Redis.
	Umem-cache's hit throughput is 10% higher than Memcached and 7% higher than Redis.

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
	server:     4096    warmup:    65536    get:    65536    hit:    30595
	VmHWM: 2117068 kB   hit_rate: 46.68%    per_memory_hit_rate: 46.25%
	173.567s	    output:  806 Mb/s   input:  775 Mb/s
	======================================================================
	   65536	   2648429 ns/op	       175 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	349.764s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 0 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34586
	VmHWM: 2098508 kB   hit_rate: 52.77%    per_memory_hit_rate: 52.74%
	178.876s	    output:  690 Mb/s   input:  843 Mb/s
	======================================================================
	   65536	   2729426 ns/op	       193 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	359.253s

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
	server:     4096    warmup:    65536    get:    65536    hit:    31551
	VmHWM: 1080268 kB   hit_rate: 48.14%    per_memory_hit_rate: 46.46%
	VmHWM: 1092860 kB
	168.703s	    output:  805 Mb/s   input:  821 Mb/s
	======================================================================
	   65536	   2574206 ns/op	       180 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	341.549s

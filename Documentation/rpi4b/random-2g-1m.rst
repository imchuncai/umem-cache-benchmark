.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-2g-1m
======================

Conclusion
==========
::

	Umem-cache's hit rate is 15% higher than Memcached and 12% higher than Redis.
	Umem-cache's hit throughput is 15% higher than Memcached and 9% higher than Redis.

	Note: network throughput is nearing the limit of gigabit networks.

Memcached
=========
::

	commit e44dd0b01234bc0faf970e9225e3423e98022129

Build Command
-------------
::

	./autogen.sh
	./configure --enable-tls
	make -j

Run Command
-----------
::

	taskset -c 1 ./memcached --memory-limit=2048 \
	--max-item-size=2097152 -t 1

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30572
	VmHWM: 2121028 kB   hit_rate: 46.65%    per_memory_hit_rate: 46.12%
	175.217s	    output:  802 Mb/s   input:  736 Mb/s
	======================================================================
	   65536	   2673594 ns/op	       173 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	353.100s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

Build Command
-------------
::

	make MEM_LIMIT=2147483648 THREAD_NR=1 MAX_CONN=48

Run Command
-----------
::

	taskset -c 1 ./umem-cache 10047

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34831
	VmHWM: 2098304 kB   hit_rate: 53.15%    per_memory_hit_rate: 53.12%
	175.025s	    output:  704 Mb/s   input:  836 Mb/s
	======================================================================
	   65536	   2670673 ns/op	       199 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	353.279s

Redis
=====
::

	commit 6bf6224c3dad518329ddc893ef9c5d58dcbabdeb

Build Command
-------------
::

	make -j BUILD_TLS=yes

Run Command
-----------
::

	taskset -c 1 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 2147483648 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6379

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31516
	VmHWM: 2132268 kB   hit_rate: 48.09%    per_memory_hit_rate: 47.30%
	170.487s	    output:  799 Mb/s   input:  782 Mb/s
	======================================================================
	   65536	   2601418 ns/op	       182 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	344.792s

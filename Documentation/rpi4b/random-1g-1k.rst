.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-1g-1k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 10% higher than Memcached and 16% higher than Redis.
	Umem-cache's hit throughput is 15% higher than Memcached and 66% higher than Redis.

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

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=1024 \
	--max-item-size=1048576 -t 3

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20521042
	VmHWM: 1077388 kB   hit_rate: 61.16%    per_memory_hit_rate: 59.52%
	897.954s	    output:   59 Mb/s   input:   96 Mb/s
	======================================================================
	33554432	     26761 ns/op	     22242 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1827.614s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

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
	-args true 1073741824 1024 16 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22008957
	VmHWM: 1049940 kB   hit_rate: 65.59%    per_memory_hit_rate: 65.51%
	859.583s	    output:   54 Mb/s   input:  108 Mb/s
	======================================================================
	33554432	     25618 ns/op	     25571 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1756.945s

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

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6380

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6381

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis3$ -benchtime=33554432x \
	-args true 1073741824 1024 16 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis3-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20122968
	VmHWM:  372400 kB   hit_rate: 59.97%    per_memory_hit_rate: 56.29%
	VmHWM:  372320 kB
	VmHWM:  372428 kB
	1228.212s	    output:   44 Mb/s   input:   69 Mb/s
	======================================================================
	33554432	     36604 ns/op	     15378 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2481.734s

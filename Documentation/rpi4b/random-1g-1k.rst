.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-1g-1k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 10% higher than Memcached and 17% higher than Redis.
	Umem-cache's hit throughput is 14% higher than Memcached and 65% higher than Redis.

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

	taskset -c 1,2,3 ./memcached --memory-limit=1024 \
	--max-item-size=1048576 -t 3

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20515592
	VmHWM: 1077396 kB   hit_rate: 61.14%    per_memory_hit_rate: 59.51%
	905.427s	    output:   58 Mb/s   input:   94 Mb/s
	======================================================================
	33554432	     26984 ns/op	     22052 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1858.393s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

Build Command
-------------
::

	make MEM_LIMIT=1073741824 THREAD_NR=3 MAX_CONN=144

Run Command
-----------
::

	taskset -c 1,2,3 ./umem-cache 10047

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22007175
	VmHWM: 1049784 kB   hit_rate: 65.59%    per_memory_hit_rate: 65.51%
	877.269s	    output:   53 Mb/s   input:  104 Mb/s
	======================================================================
	33554432	     26145 ns/op	     25057 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1816.785s

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
	--maxmemory 357913941 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6380

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6381

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=33554432x \
	-args true 1073741824 1024 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20113957
	VmHWM:  374368 kB   hit_rate: 59.94%    per_memory_hit_rate: 56.18%
	VmHWM:  372212 kB
	VmHWM:  372260 kB
	1240.469s	    output:   44 Mb/s   input:   67 Mb/s
	======================================================================
	33554432	     36969 ns/op	     15196 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2525.910s

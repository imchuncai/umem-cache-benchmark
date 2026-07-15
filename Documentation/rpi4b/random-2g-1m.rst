.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-2g-1m
======================

Conclusion
==========
::

	Umem-cache's hit rate is 14% higher than Memcached and 19% higher than Redis.
	Umem-cache's hit throughput is 15% higher than Memcached and 23% higher than Redis.

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

	taskset -c 1,2,3 ./memcached --memory-limit=2048 \
	--max-item-size=2097152 -t 3

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30639
	VmHWM: 2123108 kB   hit_rate: 46.75%    per_memory_hit_rate: 46.18%
	177.412s	    output:  792 Mb/s   input:  728 Mb/s
	======================================================================
	   65536	   2707091 ns/op	       171 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	357.110s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

Build Command
-------------
::

	make MEM_LIMIT=2147483648 THREAD_NR=3 MAX_CONN=144

Run Command
-----------
::

	taskset -c 1,2,3 ./umem-cache 10047

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34594
	VmHWM: 2098168 kB   hit_rate: 52.79%    per_memory_hit_rate: 52.76%
	176.642s	    output:  703 Mb/s   input:  822 Mb/s
	======================================================================
	   65536	   2695336 ns/op	       196 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	357.620s

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
	--maxmemory 715827882 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6380

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 6381

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=65536x \
	-args true 2147483648 1048576 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31372
	VmHWM:  749052 kB   hit_rate: 47.87%    per_memory_hit_rate: 44.35%
	VmHWM:  738556 kB
	VmHWM:  775816 kB
	181.998s	    output:  754 Mb/s   input:  727 Mb/s
	======================================================================
	   65536	   2777064 ns/op	       160 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	365.863s

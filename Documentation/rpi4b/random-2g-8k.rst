.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-2g-8k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 9% higher than Memcached and 13% higher than Redis.
	Umem-cache's hit throughput is 23% higher than Memcached and 26% higher than Redis.

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
	--max-item-size=1048576 -t 3

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4972862
	VmHWM: 2120236 kB   hit_rate: 59.28%    per_memory_hit_rate: 58.64%
	423.189s	    output:  254 Mb/s   input:  380 Mb/s
	======================================================================
	 8388608	     50448 ns/op	     11623 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	860.649s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5364087
	VmHWM: 2098388 kB   hit_rate: 63.94%    per_memory_hit_rate: 63.91%
	375.695s	    output:  252 Mb/s   input:  462 Mb/s
	======================================================================
	 8388608	     44786 ns/op	     14269 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	767.114s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=8388608x \
	-args true 2147483648 8192 16 3 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4958147
	VmHWM:  731936 kB   hit_rate: 59.11%    per_memory_hit_rate: 56.51%
	VmHWM:  730552 kB
	VmHWM:  731104 kB
	419.937s	    output:  256 Mb/s   input:  383 Mb/s
	======================================================================
	 8388608	     50060 ns/op	     11288 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	852.510s

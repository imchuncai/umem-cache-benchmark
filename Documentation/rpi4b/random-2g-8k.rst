.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

======================
Benchmark-random-2g-8k
======================

Conclusion
==========
::

	Umem-cache's hit rate is 9% higher than Memcached and 13% higher than Redis.
	Umem-cache's hit throughput is 25% higher than Memcached and 26% higher than Redis.

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

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=1048576 -t 3

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4971376
	VmHWM: 2120328 kB   hit_rate: 59.26%    per_memory_hit_rate: 58.62%
	434.056s	    output:  248 Mb/s   input:  379 Mb/s
	======================================================================
	 8388608	     51743 ns/op	     11328 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	879.957s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

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
	-args true 2147483648 8192 16 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5364050
	VmHWM: 2098496 kB   hit_rate: 63.94%    per_memory_hit_rate: 63.90%
	378.923s	    output:  250 Mb/s   input:  468 Mb/s
	======================================================================
	 8388608	     45171 ns/op	     14147 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	767.017s

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
	--maxmemory 715827882 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6379

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6380

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 6381

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis3$ -benchtime=8388608x \
	-args true 2147483648 8192 16 0 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis3-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4956436
	VmHWM:  731768 kB   hit_rate: 59.09%    per_memory_hit_rate: 56.47%
	VmHWM:  731784 kB
	VmHWM:  730724 kB
	422.333s	    output:  255 Mb/s   input:  390 Mb/s
	======================================================================
	 8388608	     50346 ns/op	     11216 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	855.271s

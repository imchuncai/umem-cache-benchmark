.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-2g-8k
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 8% higher than Memcached and 13% higher than Redis.
	Umem-cache's hit throughput is 12% higher than Memcached and 51% higher than Redis.

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

	./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=1048576 -t 4 -u root \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-4   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4971910
	VmHWM: 2115084 kB   hit_rate: 59.27%    per_memory_hit_rate: 58.77%
	501.834s	    output:  215 Mb/s   input:  328 Mb/s
	======================================================================
	 8388608	     59823 ns/op	      9823 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1012.612s

Umem-cache
==========
::

	commit eeeac62ec7a2ec135b4f2e419a533ba8b3282ccc

Build Command
-------------
::

	make MEM_LIMIT=2147483648 TLS=1

Run Command
-----------
::

	./umem-cache 10047 cert.pem key.pem ca-cert.pem

Test Result
-----------
::

	go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-4   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5355771
	VmHWM: 2104448 kB   hit_rate: 63.85%    per_memory_hit_rate: 63.62%
	486.373s	    output:  196 Mb/s   input:  364 Mb/s
	======================================================================
	 8388608	     57980 ns/op	     10973 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	984.358s

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
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

Test Result
-----------
::

	go test -bench=^BenchmarkRedis2$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-4   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4913728
	VmHWM: 1089072 kB   hit_rate: 58.58%    per_memory_hit_rate: 56.35%
	VmHWM: 1090984 kB
	649.910s	    output:  168 Mb/s   input:  251 Mb/s
	======================================================================
	 8388608	     77475 ns/op	      7273 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1314.784s

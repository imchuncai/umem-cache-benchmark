.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-2g-8k
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 9% higher than Memcached and 13% higher than Redis.
	Umem-cache's hit throughput is 37% higher than Memcached and 29% higher than Redis.

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
	--max-item-size=1048576 -t 1 \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4972578
	VmHWM: 2127232 kB   hit_rate: 59.28%    per_memory_hit_rate: 58.44%
	1279.692s	    output:   84 Mb/s   input:  126 Mb/s
	======================================================================
	 8388608	    152551 ns/op	      3831 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2594.170s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

Build Command
-------------
::

	make MEM_LIMIT=2147483648 THREAD_NR=1 MAX_CONN=48 TLS=1

Run Command
-----------
::

	taskset -c 1 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5356171
	VmHWM: 2105124 kB   hit_rate: 63.85%    per_memory_hit_rate: 63.61%
	1018.833s	    output:   93 Mb/s   input:  170 Mb/s
	======================================================================
	 8388608	    121454 ns/op	      5237 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2041.214s

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
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4853427
	VmHWM: 2159808 kB   hit_rate: 57.86%    per_memory_hit_rate: 56.18%
	1163.205s	    output:   95 Mb/s   input:  135 Mb/s
	======================================================================
	 8388608	    138665 ns/op	      4051 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2340.722s

.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-1g-1k
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 10% higher than Memcached and 15% higher than Redis.
	Umem-cache's hit throughput is 71% higher than Memcached and 64% higher than Redis.

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

	taskset -c 1 ./memcached --memory-limit=1024 \
	--max-item-size=1048576 -t 1 \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20529747
	VmHWM: 1083720 kB   hit_rate: 61.18%    per_memory_hit_rate: 59.20%
	3529.626s	    output:   15 Mb/s   input:   24 Mb/s
	======================================================================
	33554432	    105191 ns/op	      5628 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	7145.379s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

Build Command
-------------
::

	make MEM_LIMIT=1073741824 THREAD_NR=1 MAX_CONN=48 TLS=1

Run Command
-----------
::

	taskset -c 1 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22017478
	VmHWM: 1057196 kB   hit_rate: 65.62%    per_memory_hit_rate: 65.08%
	2267.032s	    output:   20 Mb/s   input:   40 Mb/s
	======================================================================
	33554432	     67563 ns/op	      9633 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	4599.149s

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
	--maxmemory 1073741824 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 19807711
	VmHWM: 1089264 kB   hit_rate: 59.03%    per_memory_hit_rate: 56.83%
	3251.190s	    output:   17 Mb/s   input:   25 Mb/s
	======================================================================
	33554432	     96893 ns/op	      5865 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	6530.094s

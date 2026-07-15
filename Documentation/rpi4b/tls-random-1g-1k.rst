.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-1g-1k
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 10% higher than Memcached and 17% higher than Redis.
	Umem-cache's hit throughput is 53% higher than Memcached and 81% higher than Redis.

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
	--max-item-size=1048576 -t 3 \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 3 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20524150
	VmHWM: 1086204 kB   hit_rate: 61.17%    per_memory_hit_rate: 59.05%
	1364.548s	    output:   39 Mb/s   input:   62 Mb/s
	======================================================================
	33554432	     40667 ns/op	     14520 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2777.233s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

Build Command
-------------
::

	make MEM_LIMIT=1073741824 THREAD_NR=3 MAX_CONN=144 TLS=1

Run Command
-----------
::

	taskset -c 1,2,3 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 3 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22007175
	VmHWM: 1056476 kB   hit_rate: 65.59%    per_memory_hit_rate: 65.10%
	986.141s	    output:   47 Mb/s   input:   92 Mb/s
	======================================================================
	33554432	     29389 ns/op	     22150 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2037.772s

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
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 357913941 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6381 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=33554432x \
	-args true 1073741824 1024 16 3 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20038974
	VmHWM:  376316 kB   hit_rate: 59.72%    per_memory_hit_rate: 55.46%
	VmHWM:  376328 kB
	VmHWM:  376408 kB
	1519.969s	    output:   36 Mb/s   input:   55 Mb/s
	======================================================================
	33554432	     45299 ns/op	     12244 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3092.353s

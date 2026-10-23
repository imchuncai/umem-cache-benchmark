.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-2g-1m
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 15% higher than Memcached and 13% higher than Redis.
	Umem-cache's hit throughput is 34% higher than Memcached and 34% higher than Redis.

	Umem-cache's P90   latency is 0% lower than Memcached and 0% lower than Redis.
	Umem-cache's P99.9 latency is 0% lower than Memcached and 0% lower than Redis.

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
	--max-item-size=2097152 -t 1 \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30636
	VmHWM: 2127924 kB   hit_rate: 46.75%    per_memory_hit_rate: 46.07%
	P1: 9999 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	301.911s	    output:  465 Mb/s   input:  428 Mb/s
	======================================================================
	   65536	   4606804 ns/op	       100 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	609.124s

Umem-cache
==========
::

	commit 079c92dc41d66b84753ccc3d5ae78974460009c0

Build Command
-------------
::

	make MEM_LIMIT=2147483648 THREAD_NR=1 MAX_CONN=24 TLS=1

Run Command
-----------
::

	taskset -c 1 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34872
	VmHWM: 2104832 kB   hit_rate: 53.21%    per_memory_hit_rate: 53.02%
	P1: 9999 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	260.171s	    output:  473 Mb/s   input:  563 Mb/s
	======================================================================
	   65536	   3969891 ns/op	       134 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	524.758s

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
	--maxmemory 2147483648 --maxclients 24 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=65536x \
	-args true 2147483648 1048576 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31100
	VmHWM: 2119440 kB   hit_rate: 47.45%    per_memory_hit_rate: 46.96%
	P1: 5761 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	308.323s	    output:  447 Mb/s   input:  427 Mb/s
	======================================================================
	   65536	   4704630 ns/op	       100 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	619.033s

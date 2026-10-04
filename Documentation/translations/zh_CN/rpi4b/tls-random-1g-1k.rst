.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

========================
基准测试-tls-random-1g-1k
========================

结论
====
::

	Umem-cache的命中率比Memcached高10%，比Redis高14%。
	Umem-cache的命中吞吐量比Memcached高75%，比Redis高60%。

	Umem-cache的P90  延迟比Memcached低46%，比Redis低48%。
	Umem-cache的P99.9延迟比Memcached低70%，比Redis低57%。

Memcached
=========
::

	commit e44dd0b01234bc0faf970e9225e3423e98022129

编译命令
-------
::

	./autogen.sh
	./configure --enable-tls
	make -j

运行命令
-------
::

	taskset -c 1 ./memcached --memory-limit=1024 \
	--max-item-size=1048576 -t 1 \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20530021
	VmHWM: 1083464 kB   hit_rate: 61.18%    per_memory_hit_rate: 59.21%
	P1: 1141 us  P50: 2050 us  P90: 3236 us  P99: 6154 us  P99.9: 7498 us
	3459.679s	    output:   15 Mb/s   input:   24 Mb/s
	======================================================================
	33554432	    103106 ns/op	      5743 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	6994.881s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

编译命令
-------
::

	make MEM_LIMIT=1073741824 THREAD_NR=1 MAX_CONN=24 TLS=1

运行命令
-------
::

	taskset -c 1 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22017475
	VmHWM: 1056928 kB   hit_rate: 65.62%    per_memory_hit_rate: 65.10%
	P1: 407 us  P50: 1561 us  P90: 1761 us  P99: 1969 us  P99.9: 2267 us
	2175.630s	    output:   21 Mb/s   input:   42 Mb/s
	======================================================================
	33554432	     64839 ns/op	     10040 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	4421.300s

Redis
=====
::

	commit 6bf6224c3dad518329ddc893ef9c5d58dcbabdeb

编译命令
-------
::

	make -j BUILD_TLS=yes

运行命令
-------
::

	taskset -c 1 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 24 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=33554432x \
	-args true 1073741824 1024 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 19828554
	VmHWM: 1088364 kB   hit_rate: 59.09%    per_memory_hit_rate: 56.93%
	P1: 607 us  P50: 1887 us  P90: 3397 us  P99: 3979 us  P99.9: 5280 us
	3039.350s	    output:   18 Mb/s   input:   27 Mb/s
	======================================================================
	33554432	     90580 ns/op	      6285 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	6106.233s

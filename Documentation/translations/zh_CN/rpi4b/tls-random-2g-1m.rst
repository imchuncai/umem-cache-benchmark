.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

========================
基准测试-tls-random-2g-1m
========================

结论
====
::

	Umem-cache的命中率比Memcached高15%，比Redis高13%。
	Umem-cache的命中吞吐量比Memcached高34%，比Redis高35%。

	Umem-cache的P90  延迟比Memcached低0%，比Redis低0%。
	Umem-cache的P99.9延迟比Memcached低0%，比Redis低0%。

	注意：网络吞吐量已经接近千兆网络的极限。

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

	taskset -c 1 ./memcached --memory-limit=2048 \
	--max-item-size=2097152 -t 1 \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30611
	VmHWM: 2128280 kB   hit_rate: 46.71%    per_memory_hit_rate: 46.03%
	P1: 9999 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	301.264s	    output:  466 Mb/s   input:  429 Mb/s
	======================================================================
	   65536	   4596917 ns/op	       100 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	607.994s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

编译命令
-------
::

	make MEM_LIMIT=2147483648 THREAD_NR=1 MAX_CONN=24 TLS=1

运行命令
-------
::

	taskset -c 1 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34871
	VmHWM: 2104856 kB   hit_rate: 53.21%    per_memory_hit_rate: 53.01%
	P1: 9999 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	259.353s	    output:  475 Mb/s   input:  565 Mb/s
	======================================================================
	   65536	   3957411 ns/op	       134 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	522.911s

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
	--maxmemory 2147483648 --maxclients 24 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=65536x \
	-args true 2147483648 1048576 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31006
	VmHWM: 2115044 kB   hit_rate: 47.31%    per_memory_hit_rate: 46.91%
	P1: 6178 us  P50: 9999 us  P90: 9999 us  P99: 9999 us  P99.9: 9999 us
	310.056s	    output:  446 Mb/s   input:  423 Mb/s
	======================================================================
	   65536	   4731086 ns/op	        99 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	621.941s

.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

========================
基准测试-tls-random-2g-1m
========================

结论
====
::

	Umem-cache的命中率比Memcached高15%，比Redis高13%。
	Umem-cache的命中吞吐量比Memcached高30%，比Redis高36%。

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
	-args true 2147483648 1048576 16 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30674
	VmHWM: 2127840 kB   hit_rate: 46.80%    per_memory_hit_rate: 46.13%
	295.975s	    output:  474 Mb/s   input:  437 Mb/s
	======================================================================
	   65536	   4516217 ns/op	       102 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	596.123s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

编译命令
-------
::

	make MEM_LIMIT=2147483648 THREAD_NR=1 MAX_CONN=48 TLS=1

运行命令
-------
::

	taskset -c 1 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34874
	VmHWM: 2105328 kB   hit_rate: 53.21%    per_memory_hit_rate: 53.01%
	261.039s	    output:  472 Mb/s   input:  561 Mb/s
	======================================================================
	   65536	   3983141 ns/op	       133 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	523.781s

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
	--maxmemory 2147483648 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31014
	VmHWM: 2121204 kB   hit_rate: 47.32%    per_memory_hit_rate: 46.79%
	312.056s	    output:  443 Mb/s   input:  421 Mb/s
	======================================================================
	   65536	   4761597 ns/op	        98 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	626.513s

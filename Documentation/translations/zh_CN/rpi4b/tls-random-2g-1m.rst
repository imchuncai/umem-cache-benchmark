.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

========================
基准测试-tls-random-2g-1m
========================

结论
====
::

	Umem-cache的命中率比Memcached高14%，比Redis高14%。
	Umem-cache的命中吞吐量比Memcached高11%，比Redis高9%。

	注意：网络吞吐量已经接近千兆网络的极限。

Memcached
=========
::

	commit f1674f0231e5d291db474c4ad297f5f069d32521

编译命令
-------
::

	./autogen.sh
	./configure --enable-tls
	make -j

运行命令
-------
::

	./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=2097152 -t 4 -u root \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

测试结果
-------
::

	go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-4   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30585
	VmHWM: 2123716 kB   hit_rate: 46.67%    per_memory_hit_rate: 46.09%
	188.567s	    output:  742 Mb/s   input:  713 Mb/s
	======================================================================
	   65536	   2877304 ns/op	       160 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	375.205s

Umem-cache
==========
::

	commit eeeac62ec7a2ec135b4f2e419a533ba8b3282ccc

编译命令
-------
::

	make MEM_LIMIT=2147483648 TLS=1

运行命令
-------
::

	./umem-cache 10047 cert.pem key.pem ca-cert.pem

测试结果
-------
::

	go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-4   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34455
	VmHWM: 2103688 kB   hit_rate: 52.57%    per_memory_hit_rate: 52.41%
	193.408s	    output:  641 Mb/s   input:  777 Mb/s
	======================================================================
	   65536	   2951174 ns/op	       178 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	380.972s

Redis
=====
::

	commit 138263a1b480fcd2e756be27f369203a46481d06

编译命令
-------
::

	make -j BUILD_TLS=yes

运行命令
-------
::

	./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

测试结果
-------
::

	go test -bench=^BenchmarkRedis2$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-4   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31350
	VmHWM: 1092152 kB   hit_rate: 47.84%    per_memory_hit_rate: 45.92%
	VmHWM: 1092648 kB
	184.488s	    output:  739 Mb/s   input:  748 Mb/s
	======================================================================
	   65536	   2815057 ns/op	       163 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	369.487s

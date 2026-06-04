.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

========================
基准测试-tls-random-1g-1k
========================

结论
====
::

	Umem-cache的命中率比Memcached高9%，比Redis高15%。
	Umem-cache的命中吞吐量比Memcached高33%，比Redis高81%。

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

	./memcached --conn-limit=512 --memory-limit=1024 \
	--max-item-size=1048576 -t 4 -u root \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

测试结果
-------
::

	go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-4   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20528597
	VmHWM: 1075232 kB   hit_rate: 61.18%    per_memory_hit_rate: 59.66%
	1357.080s	    output:   39 Mb/s   input:   64 Mb/s
	======================================================================
	33554432	     40444 ns/op	     14752 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2754.501s

Umem-cache
==========
::

	commit eeeac62ec7a2ec135b4f2e419a533ba8b3282ccc

编译命令
-------
::

	make MEM_LIMIT=1073741824 TLS=1

运行命令
-------
::

	./umem-cache 10047 cert.pem key.pem ca-cert.pem

测试结果
-------
::

	go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-4   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22018166
	VmHWM: 1056076 kB   hit_rate: 65.62%    per_memory_hit_rate: 65.15%
	1113.776s	    output:   42 Mb/s   input:   83 Mb/s
	======================================================================
	33554432	     33193 ns/op	     19629 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2285.424s

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
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

测试结果
-------
::

	go test -bench=^BenchmarkRedis2$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-4   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 19996865
	VmHWM:  550312 kB   hit_rate: 59.60%    per_memory_hit_rate: 56.77%
	VmHWM:  550404 kB
	1758.592s	    output:   31 Mb/s   input:   48 Mb/s
	======================================================================
	33554432	     52410 ns/op	     10832 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3539.147s

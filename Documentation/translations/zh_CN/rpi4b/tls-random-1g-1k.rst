.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

========================
基准测试-tls-random-1g-1k
========================

结论
====
::

	Umem-cache的命中率比Memcached高9%，比Redis高15%。
	Umem-cache的命中吞吐量比Memcached高60%，比Redis高121%。

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

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=1024 \
	--max-item-size=1048576 -t 3 -u root \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20529037
	VmHWM: 1074748 kB   hit_rate: 61.18%    per_memory_hit_rate: 59.69%
	1346.331s	    output:   39 Mb/s   input:   64 Mb/s
	======================================================================
	33554432	     40124 ns/op	     14877 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2728.265s

Umem-cache
==========
::

	commit e4a931ea3ee8fc82b3693333fa12a264e36f3dd0

编译命令
-------
::

	make MEM_LIMIT=1073741824 THREAD_NR=3 TLS=1

运行命令
-------
::

	taskset -c 1,2,3 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22008664
	VmHWM: 1055556 kB   hit_rate: 65.59%    per_memory_hit_rate: 65.16%
	918.203s	    output:   51 Mb/s   input:  101 Mb/s
	======================================================================
	33554432	     27365 ns/op	     23811 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1878.615s

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

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 19994946
	VmHWM:  550044 kB   hit_rate: 59.59%    per_memory_hit_rate: 56.82%
	VmHWM:  549700 kB
	1768.016s	    output:   31 Mb/s   input:   48 Mb/s
	======================================================================
	33554432	     52691 ns/op	     10783 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3559.069s

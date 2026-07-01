.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

========================
基准测试-tls-random-2g-1m
========================

结论
====
::

	Umem-cache的命中率比Memcached高14%，比Redis高14%。
	Umem-cache的命中吞吐量比Memcached高8%，比Redis高10%。

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

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=2097152 -t 3 -u root \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30604
	VmHWM: 2123272 kB   hit_rate: 46.70%    per_memory_hit_rate: 46.12%
	178.877s	    output:  781 Mb/s   input:  752 Mb/s
	======================================================================
	   65536	   2729445 ns/op	       169 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	357.915s

Umem-cache
==========
::

	commit e4a931ea3ee8fc82b3693333fa12a264e36f3dd0

编译命令
-------
::

	make MEM_LIMIT=2147483648 THREAD_NR=3 TLS=1

运行命令
-------
::

	taskset -c 1,2,3 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34601
	VmHWM: 2103608 kB   hit_rate: 52.80%    per_memory_hit_rate: 52.63%
	189.687s	    output:  651 Mb/s   input:  795 Mb/s
	======================================================================
	   65536	   2894389 ns/op	       182 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	373.517s

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
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31415
	VmHWM: 1085644 kB   hit_rate: 47.94%    per_memory_hit_rate: 46.07%
	VmHWM: 1096264 kB
	183.040s	    output:  744 Mb/s   input:  755 Mb/s
	======================================================================
	   65536	   2792962 ns/op	       165 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	367.961s

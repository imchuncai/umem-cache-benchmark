.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

========================
基准测试-tls-random-2g-8k
========================

结论
====
::

	Umem-cache的命中率比Memcached高8%，比Redis高13%。
	Umem-cache的命中吞吐量比Memcached高30%，比Redis高70%。

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
	--max-item-size=1048576 -t 3 -u root \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4972452
	VmHWM: 2115948 kB   hit_rate: 59.28%    per_memory_hit_rate: 58.75%
	512.419s	    output:  210 Mb/s   input:  321 Mb/s
	======================================================================
	 8388608	     61085 ns/op	      9618 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1032.730s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5364259
	VmHWM: 2104200 kB   hit_rate: 63.95%    per_memory_hit_rate: 63.73%
	428.651s	    output:  221 Mb/s   input:  414 Mb/s
	======================================================================
	 8388608	     51099 ns/op	     12472 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	867.262s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4914309
	VmHWM: 1087888 kB   hit_rate: 58.58%    per_memory_hit_rate: 56.47%
	VmHWM: 1087864 kB
	644.413s	    output:  169 Mb/s   input:  254 Mb/s
	======================================================================
	 8388608	     76820 ns/op	      7351 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1302.925s

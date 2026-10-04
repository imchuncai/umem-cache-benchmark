.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

========================
基准测试-tls-random-2g-8k
========================

结论
====
::

	Umem-cache的命中率比Memcached高9%，比Redis高13%。
	Umem-cache的命中吞吐量比Memcached高36%，比Redis高33%。

	Umem-cache的P90  延迟比Memcached低30%，比Redis低36%。
	Umem-cache的P99.9延迟比Memcached低55%，比Redis低37%。

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
	--max-item-size=1048576 -t 1 \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

测试结果
-------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4972864
	VmHWM: 2125764 kB   hit_rate: 59.28%    per_memory_hit_rate: 58.48%
	P1: 1441 us  P50: 3609 us  P90: 4840 us  P99: 8811 us  P99.9: 9552 us
	1246.616s	    output:   86 Mb/s   input:  129 Mb/s
	======================================================================
	 8388608	    148608 ns/op	      3935 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2525.941s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5356169
	VmHWM: 2103492 kB   hit_rate: 63.85%    per_memory_hit_rate: 63.66%
	P1: 854 us  P50: 2849 us  P90: 3365 us  P99: 3802 us  P99.9: 4300 us
	995.393s	    output:   95 Mb/s   input:  174 Mb/s
	======================================================================
	 8388608	    118660 ns/op	      5365 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1991.117s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=8388608x \
	-args true 2147483648 8192 8 1 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4852507
	VmHWM: 2160248 kB   hit_rate: 57.85%    per_memory_hit_rate: 56.16%
	P1: 1192 us  P50: 2936 us  P90: 5220 us  P99: 5982 us  P99.9: 6826 us
	1169.717s	    output:   95 Mb/s   input:  135 Mb/s
	======================================================================
	 8388608	    139441 ns/op	      4027 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2359.410s

.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-1g-1k
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 9% higher than Memcached and 15% higher than Redis.
	Umem-cache's hit throughput is 58% higher than Memcached and 113% higher than Redis.

Memcached
=========
::

	commit f1674f0231e5d291db474c4ad297f5f069d32521

Build Command
-------------
::

	./autogen.sh
	./configure --enable-tls
	make -j

Run Command
-----------
::

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=1024 \
	--max-item-size=1048576 -t 3 -u root \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20522647
	VmHWM: 1074656 kB   hit_rate: 61.16%    per_memory_hit_rate: 59.68%
	1335.893s	    output:   39 Mb/s   input:   65 Mb/s
	======================================================================
	33554432	     39813 ns/op	     14990 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	2717.115s

Umem-cache
==========
::

	commit 855aee6d8e727184d9a23806597403ded1d941b4

Build Command
-------------
::

	make MEM_LIMIT=1073741824 THREAD_NR=3 TLS=1

Run Command
-----------
::

	taskset -c 1,2,3 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 22008671
	VmHWM: 1055476 kB   hit_rate: 65.59%    per_memory_hit_rate: 65.16%
	923.517s	    output:   50 Mb/s   input:  100 Mb/s
	======================================================================
	33554432	     27523 ns/op	     23676 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1888.701s

Redis
=====
::

	commit 138263a1b480fcd2e756be27f369203a46481d06

Build Command
-------------
::

	make -j BUILD_TLS=yes

Run Command
-----------
::

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 536870912 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=33554432x \
	-args true 1073741824 1024 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:  2097152    warmup: 33554432    get: 33554432    hit: 20005314
	VmHWM:  550396 kB   hit_rate: 59.62%    per_memory_hit_rate: 56.83%
	VmHWM:  549572 kB
	1714.169s	    output:   32 Mb/s   input:   49 Mb/s
	======================================================================
	33554432	     51086 ns/op	     11125 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	3452.391s

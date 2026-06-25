.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-2g-8k
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 8% higher than Memcached and 13% higher than Redis.
	Umem-cache's hit throughput is 30% higher than Memcached and 77% higher than Redis.

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

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=1048576 -t 3 -u root \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4972451
	VmHWM: 2114632 kB   hit_rate: 59.28%    per_memory_hit_rate: 58.79%
	500.112s	    output:  215 Mb/s   input:  329 Mb/s
	======================================================================
	 8388608	     59618 ns/op	      9860 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1021.187s

Umem-cache
==========
::

	commit 855aee6d8e727184d9a23806597403ded1d941b4

Build Command
-------------
::

	make MEM_LIMIT=2147483648 THREAD_NR=3 TLS=1

Run Command
-----------
::

	taskset -c 1,2,3 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5364258
	VmHWM: 2104120 kB   hit_rate: 63.95%    per_memory_hit_rate: 63.74%
	416.057s	    output:  228 Mb/s   input:  427 Mb/s
	======================================================================
	 8388608	     49598 ns/op	     12850 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	838.213s

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
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 1073741824 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4912989
	VmHWM: 1088212 kB   hit_rate: 58.57%    per_memory_hit_rate: 56.40%
	VmHWM: 1089436 kB
	652.618s	    output:  167 Mb/s   input:  250 Mb/s
	======================================================================
	 8388608	     77798 ns/op	      7250 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1317.370s

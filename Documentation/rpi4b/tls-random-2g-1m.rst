.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-2g-1m
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 14% higher than Memcached and 14% higher than Redis.
	Umem-cache's hit throughput is 5% higher than Memcached and 7% higher than Redis.

	Note: network throughput is nearing the limit of gigabit networks.

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
	--max-item-size=2097152 -t 3 -u root \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30593
	VmHWM: 2123144 kB   hit_rate: 46.68%    per_memory_hit_rate: 46.11%
	177.910s	    output:  785 Mb/s   input:  756 Mb/s
	======================================================================
	   65536	   2714698 ns/op	       170 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	361.525s

Umem-cache
==========
::

	commit eeeac62ec7a2ec135b4f2e419a533ba8b3282ccc

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

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34599
	VmHWM: 2103788 kB   hit_rate: 52.79%    per_memory_hit_rate: 52.63%
	192.653s	    output:  641 Mb/s   input:  783 Mb/s
	======================================================================
	   65536	   2939658 ns/op	       179 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	387.823s

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

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis2$ -benchtime=65536x \
	-args true 2147483648 1048576 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31441
	VmHWM: 1087372 kB   hit_rate: 47.98%    per_memory_hit_rate: 46.26%
	VmHWM: 1087728 kB
	181.950s	    output:  749 Mb/s   input:  758 Mb/s
	======================================================================
	   65536	   2776337 ns/op	       167 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	365.497s

.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-2g-1m
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 14% higher than Memcached and 18% higher than Redis.
	Umem-cache's hit throughput is 17% higher than Memcached and 19% higher than Redis.

	Note: network throughput is nearing the limit of gigabit networks.

Memcached
=========
::

	commit e44dd0b01234bc0faf970e9225e3423e98022129

Build Command
-------------
::

	./autogen.sh
	./configure --enable-tls
	make -j

Run Command
-----------
::

	taskset -c 1,2,3 ./memcached --memory-limit=2048 \
	--max-item-size=2097152 -t 3 \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=65536x \
	-args true 2147483648 1048576 16 3 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    30602
	VmHWM: 2129664 kB   hit_rate: 46.69%    per_memory_hit_rate: 45.98%
	197.174s	    output:  708 Mb/s   input:  683 Mb/s
	======================================================================
	   65536	   3008630 ns/op	       153 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	388.573s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

Build Command
-------------
::

	make MEM_LIMIT=2147483648 THREAD_NR=3 MAX_CONN=144 TLS=1

Run Command
-----------
::

	taskset -c 1,2,3 ./umem-cache 10047 cert.pem key.pem ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkUmemCache$ -benchtime=65536x \
	-args true 2147483648 1048576 16 3 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    34571
	VmHWM: 2103544 kB   hit_rate: 52.75%    per_memory_hit_rate: 52.59%
	192.765s	    output:  641 Mb/s   input:  782 Mb/s
	======================================================================
	   65536	   2941368 ns/op	       179 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	385.953s

Redis
=====
::

	commit 6bf6224c3dad518329ddc893ef9c5d58dcbabdeb

Build Command
-------------
::

	make -j BUILD_TLS=yes

Run Command
-----------
::

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 48 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6381 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis$ -benchtime=65536x \
	-args true 2147483648 1048576 16 3 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis-3   	
	======================================================================
	server:     4096    warmup:    65536    get:    65536    hit:    31394
	VmHWM:  757092 kB   hit_rate: 47.90%    per_memory_hit_rate: 44.74%
	VmHWM:  744020 kB
	VmHWM:  744440 kB
	193.750s	    output:  703 Mb/s   input:  713 Mb/s
	======================================================================
	   65536	   2956388 ns/op	       151 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	391.281s

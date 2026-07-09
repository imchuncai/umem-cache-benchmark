.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========================
Benchmark-tls-random-2g-8k
==========================

Conclusion
==========
::

	Umem-cache's hit rate is 9% higher than Memcached and 14% higher than Redis.
	Umem-cache's hit throughput is 27% higher than Memcached and 45% higher than Redis.

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

	taskset -c 1,2,3 ./memcached --conn-limit=512 --memory-limit=2048 \
	--max-item-size=1048576 -t 3 \
	--enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem \
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkMemcached$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4972898
	VmHWM: 2128248 kB   hit_rate: 59.28%    per_memory_hit_rate: 58.42%
	507.586s	    output:  212 Mb/s   input:  324 Mb/s
	======================================================================
	 8388608	     60509 ns/op	      9654 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1031.638s

Umem-cache
==========
::

	commit 32cedee7c65bf3af587956f8a80e66efce4643b3

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
	-args true 2147483648 8192 16 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  5364051
	VmHWM: 2105316 kB   hit_rate: 63.94%    per_memory_hit_rate: 63.70%
	435.192s	    output:  218 Mb/s   input:  408 Mb/s
	======================================================================
	 8388608	     51879 ns/op	     12278 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	881.434s

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
	--maxmemory 715827882 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6379 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6380 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

	taskset -c 1,2,3 ./src/redis-server --protected-mode no --appendonly no --save "" \
	--maxmemory 715827882 --maxclients 512 --maxmemory-policy allkeys-lfu \
	--port 0 --tls-port 6381 --tls-cert-file cert.pem \
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

Test Result
-----------
::

	taskset -c 1,2,3 go test -bench=^BenchmarkRedis3$ -benchtime=8388608x \
	-args true 2147483648 8192 16 1 [fe80::179:7fda:ca6e:7c1e%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis3-3   	
	======================================================================
	server:   524288    warmup:  8388608    get:  8388608    hit:  4930576
	VmHWM:  736328 kB   hit_rate: 58.78%    per_memory_hit_rate: 55.86%
	VmHWM:  736136 kB
	VmHWM:  734092 kB
	551.710s	    output:  197 Mb/s   input:  297 Mb/s
	======================================================================
	 8388608	     65769 ns/op	      8494 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	1113.981s

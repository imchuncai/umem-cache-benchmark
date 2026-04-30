.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

============================
BENCHMARK-TLS-RANDOM-100M-1K
============================

CONCLUSION
==========

The test results showed that performance of UMEM-CACHE is the best, is about
28% faster than MEMCACHED and about 99% faster than REDIS.

MEMCACHED
=========
::

	commit 4b9e6198fc44c9eb3ae80802a1b0dcbaf9602969

BUILD COMMAND
-------------
::

	./autogen.sh
	./configure --enable-tls
	make -j

RUN COMMAND
-----------
::

	./memcached --conn-limit=512 --memory-limit=100 --max-item-size=1048576	\
	-t 4 -u root --enable-ssl -o ssl_chain_cert=cert.pem -o ssl_key=key.pem	\
	-o ssl_ca_cert=ca-cert.pem -o ssl_kernel_tls -o ssl_verify_mode=2

TEST RESULT
-----------
::

	go test -bench=^BenchmarkMemcached$ -benchtime=3276800x			\
	-args true 104857600 1024 20 80 80 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkMemcached-4   	
	=======================================================
	case:  819200    hot:  163840(20%)    hot_access: 80% 
	get: 3276800    hit: 1995992    hit_rate: 60.91% 
	hot: 2621031    hit: 1943193    hit_rate: 74.14% 
	VmHWM:	113784 kB    per_memory_hit_rate: 54.82%
	140.254s
	=======================================================
	3276800	     42802 ns/op	      12807 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	283.179s

UMEM-CACHE
==========
::

	commit 53f97eb219364fb18e15431e069b2ceef877b5d9

BUILD COMMAND
-------------
::

	make TLS=1

RUN COMMAND
-----------
::

	./umem-cache 10047 cert.pem key.pem ca-cert.pem

TEST RESULT
-----------
::

	go test -bench=^BenchmarkUmemCache$ -benchtime=3276800x			\
	-args true 104857600 1024 20 80 80 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkUmemCache-4   	
	=======================================================
	case:  819200    hot:  163840(20%)    hot_access: 80% 
	get: 3276800    hit: 1994629    hit_rate: 60.87% 
	hot: 2621031    hit: 1941555    hit_rate: 74.08% 
	VmHWM:	109680 kB    per_memory_hit_rate: 56.83%
	113.924s
	=======================================================
	3276800	     34767 ns/op	      16346 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	233.692s

REDIS
=====
::

	commit e6e0cf5764c99fc1414e46197126e84360536df6

BUILD COMMAND
-------------
::

	make -j BUILD_TLS=yes

RUN COMMAND
-----------
::

	./src/redis-server --protected-mode no --appendonly no --save ""	\
	--maxmemory 52428800 --maxclients 512 --maxmemory-policy allkeys-lru	\
	--port 0 --tls-port 6379 --tls-cert-file cert.pem			\
	--tls-key-file key.pem--tls-ca-cert-file ca-cert.pem

	./src/redis-server --protected-mode no --appendonly no --save ""	\
	--maxmemory 52428800 --maxclients 512 --maxmemory-policy allkeys-lru	\
	--port 0 --tls-port 6380 --tls-cert-file cert.pem			\
	--tls-key-file key.pem --tls-ca-cert-file ca-cert.pem

TEST RESULT
-----------
::

	go test -bench=^BenchmarkRedis2$ -benchtime=3276800x			\
	-args true 104857600 1024 20 80 80 16 1 [fe80::4038:6954:f1a3:4d0f%end0]
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkRedis2-4   	
	=======================================================
	case:  819200    hot:  163840(20%)    hot_access: 80% 
	get: 3276800    hit: 1902724    hit_rate: 58.07% 
	hot: 2621031    hit: 1851111    hit_rate: 70.63% 
	VmHWM:	 67664 kB    per_memory_hit_rate: 43.89%
	VmHWM:	 67808 kB
	174.738s
	=======================================================
	3276800	     53326 ns/op	      8231 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	352.566s

POGOCACHE
=========
::

	commit 71972a9f161d96d91b0f67bfe28897d00bfbd49b

BUILD COMMAND
-------------
::

	make -j NOMIMALLOC=1

RUN COMMAND
-----------
::

	./pogocache --threads=4 --maxmemory=104857600 --maxconns=512 --port=0	\
	-h 192.168.1.70 --tlsport=9401 --tlscert=cert.pem --tlskey=key.pem	\
	--tlscacert=ca-cert.pem

TEST RESULT
-----------
::

	go test -bench=^BenchmarkPogocache$ -benchtime=3276800x			\
	-args true 104857600 1024 20 80 80 16 1 192.168.1.70
	goos: linux
	goarch: arm64
	pkg: github.com/imchuncai/umem-cache-benchmark
	BenchmarkPogocache-4   	
	=======================================================
	case:  819200    hot:  163840(20%)    hot_access: 80% 
	get: 3276800    hit: 1780203    hit_rate: 54.33% 
	hot: 2621031    hit: 1724681    hit_rate: 65.80% 
	VmHWM:	111940 kB    per_memory_hit_rate: 49.70%
	158.158s
	=======================================================
	3276800	     48266 ns/op	      10297 hit/s/mem
	PASS
	ok  	github.com/imchuncai/umem-cache-benchmark	318.728s

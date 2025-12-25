.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
UMEM-CACHE-BENCHMARK
====================

UMEM-CACHE-BENCHMARK is a benchmark project for user space in memory cache.

Multilingual 多语言
==================

- `简体中文 <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/translations/zh_CN/README.rst>`_

RANDOM SIZE BENCHMARK
=====================

We limit server memory to MEM_LIMIT, and test case size is set to 4 times of
that, and the size of key and value is random in the range[0, KV_LIMIT]. We
first get the value from the server, if the get missed, we store it. And 80% of
the time the first 20% of the test cases are used.

FIXED SIZE BENCHMARK
====================

We limit server memory to MEM_LIMIT, and test case size is set to 4 times of
that, and the size of key and value is KV_LIMIT. We first get the value from the
server, if the get missed, we store it. And 80% of the time the first 20% of the
test cases are used.

SUPPORTED APPS
==============

- Memcached
- UmemCache
- Pogocache
- Redis1
- Redis2
- Redis3
- Redis4

Note: because Redis is designed to be single-threaded, it may not be able to
fully utilize the performance of the test machine, so we support it running
multiple instances and distributing keys evenly across these instances.

Note: we use APP's default port

TEST COMMAND
============

RANDOM-2G-1M
------------
::

	make test-{APP}-random-2g-1m REMOTE_IP=[::1]

RANDOM-100M-1K
--------------
::

	make test-{APP}-random-100m-1k REMOTE_IP=[::1]

2G-513K
-------
::

	make test-{APP}-2g-513k REMOTE_IP=[::1]

100M-512B
---------
::

	make test-{APP}-100m-512b REMOTE_IP=[::1]

TEST RESULT
===========

We discovered a serious issue with memcached in our benchmark test, there is a
corner case that your set of a key will never succeed. Specifically, id you
exhausted slab's storage space with chunk data allocated by big keys before
storing any keys into it, you'll unable to store keys that meet the slab's items
size. It can be reproduced by following commands:

::
	# bash A

	./memcached --conn-limit=512 --memory-limit=100 --max-item-size=1048576 -t 4 -u root

::

	# bash B
	
	a_540k=$(for i in {1..552960}; do printf "a"; done)
	a_20000=$(for i in {1..20000}; do printf "a"; done)
	a_30000=$(for i in {1..30000}; do printf "a"; done)

	for i in {1..180}; do
		printf "set ${i} 0 0 552960\r\n${a_540k}\r\n" | nc 127.0.0.1 11211
	done

	for i in {201..344}; do
		printf "set ${i} 0 0 20000\r\n${a_20000}\r\n" | nc 127.0.0.1 11211
	done

	# this set can never be stored
	printf "set 400 0 0 30000\r\n${a_30000}\r\n" | nc 127.0.0.1 11211

RPI4B
-----

Two 4GB version of Raspberry Pi 4 Model B connected in LAN with Gigabit network.
One used as a server and the other as a client. And the installed operating
system is Fedora-Server-40-1.14.aarch64.

Test result is at `rpi4b <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/rpi4b>`_ .

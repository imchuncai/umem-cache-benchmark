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

RPI4B
-----

Two 4GB version of Raspberry Pi 4 Model B connected in LAN with Gigabit network.
One used as a server and the other as a client. And the installed operating
system is Fedora-Server-40-1.14.aarch64.

Test result is at `rpi4b <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/rpi4b>`_ .

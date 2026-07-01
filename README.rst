.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
Umem-cache-benchmark
====================

Umem-cache-benchmark is a benchmark project for user space in memory cache.

It is built to prove that `Umem-cache <https://github.com/imchuncai/umem-cache>`_ is the best cache in the world.

Multilingual 多语言
==================

- `简体中文 <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/translations/zh_CN/README.rst>`_

Benchmark
=========

We will deploy the server and client on two separate machines.

In the following statement, the word *CAP* represents the theoretically maximum
number of key-value pairs that the cache can hold.

Test Dataset
------------

The test dataset was generated using the Zipf variable generator from the Go
standard library, with s=1.0001 and v=1.0.

The size of the test dataset is 1000 times that of *CAP*.

The length of each key in the test dataset is randomized between 16 and 47 bytes,
a range representative of a production environment.

Test Process
------------

Each batch of our tests contains *N* requests, where *N* is 16 times the *CAP*.
For each request, we first get the value from the server, if the get missed,
we store it.

We first make *N* requests to the server to warm up the cache,
and then make *N* more requests and collect statistics.

Supported Apps
==============

- Memcached
- UmemCache
- Redis1
- Redis2
- Redis3
- Redis4

Note: because Redis is designed to be single-threaded, it may not be able to
fully utilize the performance of the test machine, so we support it running
multiple instances and distributing keys evenly across these instances.

Note: we use APP's default port

RPI4B
=====

Two 4GB version of Raspberry Pi 4 Model B with fans connected in LAN with
Gigabit network. One used as a server and the other as a client. And the
installed operating system is Fedora-Server-40-1.14.aarch64.

Test Result
-----------

Umem-cache demonstrates significant advantages across various aspects:

.. [#] Hit rate is 8% to 14% higher than Memcached and 12% to 15% higher than Redis.
.. [#] Without TLS enabled, hit throughput is 2% to 24% higher than Memcached and 7% to 59% higher than Redis.
.. [#] With TLS enabled, hit throughput is 8% to 60% higher than Memcached and 10% to 121% higher than Redis.

The details is at `rpi4b <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/rpi4b>`_ .

Redis Issues
============

In tests with a fixed key-value pair size of 513k bytes, Redis only used about
80% of the memory.

Memcached Issues
================

We discovered a serious issue with Memcached in our benchmark test, there is a
corner case that your set of a key will never succeed. Specifically, if you
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

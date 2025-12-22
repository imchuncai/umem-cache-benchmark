.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
UMEM-CACHE-BENCHMARK
====================

UMEM-CACHE-BENCHMARK是一个为用户空间缓存打造的基准测试项目

Multilingual 多语言
==================

- `简体中文 <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/translations/zh_CN/README.rst>`_

随机大小基准测试
==============

我们把服务器的内存限制为MEM_LIMIT，测试集的大小为该限制的4倍，包含了大小在[0, KV_LIMIT]区间
的随机大小键值对。我们首先向服务端请求对应键的值，如果未请求到，我们将缓存该值。同时我们发出的
80%的请求使用的是前20%的键值对。

固定大小基准测试
==============

我们把服务器的内存限制为MEM_LIMIT，测试集的大小为该限制的4倍，包含了大小为KV_LIMIT的键值对。
我们首先向服务端请求对应键的值，如果未请求到，我们将缓存该值。同时我们发出的 80%的请求使用的是
前20%的键值对。

支持的软件
=========

- Memcached
- UmemCache
- Pogocache
- Redis1
- Redis2
- Redis3
- Redis4

注意： 因为Redis被设计为单线程，可能无法充分利用测试机的性能，所以我们支持了它跑多个实例，然后
把所有的键值对均匀地分散到这些实例上。

注意：我们使用软件的默认端口

测试命令
=======

随机-2G-1M
----------
::

	make test-{APP}-random-2g-1m REMOTE_IP=[::1]

随机-100M-1K
------------
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

测试结果
=======

RPI4B
-----

两台4GB版本的树莓派4 Model B用千兆网络在局域网连接，一台用作服务端，另一台用作客户端。两台机
器所安装的操作系统都为Fedora-Server-40-1.14.aarch64。

测试结果在 `rpi4b <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/translations/zh_CN/rpi4b>`_ 。

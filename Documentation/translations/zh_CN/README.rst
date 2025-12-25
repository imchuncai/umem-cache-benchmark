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

我们在基准测试中发现memcached存在一个严重问题：在特别情况下，你将永远无法成功设置某个键。具体
来说， 如果你在往slab存储任何键之前，已经用大键分配的数据块耗尽了该slab的存储，那么你将无法存
储满足该slab大小的任何键。用以下命令可以复现这个问题：

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

两台4GB版本的树莓派4 Model B用千兆网络在局域网连接，一台用作服务端，另一台用作客户端。两台机
器所安装的操作系统都为Fedora-Server-40-1.14.aarch64。

测试结果在 `rpi4b <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/translations/zh_CN/rpi4b>`_ 。

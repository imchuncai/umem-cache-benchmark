.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====================
Umem-cache-benchmark
====================

Umem-cache-benchmark是一个为用户空间缓存打造的基准测试项目。

它的存在是为了证明 `Umem-cache <https://github.com/imchuncai/umem-cache>`_ 是世界上最好的用户空间缓存。

Multilingual 多语言
==================

- `简体中文 <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/translations/zh_CN/README.rst>`_

基准测试
=======

我们会将服务端和客户端分别部署于两台独立的机器上。

在以下表述中， *CAP* 表示缓存在理论上能容纳的键值对数量。

测试数据集
--------

测试数据集是用Go标准库的Zipf变量生成器生成的，s=1.0001，v=1.0。

测试数据集的大小是 *CAP* 的1000倍。

测试数据集中的每个键的长度在16到47字节之间随机，这个长度范围在生产环境中具有代表性。

测试流程
-------

我们每个批次的测试包含了 *N* 次请求, *N* 的大小为 *CAP* 的16倍。
对于每次请求，我们首先向服务端请求对应键的值，如果未请求到，我们将缓存该值。

我们先向服务端请求 *N* 次用于预热缓存，然后再请求 *N* 次并统计数据。

支持的软件
=========

- Memcached
- UmemCache
- Redis1
- Redis2
- Redis3
- Redis4

注意： 因为Redis被设计为单线程，可能无法充分利用测试机的性能，所以我们支持了它跑多个实例，然后
把所有的键值对均匀地分散到这些实例上。

注意：我们使用软件的默认端口

RPI4B
=====

两台带风扇的4GB版本的树莓派4 Model B用千兆网络在局域网连接，一台用作服务端，另一台用作客户端。两台机
器所安装的操作系统都为Fedora-Server-7.1.3-200.fc44.aarch64。

测试结果
-------

Umem-cache在各个方面都有明显的领先：

.. [#] 命中率比Memcached高9%到14%，比Redis高12%到18%。
.. [#] 在未开启TLS的情况下，命中吞吐量比Memcached高11%到24%，比Redis高12%到62%。
.. [#] 在开启TLS的情况下，命中吞吐量比Memcached高17%到54%，比Redis高19%到81%。

测试详细结果在 `rpi4b <https://github.com/imchuncai/umem-cache-benchmark/tree/master/Documentation/translations/zh_CN/rpi4b>`_ 。

Redis存在的问题
==============

在键值对固定大小为513k字节的测试中，Redis存在内存使用不充分的问题，大约只用了80%的内存。

Memcached存在的问题
==================

我们在基准测试中发现Memcached存在一个严重问题：在特别情况下，你将永远无法成功设置某个键。具体
来说，如果你在往slab存储任何键之前，已经用大键分配的数据块耗尽了该slab的存储，那么你将无法存
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

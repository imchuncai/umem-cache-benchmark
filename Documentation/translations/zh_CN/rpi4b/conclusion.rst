.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

====
结论
====

UMEM-CACHE
==========

优点
----

- 对内存限制十分尊重
- 高命中率
- 快速
- 内置反缓存击穿
- 键可由任意字节组成
- 极少的可配置项

缺点
----

- 键的长度上限是255字节
- 不支持pipeline

MEMCACHED
=========

优点
----

- 尊重内存限制
- 高命中率
- 快速
- 支持pipeline

缺点
----

- 存在一个严重问题，会导致et某些键会一直失败
- 键的长度上限是250字节
- 键只能由可见字符组成

REDIS
=====

优点
----

- 尊重内存限制
- 键可由任意字节组成
- 键的长度上限是512兆字节
- 支持pipeline

缺点
----

- 低命中率
- 慢
- 存在内存使用不充分的问题

POGOCACHE
=========

优点
----

- 支持pipeline

缺点
----

- 对内存限制的不尊重程度令人震惊
- 低命中率
- 慢
- 键的长度上限是250字节
- 键只能由可见字符组成

.. SPDX-License-Identifier: BSD-3-Clause
.. Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

==========
CONCLUSION
==========

UMEM-CACHE
==========

Pros
----

- impressive respects memory limit
- high hit rate
- fast
- built-in anti-dogpiling
- arbitrary key
- few configurable items

Cons
----

- key size limit to 255 bytes
- not support pipeline

MEMCACHED
=========

Pros
----

- respects memory limit
- high hit rate
- fast
- support pipeline

Cons
----

- there is serious issue where set command can never be stored
- key size limit to 250 bytes
- key only takes visible characters

REDIS
=====

Pros
----

- respects memory limit
- arbitrary key
- the maximum allowed key size is 512 MB
- support pipeline

Cons
----

- low hit rate
- slow
- has inefficient use of memory issue

POGOCACHE
=========

Pros
----

- support pipeline

Cons
----

- impressive not respects memory limit
- low hit rate
- slow
- key size limit to 250 bytes
- key only takes visible characters

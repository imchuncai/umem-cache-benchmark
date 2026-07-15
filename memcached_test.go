// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"testing"
	"unsafe"

	"github.com/bradfitz/gomemcache/memcache"
)

const MEMCACHED_PORT = 11211

type MemcachedClient struct {
	client *memcache.Client
}

func (c *MemcachedClient) Init(remoteIPV6 string, threadNR int, config *tls.Config) error {
	client := memcache.New(fmt.Sprintf("%s:%d", remoteIPV6, MEMCACHED_PORT))
	client.Timeout = TIMEOUT
	client.MaxIdleConns = math.MaxInt
	if config != nil {
		client.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			var td tls.Dialer
			td.Config = config
			return td.DialContext(ctx, network, addr)
		}
	}
	c.client = client
	return nil
}

func (c *MemcachedClient) GetOrSet(key []byte, fallbackVal func() []byte) ([]byte, error) {
	strKey := stringKey(key)
	item, err := c.client.Get(strKey)
	if err == nil {
		return item.Value, nil
	}
	if err != memcache.ErrCacheMiss {
		return nil, err
	}

	val := fallbackVal()
	c.client.Set(&memcache.Item{Key: strKey, Value: val})
	// Memcached server may return error: "SERVER_ERROR out of memory storing object",
	// and we just ignored it instead of retry,
	// because Memcached has a serious issue may cause the retry endless
	return val, nil
}

func stringKey(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func BenchmarkMemcached(b *testing.B) {
	parallel[MemcachedClient](b)
}

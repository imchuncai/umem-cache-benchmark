// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"crypto/tls"
	"fmt"
	"testing"

	client "github.com/imchuncai/umem-cache-client-Go"
)

const UMEM_CACHE_PORT = 10047

type UmemCacheClient struct {
	client *client.Client
}

func (c *UmemCacheClient) Init(remoteIPV6 string, threadNR int, config *tls.Config) error {
	var err error
	c.client, err = client.New(
		fmt.Sprintf("%s:%d", remoteIPV6, UMEM_CACHE_PORT),
		client.Config{TIMEOUT, threadNR, 0, config},
	)
	if err != nil {
		return fmt.Errorf("new client failed: %w", err)
	}
	return nil
}

func (c *UmemCacheClient) GetOrSet(key []byte, i uint64, fallbackVal func() []byte) ([]byte, error) {
	get := func([]byte) ([]byte, error) {
		return fallbackVal(), nil
	}
	return c.client.GetOrSet(key, get)
}

func BenchmarkUmemCache(b *testing.B) {
	parallel[UmemCacheClient](b)
}

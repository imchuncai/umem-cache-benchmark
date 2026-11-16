// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"crypto/tls"
	"fmt"
	"testing"

	client "github.com/imchuncai/umem-cache-client-Go"
)

type UmemCacheCluster struct {
	client *client.ClusterClient
}

func (c *UmemCacheCluster) Init(remoteIPV6 string, threadNR int, config *tls.Config) error {
	var err error
	c.client, err = client.NewCluster(
		[]string{fmt.Sprintf("%s:%d", remoteIPV6, UMEM_CACHE_PORT)},
		client.Config{TIMEOUT, 0, config},
	)
	if err != nil {
		return fmt.Errorf("new client failed: %w", err)
	}
	return nil
}

func (c *UmemCacheCluster) GetOrSet(key []byte, fallbackVal func() []byte) ([]byte, error) {
	get := func([]byte) ([]byte, error) {
		return fallbackVal(), nil
	}
	return c.client.GetOrSet(key, get)
}

func BenchmarkUmemCacheCluster(b *testing.B) {
	parallel[UmemCacheCluster](b)
}

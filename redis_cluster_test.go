// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

type RedisCluster struct {
	client *redis.ClusterClient
}

func (c *RedisCluster) Init(remoteIPV6 string, threadNR int, config *tls.Config) error {
	c.client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:     []string{fmt.Sprintf("%s:%d", remoteIPV6, REDIS_PORT)},
		Password:  "",
		TLSConfig: config,
	})
	return nil
}

func (c *RedisCluster) GetOrSet(key []byte, fallbackVal func() []byte) ([]byte, error) {
	strKey := stringKey(key)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	val, err := c.client.Get(ctx, strKey).Bytes()
	if err == redis.Nil {
		val = fallbackVal()
		return val, c.client.Set(ctx, strKey, val, 0).Err()
	}
	return val, err
}

func BenchmarkRedisCluster(b *testing.B) {
	parallel[RedisCluster](b)
}

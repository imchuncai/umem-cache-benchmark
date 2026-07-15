// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/bits"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/twmb/murmur3"
)

const REDIS_PORT = 6379

type RedisClient struct {
	clients []*redis.Client
}

func (c *RedisClient) Init(remoteIPV6 string, threadNR int, config *tls.Config) error {
	c.clients = make([]*redis.Client, threadNR)
	for i := range c.clients {
		port := REDIS_PORT + i
		c.clients[i] = redis.NewClient(&redis.Options{
			Addr:      fmt.Sprintf("%s:%d", remoteIPV6, port),
			Password:  "",
			DB:        0,
			TLSConfig: config,
		})
	}
	return nil
}

func (c *RedisClient) GetOrSet(key []byte, fallbackVal func() []byte) ([]byte, error) {
	h, _ := murmur3.SeedSum128(REDIS_PORT, REDIS_PORT, key)
	hi, _ := bits.Mul32(uint32(h), uint32(len(c.clients)))
	client := c.clients[hi]
	strKey := stringKey(key)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	val, err := client.Get(ctx, strKey).Bytes()
	if err == redis.Nil {
		val = fallbackVal()
		return val, client.Set(ctx, strKey, val, 0).Err()
	}
	return val, err
}

func BenchmarkRedis(b *testing.B) {
	parallel[RedisClient](b)
}

// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

const REDIS_PORT = 6379

type RedisClient struct {
	client *redis.Client
}

func (c *RedisClient) Init(remoteIPV6 string, config *tls.Config) error {
	c.client = redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%d", remoteIPV6, REDIS_PORT),
		Password:  "",
		DB:        0,
		TLSConfig: config,
	})
	return nil
}

func (c *RedisClient) GetOrSet(key []byte, i uint64, fallbackVal func() []byte) ([]byte, error) {
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

type RedisClientN struct {
	clients []*redis.Client
}

func (c *RedisClientN) Init(remoteIPV6 string, config *tls.Config) error {
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

func (c *RedisClientN) GetOrSet(key []byte, i uint64, fallbackVal func() []byte) ([]byte, error) {
	client := c.clients[int(i)%len(c.clients)]
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

func BenchmarkRedis1(b *testing.B) {
	parallel(b, &RedisClient{})
}
func BenchmarkRedis2(b *testing.B) {
	parallel(b, &RedisClientN{make([]*redis.Client, 2)})
}
func BenchmarkRedis3(b *testing.B) {
	parallel(b, &RedisClientN{make([]*redis.Client, 3)})
}
func BenchmarkRedis4(b *testing.B) {
	parallel(b, &RedisClientN{make([]*redis.Client, 4)})
}

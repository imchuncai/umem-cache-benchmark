// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/redis/go-redis/v9"
)

const POGOCACHE_PORT = 9401

func runPogocacheServer(serverMemory int, kvSizeLimit int, remoteIP string) ([]*exec.Cmd, GetOrSetFunc, error) {
	var cmd []*exec.Cmd
	if remoteIP == "" {
		cmd = []*exec.Cmd{exec.Command("./pogocache/pogocache",
			fmt.Sprintf("--threads=%d", THREAD_NR),
			fmt.Sprintf("--maxmemory=%d", serverMemory),
			"--maxconns=512",
			fmt.Sprintf("--port=%d", POGOCACHE_PORT),
		)}
		remoteIP = "127.0.0.1"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", remoteIP, POGOCACHE_PORT),
		Password: "",
		DB:       0,
	})

	getOrSet := func(key []byte, i int, fallbackVal func() []byte) ([]byte, error) {
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
	return cmd, getOrSet, nil
}

func BenchmarkPogocache(b *testing.B) {
	parallel(b, runPogocacheServer)
}

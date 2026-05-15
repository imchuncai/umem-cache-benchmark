// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os/exec"
	"testing"

	"github.com/redis/go-redis/v9"
)

const REDIS_PORT = 6379

func __runRedisServerCmd(serverMemory int, port int, tlsEnabled bool) *exec.Cmd {
	args := []string{
		"--protected-mode no",
		"--appendonly no",
		`--save ""`,
		fmt.Sprintf("--maxmemory %d", serverMemory),
		"--maxclients 512",
		"--maxmemory-policy allkeys-lfu",
	}
	if tlsEnabled {
		args = append(args,
			"--port 0",
			fmt.Sprintf("--tls-port %d", port),
			"--tls-cert-file cert.pem",
			"--tls-key-file key.pem",
			"--tls-ca-cert-file ca-cert.pem",
		)
	} else {
		args = append(args, fmt.Sprintf("--port %d", port))
	}
	return exec.Command("./redis/src/redis-server", args...)
}

func runRedisServerN(n int) RunServer {
	return func(serverMemory int, kvSizeLimit int, tlsEnabled bool, remoteIP string) ([]*exec.Cmd, GetOrSetFunc, error) {
		var cmds []*exec.Cmd
		if remoteIP == "" {
			for i := range n {
				port := REDIS_PORT + i
				cmd := __runRedisServerCmd(serverMemory/n, port, tlsEnabled)
				cmds = append(cmds, cmd)
			}
			remoteIP = "[::1]"
		}

		var tlsConfig *tls.Config
		if tlsEnabled {
			tlsConfig = TLS_CONFIG
		}
		var clients []*redis.Client
		for i := range n {
			port := REDIS_PORT + i
			clients = append(clients, redis.NewClient(&redis.Options{
				Addr:      fmt.Sprintf("%s:%d", remoteIP, port),
				Password:  "",
				DB:        0,
				TLSConfig: tlsConfig,
			}))
		}

		getOrSet := func(key []byte, i int, fallbackVal func() []byte) ([]byte, error) {
			client := clients[i%len(clients)]
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
		return cmds, getOrSet, nil
	}
}

func BenchmarkRedis1(b *testing.B) {
	parallel(b, runRedisServerN(1))
}
func BenchmarkRedis2(b *testing.B) {
	parallel(b, runRedisServerN(2))
}
func BenchmarkRedis3(b *testing.B) {
	parallel(b, runRedisServerN(3))
}
func BenchmarkRedis4(b *testing.B) {
	parallel(b, runRedisServerN(4))
}

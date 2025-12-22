// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	client "github.com/imchuncai/umem-cache-client-Go"
)

const UMEM_CACHE_PORT = 10047

func runUmemCacheServer(serverMemory int, kvSizeLimit int, remoteIP string) ([]*exec.Cmd, GetOrSetFunc, error) {
	var cmd []*exec.Cmd
	if remoteIP == "" {
		cmd = []*exec.Cmd{exec.Command("umem-cache/umem-cache", strconv.Itoa(UMEM_CACHE_PORT))}
		remoteIP = "[::1]"
	}
	client, err := client.New(fmt.Sprintf("%s:%d", remoteIP, UMEM_CACHE_PORT), client.Config{TIMEOUT, 4, 0, nil})
	if err != nil {
		return nil, nil, fmt.Errorf("new client failed: %w", err)
	}

	getOrSet := func(key []byte, i int, fallbackVal func() []byte) ([]byte, error) {
		fallbackGet := func([]byte) ([]byte, error) {
			return fallbackVal(), nil
		}
		return client.GetOrSet(key, fallbackGet)
	}
	return cmd, getOrSet, nil
}

func BenchmarkUmemCache(b *testing.B) {
	parallel(b, runUmemCacheServer)
}

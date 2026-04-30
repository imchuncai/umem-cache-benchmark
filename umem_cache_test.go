// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"crypto/tls"
	"fmt"
	"os/exec"
	"strconv"
	"testing"

	client "github.com/imchuncai/umem-cache-client-Go"
)

const UMEM_CACHE_PORT = 10047

func runUmemCacheServer(serverMemory int, kvSizeLimit int, tlsEnabled bool, remoteIP string) ([]*exec.Cmd, GetOrSetFunc, error) {
	var cmd []*exec.Cmd
	if remoteIP == "" {
		args := []string{strconv.Itoa(UMEM_CACHE_PORT)}
		if tlsEnabled {
			args = append(args, "cert.pem", "key.pem", "ca-cert.pem")
		}
		cmd = []*exec.Cmd{exec.Command("umem-cache/umem-cache", args...)}
		remoteIP = "[::1]"
	}

	var tlsConfig *tls.Config
	if tlsEnabled {
		tlsConfig = TLS_CONFIG
	}
	client, err := client.New(fmt.Sprintf("%s:%d", remoteIP, UMEM_CACHE_PORT), client.Config{TIMEOUT, 4, 0, tlsConfig})
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

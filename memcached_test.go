// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package main

import (
	"fmt"
	"math"
	"os/exec"
	"testing"
	"unsafe"

	"github.com/bradfitz/gomemcache/memcache"
)

const MEMCACHED_PORT = 11211

func runMemcachedServer(serverMemory int, kvSizeLimit int, remoteIP string) ([]*exec.Cmd, GetOrSetFunc, error) {
	var cmd []*exec.Cmd
	if remoteIP == "" {
		itemLimit := kvSizeLimit * 2 >> 20 << 20
		if itemLimit < 1<<20 {
			// can't set max-item-size below that
			itemLimit = 1 << 20
		}
		cmd = []*exec.Cmd{exec.Command("./memcached/memcached",
			fmt.Sprintf("-t %d", THREAD_NR),
			"--conn-limit=512",
			fmt.Sprintf("--memory-limit=%d", serverMemory>>20),
			fmt.Sprintf("--max-item-size=%d", itemLimit),
		)}
		remoteIP = "[::1]"
	}
	client := memcache.New(fmt.Sprintf("%s:%d", remoteIP, MEMCACHED_PORT))
	client.Timeout = TIMEOUT
	client.MaxIdleConns = math.MaxInt

	getOrSet := func(key []byte, i int, fallbackVal func() []byte) ([]byte, error) {
		strKey := stringKey(key)
		item, err := client.Get(strKey)
		if err == nil {
			return item.Value, nil
		}
		if err != memcache.ErrCacheMiss {
			return nil, err
		}

		val := fallbackVal()
		client.Set(&memcache.Item{Key: strKey, Value: val})
		// memcached server may return error: "SERVER_ERROR out of memory storing object",
		// and we just ignored it instead of retry,
		// because memcached has a serious issue may cause the retry endless
		return val, nil
	}
	return cmd, getOrSet, nil
}

func stringKey(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func BenchmarkMemcached(b *testing.B) {
	parallel(b, runMemcachedServer)
}

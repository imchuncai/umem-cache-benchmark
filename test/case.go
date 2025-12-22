// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package test

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand/v2"
)

type Case struct {
	I   int
	Hot bool
	Key []byte
	Val []byte
}

func randV(r *rand.Rand, pool []byte, n int) []byte {
	if n > len(pool) {
		panic(1)
	}

	i := r.IntN(len(pool) - n + 1)
	return pool[i : i+n]
}

func NewCase(r *rand.Rand, i int, kvSizeLimit int, kvRandSize bool, valPool []byte) Case {
	var kvSize int
	if kvRandSize {
		kvSize = r.IntN(kvSizeLimit-KEY_PREFIX_SIZE+1) + KEY_PREFIX_SIZE
	} else {
		kvSize = kvSizeLimit
	}

	keySizeLimit := min(KEY_SIZE_LIMIT, kvSize)
	keySize := r.IntN(keySizeLimit-KEY_PREFIX_SIZE+1) + KEY_PREFIX_SIZE
	valSize := kvSize - keySize

	prefix := make([]byte, 8)
	binary.BigEndian.PutUint64(prefix, uint64(i))
	key := make([]byte, keySize)
	base64.StdEncoding.Encode(key, prefix)
	copy(key[12:], randV(r, valPool, keySize)[12:])
	return Case{
		I:   i,
		Key: key,
		Val: randV(r, valPool, valSize),
	}
}

func (tc Case) Size() int {
	return len(tc.Key) + len(tc.Val)
}

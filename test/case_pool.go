// SPDX-License-Identifier: BSD-3-Clause
// Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

package test

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

const (
	SEED            = uint64(47)
	KEY_PREFIX_SIZE = 12  // base64.StdEncoding.EncodedLen(8)
	KEY_SIZE_LIMIT  = 250 // limit by MEMCACHED
)

type Pool struct {
	mu sync.Mutex
	r  *rand.Rand

	hotCaseAccessPercent int
	hotN                 int
	Cases                []Case
}

func NewPool(kvSizeLimit int, kvRandSize bool, serverMemory int, hotCasePercent int, hotCaseAccessPercent int, hotCaseServerMemoryPercent int) *Pool {
	if kvSizeLimit < KEY_PREFIX_SIZE {
		panic(fmt.Sprintf("kvSizeLimit should be at least %d", KEY_PREFIX_SIZE))
	}

	kvAverageSize := kvSizeLimit
	if kvRandSize {
		kvAverageSize /= 2
	}

	// poolKVSize * HOT_CASE_PERCENT / 100 == serverMemory * HOT_CASE_SERVER_MEMORY_PERCENT / 100
	poolKVSize := serverMemory * hotCaseServerMemoryPercent / hotCasePercent
	caseN := poolKVSize / kvAverageSize
	hotCaseN := caseN * hotCasePercent / 100

	p := Pool{
		r:                    rand.New(rand.NewPCG(SEED, SEED)),
		hotCaseAccessPercent: hotCaseAccessPercent,
		hotN:                 hotCaseN,
		Cases:                make([]Case, caseN),
	}
	kv := make([]byte, kvSizeLimit*2)
	for i := range kv {
		// MEMCACHED requires keys to be visible characters
		kv[i] = byte(p.r.IntN(127-33) + 33)
	}
	for i := 0; i < hotCaseN; i++ {
		p.Cases[i] = NewCase(p.r, i, kvSizeLimit, kvRandSize, kv)
		p.Cases[i].Hot = true
	}
	for i := hotCaseN; i < caseN; i++ {
		p.Cases[i] = NewCase(p.r, i, kvSizeLimit, kvRandSize, kv)
		p.Cases[i].Hot = false
	}
	return &p
}

func (p *Pool) randHot() Case {
	i := p.r.IntN(p.hotN)
	return p.Cases[i]
}

func (p *Pool) randCold() Case {
	i := p.r.IntN(len(p.Cases)-p.hotN) + p.hotN
	return p.Cases[i]
}

func (p *Pool) RandCase() Case {
	p.mu.Lock()
	defer p.mu.Unlock()

	i := p.r.IntN(100)
	if i < p.hotCaseAccessPercent {
		return p.randHot()
	}
	return p.randCold()
}

func (p *Pool) HotMegabytes() int {
	n := 0
	for _, tc := range p.Cases[:p.hotN] {
		n += tc.Size()
	}
	return n >> 20
}

func (p *Pool) ColdMegaBytes() int {
	n := 0
	for _, tc := range p.Cases[p.hotN:] {
		n += tc.Size()
	}
	return n >> 20
}

func (p *Pool) CaseMegabytes() int {
	n := 0
	for _, tc := range p.Cases {
		n += tc.Size()
	}
	return n >> 20
}

func (p *Pool) HotN() int {
	return p.hotN
}

func (p *Pool) ColdN() int {
	return len(p.Cases) - p.hotN
}

func (p *Pool) CaseN() int {
	return len(p.Cases)
}

func (p *Pool) Hot() []Case {
	return p.Cases[:p.hotN]
}

func (p *Pool) Cold() []Case {
	return p.Cases[p.hotN:]
}

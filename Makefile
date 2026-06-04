# SPDX-License-Identifier: BSD-3-Clause
# Copyright (C) 2025-2026, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

TEST_N_RATIO = 16
PARALLELISM = 16

ifndef APP
APPS = Memcached UmemCache Redis1 Redis2 Redis3 Redis4
$(error APP is not defined! Available apps: $(APPS))
endif

ifndef REMOTE_IPV6
$(error REMOTE_IPV6 is not defined! Example: [::1])
endif

ifdef TLS
	ifneq ($(TLS),0)
		TLS = 1
	endif
else
	TLS = 0
endif

math = $(shell echo "$$(( $(1) ))" )

define test =
	$(eval BENCHTIME := $(shell					\
		if [[ $(1) = "true" ]]; then				\
			echo $$((2 * $(TEST_N_RATIO) * $(2) / $(3)));	\
	  	else							\
			echo $$(($(TEST_N_RATIO) * $(2) / $(3)));	\
		fi))

	go test -bench=^Benchmark$(APP)$$ -benchtime=$(BENCHTIME)x	\
	-args $(1) $(2) $(3) $(PARALLELISM) $(TLS) $(REMOTE_IPV6)
	@echo ""
endef

test-random-100m-1k:
	$(call test,true,$(call math, 100 << 20),$(call math, 1 << 10))

test-random-1g-1k:
	$(call test,true,$(call math, 1 << 30),$(call math, 1 << 10))

test-random-2g-8k:
	$(call test,true,$(call math, 2 << 30),$(call math, 8 << 10))

test-random-2g-1m:
	$(call test,true,$(call math, 2 << 30),$(call math, 1 << 20))

test-1g-512b:
	$(call test,false,$(call math, 1 << 30),512)

test-2g-513k:
	$(call test,false,$(call math, 2 << 30),$(call math, 513 << 10))

# SPDX-License-Identifier: BSD-3-Clause
# Copyright (C) 2025, Shu De Zheng <imchuncai@gmail.com>. All Rights Reserved.

apps = Memcached UmemCache Pogocache Redis1 Redis2 Redis3 Redis4

# 80% of the cache accesses comes from 20% of the keys 
HOT_CASE_PERCENT = 20
HOT_CASE_ACCESS_PERCENT = 80

# Interestingly, if HOT_CASE_SERVER_MEMORY_PERCENT is set to 60, the cache hit rate becomes 80%
HOT_CASE_SERVER_MEMORY_PERCENT = 80
PARALLELISM = 16
# REMOTE_IP = [::1]

math = $(shell echo "$$(( $(1) ))" )

define test =
	$(eval BENCHTIME := $(shell										\
		if [[ $(2) = "true" ]]; then									\
			echo $$((8 * $(3) * $(HOT_CASE_SERVER_MEMORY_PERCENT) / $(HOT_CASE_PERCENT) / $(4)));	\
	  	else												\
			echo $$((4 * $(3) * $(HOT_CASE_SERVER_MEMORY_PERCENT) / $(HOT_CASE_PERCENT) / $(4)));	\
		fi))

	@if [[ $(1) = "UmemCache" ]]; then				\
		cd umem-cache && make -s MEM_LIMIT=$(3);		\
	fi

	go test -bench=^Benchmark$(1)$$ -benchtime=$(BENCHTIME)x  \
	-args $(2) $(3) $(4) $(HOT_CASE_PERCENT) $(HOT_CASE_ACCESS_PERCENT) $(HOT_CASE_SERVER_MEMORY_PERCENT) $(PARALLELISM) $(REMOTE_IP)
	@echo ""
endef

test-random-100m-1k test-random-2g-1m test-2g-513k test-100m-512b: test-%: $(foreach app, $(apps), test-$(app)-%);

test-%-random-100m-1k:
	$(call test,$(*),true,$(call math, 100 << 20),$(call math, 1 << 10))

test-%-random-2g-1m:
	$(call test,$(*),true,$(call math, 2 << 30),$(call math, 1 << 20))

test-%-2g-513k:
	$(call test,$(*),false,$(call math, 2 << 30),$(call math, 513 << 10))

test-%-100m-512b:
	$(call test,$(*),false,$(call math, 100 << 20),512)

update:
	git submodule update --init --recursive --remote
	cd memcached && ./autogen.sh && ./configure
	$(MAKE) -j -C memcached
	$(MAKE) -j -C redis
	$(MAKE) -j -C pogocache

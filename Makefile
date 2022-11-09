LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
	TARGET_ARCH_LOCAL=amd64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),armv8)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 4),armv)
	TARGET_ARCH_LOCAL=arm
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 5),arm64)
	TARGET_ARCH_LOCAL=arm64
else ifeq ($(shell echo $(LOCAL_ARCH) | head -c 7),aarch64)
	TARGET_ARCH_LOCAL=arm64
else
	TARGET_ARCH_LOCAL=amd64
endif

LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   TARGET_OS_LOCAL = linux
else ifeq ($(LOCAL_OS),Darwin)
   TARGET_OS_LOCAL = darwin
else
   TARGET_OS_LOCAL = windows
endif

K8S_CLUSTER ?= kind
KO_DOCKER_REPO ?= localhost:5001
#### TARGET setup-tools required tools based on OS ####
setup-tools:
	./scripts/tools/setup_$(TARGET_OS_LOCAL)_$(TARGET_ARCH_LOCAL).sh
#### TARGET setup-k8s ####
setup-k8s:
	./infra/$(K8S_CLUSTER)/setup.sh

define genRunTestApp
.PHONY: run-test-$(1)
run-test-$(1):
	vcluster create $(1) --expose
endef

TESTS := $(shell ls tests)
# Generate run test targets
$(foreach ITEM,$(TESTS),$(eval $(call genRunTestApp,$(ITEM))))
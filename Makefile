all::

.DEFAULT_GOAL := help

##@ General

.PHONY: help
help: ## Show this help message.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

include Makefile.common

BIN ?= sdtd_exporter

PROMTOOL_VERSION ?= 2.51.0
PROMTOOL_URL     ?= https://github.com/prometheus/prometheus/releases/download/v$(PROMTOOL_VERSION)/prometheus-$(PROMTOOL_VERSION).$(GO_BUILD_PLATFORM).tar.gz
PROMTOOL         ?= $(FIRST_GOPATH)/bin/promtool

DOCKER_IMAGE_NAME       ?= $(BIN)
MACH                    ?= $(shell uname -m)

ifeq($(MACH),x86_64)
ARCH := amd64
else
ifeq($(MACH),aarch64)
ARCH := arm64
endif
endif

STATICCHECK_IGNORE =

PROMU_CONF := .promu.yml
PROMU := $(FIRST_GOPATH)/bin/promu --config $(PROMU_CONF)

.PHONY: build
build: promu $(BIN)
$(BIN): *.go
	$(PROMU) build --prefix=output

fmt:
	@echo ">> Running fmt"
	gofmt -l -w -s .

crossbuild: promu
	@echo ">> Running crossbuild"
	GOARCH=amd64 $(PROMU) build --prefix=output/amd64
	GOARCH=arm64 $(PROMU) build --prefix=output/arm64

clean:
	@echo ">> Running clean"
	rm -rf $(BIN) output

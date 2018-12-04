GO_VERBOSE := -v
CS_BUILD_DIR := $(PWD)/build/container-suseconnect

ifneq "$(VERBOSE)" "1"
GO_VERBOSE=
.SILENT:
endif

.PHONY: test

all: deps
	rm -rf $(CS_BUILD_DIR)
	mkdir -p $(CS_BUILD_DIR)
	GOBIN=$(CS_BUILD_DIR) go install -a $(GO_VERBOSE) ./...

test: deps
	go test $(GO_VERBOSE) ./...
	build/ci/climate -t 80 .

deps:
	go mod vendor

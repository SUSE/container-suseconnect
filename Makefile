GO_VERBOSE := -v
CS_BUILD_DIR := $(PWD)/build/container-suseconnect

export GO111MODULE=auto

ifneq "$(VERBOSE)" "1"
GO_VERBOSE=
.SILENT:
endif

.PHONY: test

all:
	rm -rf $(CS_BUILD_DIR)
	mkdir -p $(CS_BUILD_DIR)
	GOBIN=$(CS_BUILD_DIR) go install -ldflags='-w -s' -a $(GO_VERBOSE) ./...

test:
	go test $(GO_VERBOSE) ./...
	build/ci/climate -t 80 .

mod:
	export GO111MODULE=on \
		go mod tidy && \
		go mod vendor && \
		go mod verify

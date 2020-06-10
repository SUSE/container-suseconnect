GO_VERBOSE := -v
CS_BUILD_DIR := $(PWD)/build/container-suseconnect

export GO111MODULE=auto

ifneq "$(VERBOSE)" "1"
GO_VERBOSE=
.SILENT:
endif

all:
	rm -rf $(CS_BUILD_DIR)
	mkdir -p $(CS_BUILD_DIR)
	GOBIN=$(CS_BUILD_DIR) go install -ldflags='-w -s' -a $(GO_VERBOSE) ./...

.PHONY: test
test: test-unit validate-go

.PHONY: test-unit
test-unit:
	go test $(GO_VERBOSE) ./...


.PHONY: validate-go
validate-go:
	build/ci/climate -t 80 internal
	build/ci/climate -t 80 internal/regionsrv

	@which gofmt >/dev/null 2>/dev/null || (echo "ERROR: gofmt not found." && false)
	test -z "$$(gofmt -s -l . | grep -vE '^vendor/' | tee /dev/stderr)"

mod:
	export GO111MODULE=on \
		go mod tidy && \
		go mod vendor && \
		go mod verify

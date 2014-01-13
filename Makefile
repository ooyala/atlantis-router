PROJECT_ROOT := $(shell pwd)
ifeq ($(shell pwd | xargs dirname | xargs basename),lib)
	VENDOR_PATH := $(shell pwd | xargs dirname | xargs dirname)/vendor
else
	VENDOR_PATH := $(PROJECT_ROOT)/vendor
endif

GOPATH := $(PROJECT_ROOT):$(VENDOR_PATH)
export GOPATH

all:
	@echo "make fmt|install-deps|test|annotate|example|routertest|clean"

install-deps:
	@echo "Installing Dependencies..."
	@rm -rf $(VENDOR_PATH)
	@mkdir -p $(VENDOR_PATH) || exit 2
	@GOPATH=$(VENDOR_PATH) go get launchpad.net/gozk
	@GOPATH=$(VENDOR_PATH) go get code.google.com/p/go.tools/cmd/cover
	@echo "Done."

test:
ifdef TEST_PACKAGE
	@echo "Testing $$TEST_PACKAGE..."
	@go test $$TEST_PACKAGE $$VERBOSE $$EXTRA_FLAGS
else
	@for p in `find ./src -type f -name "*.go" |sed 's-\./src/\(.*\)/.*-\1-' |sort -u`; do \
		echo "Testing $$p..."; \
		go test $$p -cover || exit 1; \
	done
	@echo
	@echo "ok."
endif

annotate:
ifdef TEST_PACKAGE
	@echo "Annotating $$TEST_PACKAGE..."
	@go test $$TEST_PACKAGE $$VERBOSE $$EXTRA_FLAGS -coverprofile=cover.out
	@go tool cover -html=cover.out
	@rm -f cover.out
else
	@echo "Specify package!"
endif

.PHONY: example
example:
	@go build -o example/router example/router.go

.PHONY: routertest
routertest:
	@go build -o bm/routertest/routertest bm/routertest/routertest.go

clean:
	@rm -f bm/routertest/routertest example/router
fmt:
	@find src -name \*.go -exec gofmt -l -w {} \;
	@find example -name \*.go -exec gofmt -l -w {} \;
	@find bm -name \*.go -exec gofmt -l -w {} \;

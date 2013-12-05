PROJECT_ROOT := $(shell pwd)
ifeq ($(shell pwd | xargs dirname | xargs basename),"lib")
	VENDOR_PATH := $(PROJECT_ROOT)/vendor
else
	VENDOR_PATH := $(shell pwd | xargs dirname | xargs dirname)/vendor
endif

GOPATH := $(PROJECT_ROOT):$(VENDOR_PATH)
export GOPATH

GOCOV := $(VENDOR_PATH)/bin/gocov
export GOCOV

all:
	@echo "make fmt|install-deps|test|coverage|annotate|example|routertest|clean"

install-deps:
	@echo "Installing Dependencies..."
	@rm -rf $(VENDOR_PATH)
	@mkdir -p $(VENDOR_PATH) || exit 2
	@GOPATH=$(VENDOR_PATH) go get github.com/axw/gocov/gocov
	@GOPATH=$(VENDOR_PATH) go get launchpad.net/gozk
	@echo "Done."

test:
ifdef TEST_PACKAGE
	@echo "Testing $$TEST_PACKAGE..."
	@go test $$TEST_PACKAGE $$VERBOSE $$RACE
else
	@for p in `find ./src -type f -name "*.go" |sed 's-\./src/\(.*\)/.*-\1-' |sort -u`; do \
		echo "Testing $$p..."; \
		go test $$p || exit 1; \
	done
	@echo
	@echo "ok."
endif

coverage:
ifdef TEST_PACKAGE
	@echo "Coveraging $$TEST_PACKAGE..."
	@$(GOCOV) test $$TEST_PACKAGE | $(GOCOV) report
else
	@for p in `find ./src -type f -name "*.go" |sed 's-\./src/\(.*\)/.*-\1-' |sort -u`; do \
		echo "Coveraging $$p..."; \
		$(GOCOV) test $$p >coverage.json; \
		if [[ "$$?" != "0" ]]; then \
			rm -f coverage.json; \
			continue; \
		fi; \
		$(GOCOV) report <coverage.json; \
		rm -f coverage.json; \
	done
endif

annotate:
ifdef TEST_PACKAGE
	@echo "Annotating $$TEST_PACKAGE..."
	@$(GOCOV) test $$TEST_PACKAGE >coverage.json
	@$(GOCOV) annotate coverage.json
	@rm -f coverage.json
else
	@for p in `find ./src -type f -name "*.go" |sed 's-\./src/\(.*\)/.*-\1-' |sort -u`; do \
		echo "Annotating $$p..."; \
		$(GOCOV) test $$p >coverage.json; \
		if [[ "$$?" != "0" ]]; then \
			rm -f coverage.json; \
			continue; \
		fi; \
		$(GOCOV) annotate coverage.json; \
		rm -f coverage.json; \
	done
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

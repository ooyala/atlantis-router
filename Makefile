## Copyright 2014 Ooyala, Inc. All rights reserved.
##
## This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
## except in compliance with the License. You may obtain a copy of the License at
## http://www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software distributed under the License is
## distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
## See the License for the specific language governing permissions and limitations under the License.

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
	@GOPATH=$(VENDOR_PATH) go get github.com/gorilla/mux
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

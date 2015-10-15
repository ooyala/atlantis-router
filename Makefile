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
PROJECT_NAME := $(shell pwd | xargs basename)

DEB_STAGING := $(PROJECT_ROOT)/staging
PKG_INSTALL_DIR := $(DEB_STAGING)/$(PROJECT_NAME)/opt/atlantis-router

ifndef VERSION
	VERSION := "0.1.0"
endif

GOPATH := $(PROJECT_ROOT):$(VENDOR_PATH)
export GOPATH

GOM := $(VENDOR_PATH)/bin/gom
GOM_VENDOR_NAME := vendor
export GOM_VENDOR_NAME

all:
	@echo "make fmt|install-deps|test|annotate|example|routertest|clean"

init: clean
	@mkdir bin

build: init install-deps example

install-deps:
	@echo "Installing Dependencies..."
	@rm -rf $(VENDOR_PATH)
	@mkdir -p $(VENDOR_PATH) || exit 2
	@GOPATH=$(VENDOR_PATH) go get github.com/mattn/gom
	$(GOM) install
	@echo "Done."

deb: clean build example
	@cp -a $(PROJECT_ROOT)/deb $(DEB_STAGING)
	@cp -a $(PROJECT_ROOT)/bin $(PKG_INSTALL_DIR)/
	@sed -ri "s/__VERSION__/$(VERSION)/" $(DEB_STAGING)/$(PROJECT_NAME)/DEBIAN/control
	@sed -ri "s/__PACKAGE__/atlantis-router/" $(DEB_STAGING)/$(PROJECT_NAME)/DEBIAN/control
	@cd $(DEB_STAGING) && dpkg -b $(PROJECT_NAME) $(PROJECT_ROOT)

test: clean install-deps
ifdef TEST_PACKAGE
	@echo "Testing $$TEST_PACKAGE..."
	@go test $$TEST_PACKAGE $$VERBOSE $$EXTRA_FLAGS
else
	@for p in `find ./src -type f -name "*_test.go" |sed 's-\./src/\(.*\)/.*-\1-' |sort -u`; do \
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
	@go build -o bin/router example/router.go

.PHONY: routertest
routertest:
	@go build -o bm/routertest/routertest bm/routertest/routertest.go

clean:
	@rm -rf bm/routertest/routertest bin $(DEB_STAGING) atlantis-router_*.deb
fmt:
	@find src -name \*.go -exec gofmt -l -w {} \;
	@find example -name \*.go -exec gofmt -l -w {} \;
	@find bm -name \*.go -exec gofmt -l -w {} \;

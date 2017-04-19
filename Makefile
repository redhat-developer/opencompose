#!/usr/bin/make -f
.PHONY: all
all: build

SHELL :=/bin/bash

GIT_VERSION =$(shell git describe --tags --abbrev=0 2>/dev/null)
GIT_EXACT_VERSION =$(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_COMMIT =$(shell git rev-parse --short HEAD 2>/dev/null)
GIT_TREE_STATE =$(shell (test -z "$$(git status --porcelain 2>/dev/null)" && echo clean) || echo dirty)

GOFMT :=gofmt -s
GOIMPORTS :=goimports -e

GO_IMPORT_PATH :=github.com/redhat-developer/opencompose
GO_FILES :=$(shell find . -name '*.go' -not -path './vendor/*' -print)
GO_PACKAGES =$(shell glide novendor)
ifneq ($(and $(GIT_VERSION),$(GIT_COMMIT),$(GIT_TREE_STATE)),)
GO_LDFLAGS =-ldflags="-X $(GO_IMPORT_PATH)/pkg/version.gitVersion=$(GIT_VERSION) -X $(GO_IMPORT_PATH)/pkg/version.gitCommit=$(GIT_COMMIT) -X $(GO_IMPORT_PATH)/pkg/version.gitTreeState=$(GIT_TREE_STATE)"
endif

.PHONY: build
build bin:
	go build $(GO_LDFLAGS)

.PHONY: install
install:
	go install $(GO_LDFLAGS)

.PHONY: test
test:
	go test $(GO_PACKAGES)

# compile opencompose for multiple platforms
.PHONY: cross
cross:
	gox -osarch="darwin/amd64 linux/amd64 linux/arm windows/amd64" -output="bin/opencompose-{{.OS}}-{{.Arch}}" $(GO_LDFLAGS)

.PHONY: checks
checks: check-gofmt check-goimports check-govet

.PHONY: check-gofmt
check-gofmt:
	@export files && files="$$($(GOFMT) -l $(GO_FILES))" && \
	if [ -n "$${files}" ]; then printf "ERROR: These files are not formated by $(GOFMT):\n"; printf "%s\n" $${files[@]}; exit 1; fi

.PHONY: check-goimports
check-goimports:
	@export files && files="$$($(GOIMPORTS) -l $(GO_FILES))" && \
	if [ -n "$${files}" ]; then printf "ERROR: These files are not formated by $(GOIMPORTS):\n"; printf "%s\n" $${files[@]}; exit 1; fi

.PHONY: check-govet
check-govet:
	go vet $(GO_PACKAGES)

.PHONY: check-strip-vendor
check-strip-vendor:
	@export vendors && vendors=$$(find ./vendor/ -mindepth 1 -type d -name 'vendor') && \
	if [ -n "$${vendors}" ]; then printf "ERROR: There are nested vendor directories: \n"; printf "%s\n" $${vendors[@]}; exit 1; fi
	@export files && files=$$($(do-strip-vendor) --dryrun) && \
	if [ -n "$${files}" ]; then printf "ERROR: There are unused files in ./vendor/\nRun 'make strip-vendor' to fix it.\n"; exit 1; fi

.PHONY: format
format: format-gofmt format-goimports

.PHONY: format-gofmt
format-gofmt:
	$(GOFMT) -w $(GO_FILES)

.PHONY: format-goimports
format-goimports:
	$(GOIMPORTS) -w $(GO_FILES)

do-strip-vendor :=glide-vc --only-code --no-tests --no-test-imports --no-legal-files

.PHONY: recreate-vendor
recreate-vendor:
	glide install --strip-vendor
	$(do-strip-vendor)

.PHONY: update-vendor
update-vendor:
	glide update --strip-vendor
	$(do-strip-vendor)

.PHONY: strip-vendor
strip-vendor:
	$(do-strip-vendor)

.PHONY: release
release: build
ifeq ($(and $(GIT_EXACT_VERSION),$(GIT_COMMIT)),)
	$(error GIT_EXACT_VERSION or GIT_COMMIT is not set or this is not a tagged commit)
endif
	tar -cJf opencompose-$(GIT_EXACT_VERSION)-$(GIT_COMMIT)-linux-64bit.tar.xz ./opencompose ./LICENSE ./README.md

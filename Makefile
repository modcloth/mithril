LIBS := \
  github.com/modcloth-labs/mithril
TARGETS := \
  $(LIBS) \
  github.com/modcloth-labs/mithril/mithril-server
REV_VAR := github.com/modcloth-labs/versioning.RevString
VERSION_VAR := github.com/modcloth-labs/versioning.VersionString
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR) $(REPO_REV) -X $(VERSION_VAR) $(REPO_VERSION)"


GO_TAG_ARGS ?= -tags full

ADDR := :8371
export ADDR

test: clean build
	go test $(GO_TAG_ARGS) -x $(TARGETS)

build: deps
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) -x $(TARGETS)

deps:
	go get $(GO_TAG_ARGS) -x -n $(TARGETS)

clean:
	find $${GOPATH%%:*}/pkg -name '*mithril*' -exec rm -v {} \;
	go clean -x $(TARGETS)

serve:
	$${GOPATH%%:*}/bin/mithril-server -d -a $(ADDR)

golden: test
	./runtests -v

.PHONY: build deps test clean serve

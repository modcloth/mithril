REV_VAR := github.com/modcloth/mithril.Rev
VERSION_VAR := github.com/modcloth/mithril.Version
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
GOBUILD_VERSION_ARGS := -ldflags "\
	-X $(REV_VAR) $(REPO_REV) \
	-X $(VERSION_VAR) $(REPO_VERSION)"

GO ?= go
DEPPY ?= deppy

GO_TAG_ARGS ?=

ADDR := :8371
export ADDR

.PHONY: all
all: clean golden

.PHONY: test
test: build
	go test $(GO_TAG_ARGS) -x ./...

.PHONY: build
build: deps
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) ./...
	go build -o $${GOPATH%%:*}/bin/mithril-server $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) ./mithril-server

.PHONY: deps
deps:
	$(DEPPY) restore

.PHONY: save
save:
	$(DEPPY) save ./...

.PHONY: clean
clean:
	$(GO) clean -x ./...
	if [ -d $${GOPATH%%:*}/pkg ] ; then \
		find $${GOPATH%%:*}/pkg -wholename '*modcloth/mithril*' -exec $(RM) -v {} \; ; \
	fi
	$(RM) .artifacts/*

.PHONY: distclean
distclean: clean

.PHONY: serve
serve:
	$${GOPATH%%:*}/bin/mithril-server -d -a $(ADDR)

.PHONY: golden
golden: test
	./runtests -v

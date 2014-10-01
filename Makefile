PACKAGE := github.com/modcloth-labs/mithril
SUBPACKAGES := \
	$(PACKAGE)/log \
	$(PACKAGE)/message \
	$(PACKAGE)/mithril-server \
	$(PACKAGE)/store
REV_VAR := github.com/modcloth-labs/mithril.Rev
VERSION_VAR := github.com/modcloth-labs/mithril.Version
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
	$(DEPPY) go test $(GO_TAG_ARGS) -x $(PACKAGE) $(SUBPACKAGES)

.PHONY: build
build: deps
	$(DEPPY) go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) -x $(PACKAGE) $(SUBPACKAGES)
	$(DEPPY) go build -o $${GOPATH%%:*}/bin/mithril-server $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) ./mithril-server

.PHONY: deps
deps:
	$(DEPPY) restore

.PHONY: save
save:
	$(DEPPY) save $(PACKAGE) $(SUBPACKAGES)

.PHONY: clean
clean:
	$(GO) clean -x $(PACKAGE) $(SUBPACKAGES) || true
	if [ -d $${GOPATH%%:*}/pkg ] ; then \
		find $${GOPATH%%:*}/pkg -wholename '*modcloth-labs/mithril*' -exec $(RM) -v {} \; ; \
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

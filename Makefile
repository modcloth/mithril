LIBS := mithril
REV_VAR := mithril.RevString
VERSION_VAR := mithril.VersionString
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR) $(REPO_REV) -X $(VERSION_VAR) $(REPO_VERSION)"
JOHNNY_DEPS_REV := f2af161b01bcda148859a0f7d0524769186b339b


GO_TAG_ARGS ?= -tags full

ADDR := :8371
export ADDR

all: clean golden

test: build
	go test $(GO_TAG_ARGS) -x $(LIBS)

build: deps
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) -x $(LIBS)
	go build -o $${GOPATH%%:*}/bin/mithril-server $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) ./mithril-server

deps: johnny_deps
	if [ ! -L $${GOPATH%%:*}/src/mithril ] ; then gvm linkthis ; fi
	./johnny_deps

johnny_deps:
	curl -s -o $@ https://raw.github.com/meatballhat/johnny-deps/$(JOHNNY_DEPS_REV)/bin/johnny_deps
	chmod +x $@

clean:
	go clean -x $(LIBS) || true
	if [ -d $${GOPATH%%:*}/pkg ] ; then \
		find $${GOPATH%%:*}/pkg -name '*mithril*' -exec rm -v {} \; ; \
	fi

distclean: clean
	rm -f ./johnny_deps

serve:
	$${GOPATH%%:*}/bin/mithril-server -d -a $(ADDR)

golden: test
	./runtests -v

.PHONY: all build deps test clean distclean serve

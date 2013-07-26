LIBS := mithril
REV_VAR := mithril.RevString
VERSION_VAR := mithril.VersionString
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR) $(REPO_REV) -X $(VERSION_VAR) $(REPO_VERSION)"


GO_TAG_ARGS ?= -tags full

ADDR := :8371
export ADDR

all: golden

test: clean build
	go test $(GO_TAG_ARGS) -x $(LIBS)

build: deps
	go install $(GOBUILD_VERSION_ARGS) $(GO_TAG_ARGS) -x $(LIBS)
	go build -o $${GOPATH%%:*}/bin/mithril-server ./mithril-server

deps:
	if [ ! -L $${GOPATH%%:*}/src/mithril ] ; then gvm linkthis ; fi
	./install-deps ./deps.txt

clean:
	go clean -x $(LIBS) || true
	find $${GOPATH%%:*}/pkg -name '*mithril*' -exec rm -v {} \;

serve:
	$${GOPATH%%:*}/bin/mithril-server -d -a $(ADDR)

golden: test
	./runtests -v

.PHONY: all build deps test clean serve

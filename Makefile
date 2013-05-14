LIBS := \
  github.com/modcloth-labs/mithril
TARGETS := \
  $(LIBS) \
  github.com/modcloth-labs/mithril/mithril-server

GO_TAG_ARGS := -tags pg

ADDR := :8371
export ADDR

test: build
	go test $(GO_TAG_ARGS) -x $(LIBS)

build: deps
	go install $(GO_TAG_ARGS) -x $(TARGETS)

deps:
	go get $(GO_TAG_ARGS) -x -n $(TARGETS)

clean:
	go clean -x $(TARGETS)

serve:
	$${GOPATH%%:*}/bin/mithril-server -a $(ADDR)

golden: test
	./runtests -v

.PHONY: build deps test clean serve

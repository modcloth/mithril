LIBS := \
  github.com/modcloth-labs/mithril
TARGETS := \
  $(LIBS) \
  github.com/modcloth-labs/mithril/mithril-server

ADDR := :8371

test: build
	go test -x $(LIBS)

build: deps
	go install -x $(TARGETS)

deps:
	go get -x -n $(TARGETS)

clean:
	go clean -x $(TARGETS)

serve:
	$${GOPATH%%:*}/bin/mithril-server -a $(ADDR)

.PHONY: build deps test clean serve

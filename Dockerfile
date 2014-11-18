FROM golang:1.3.3

RUN go get github.com/meatballhat/deppy

ADD . /go/src/github.com/modcloth/mithril

WORKDIR /go/src/github.com/modcloth/mithril

RUN touch Makefile \
  && make build \
  && rm -rf $GOPATH/src \
  && rm -rf $GOPATH/pkg \
  && rm -f $GOPATH/bin/deppy

ENTRYPOINT ["mithril-server"]

CMD ["-h"]

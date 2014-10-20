FROM golang:1.3.3

ADD . /go/src/github.com/modcloth/mithril

WORKDIR /go/src/github.com/modcloth/mithril

RUN go get -x github.com/meatballhat/deppy

RUN make build

ENTRYPOINT ["mithril-server"]

CMD ["-h"]

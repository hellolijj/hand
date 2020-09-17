FROM golang:1.10-alpine as build

# Copy in the go src
WORKDIR /go/src/github.com/hellolijj/hand
COPY cmd/    cmd/
COPY vendor/ vendor/
COPY pkg/    pkg/

RUN go build -o /go/bin/hand cmd/*.go

FROM alpine

COPY --from=build /go/bin/hand /usr/bin/hand

CMD ["hand"]
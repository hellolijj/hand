FROM golang:1.10-alpine

# Copy in the go src
WORKDIR /go/src/github.com/hellolijj/hand

COPY pkg/    pkg/
COPY main.go  main.go

RUN go build -o /go/bin/hand main.go

CMD ["hand"]
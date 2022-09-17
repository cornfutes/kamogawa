FROM golang:1.19 as built
WORKDIR /go/app/api
COPY . .
ENV CGO_ENABLED=0
RUN go get -d -v ./...
RUN go build -o /tmp/kamogawa-server ./*.go

FROM postgres:13.8-alpine
COPY --from=built /tmp/kamogawa-server /usr/bin/kamogawa-server
WORKDIR /usr/bin
CMD kamogawa-server;


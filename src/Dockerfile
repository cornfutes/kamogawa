FROM golang:1.19 as builder
WORKDIR /foo
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /foo/kamogawa-server

FROM alpine
COPY --from=builder /foo/kamogawa-server /usr/bin/kamogawa-server
CMD kamogawa-server;
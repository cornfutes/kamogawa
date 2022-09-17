FROM golang:1.17 as builder
WORKDIR /foo
COPY . .
RUN go mod download
# https://stackoverflow.com/questions/22267189/what-does-the-w-flag-mean-when-passed-in-via-the-ldflags-option-to-the-go-comman/22276273#22276273
RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o /foo/kamogawa-server

FROM alpine
COPY --from=builder /foo/kamogawa-server /usr/bin/kamogawa-server
CMD kamogawa-server;
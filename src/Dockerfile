FROM golang:1.19 as builder
WORKDIR /foo
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /foo/kamogawa

FROM scratch
WORKDIR /app
COPY --from=builder /foo/kamogawa ./main
EXPOSE 3000
ENTRYPOINT ["./main"]
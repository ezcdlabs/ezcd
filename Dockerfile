FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapi ./cmd/api

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/myapi .
CMD ["./myapi"]
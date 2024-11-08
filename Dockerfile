FROM golang:1.23.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ezcd-server ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/ezcd-server .
CMD ["./ezcd-server"]
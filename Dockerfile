FROM golang:1.23.2 AS builder
WORKDIR /app
RUN go mod init kino-cat-bot-go
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /app/kino-cat-bot-go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/kino-cat-bot-go .
CMD ["./kino-cat-bot-go"]

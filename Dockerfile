FROM golang:1.20 AS builder
WORKDIR /app

RUN go mod download
RUN go mod init kino-cat-bot-go
RUN go get gopkg.in/telebot.v3
RUN go get github.com/nats-io/nats.go

COPY . .
RUN CGO_ENABLED=0 go build -o /app/kino-cat-bot-go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/kino-cat-bot-go .

CMD ["./kino-cat-bot-go"]

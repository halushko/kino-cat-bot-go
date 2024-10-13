FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/kino-cat-bot-go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/kino-cat-bot-go .

CMD ["./kino-cat-bot-go"]

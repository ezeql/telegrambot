FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o oab-price-bot

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/oab-price-bot .


CMD ["./oab-price-bot"]
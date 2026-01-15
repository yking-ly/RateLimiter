FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o rate-limiter-server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/rate-limiter-server .

COPY --from=builder /app/public ./public

EXPOSE 8080

CMD ["./rate-limiter-server"]

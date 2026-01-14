FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/http-client ./cmd/http-client

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/http-client .

CMD ["./http-client"]

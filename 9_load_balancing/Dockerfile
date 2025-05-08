# Stage 1: Build
FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o main .

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .

CMD ["./main"]
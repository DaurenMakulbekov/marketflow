FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main ./cmd

FROM alpine

COPY --from=builder /app/main /main

CMD ["./main"]
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o marketflow ./cmd

FROM alpine

COPY --from=builder /app/marketflow /marketflow

CMD ["./marketflow"]

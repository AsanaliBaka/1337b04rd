FROM golang:1.22.2 AS builder

WORKDIR /app
COPY . .
RUN go build -o /app/1337b04rd ./cmd


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/1337b04rd .

CMD ["/app/1337b04rd"]
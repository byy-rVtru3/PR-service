FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/avitoTech ./cmd/server

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/avitoTech /app/avitoTech

EXPOSE 8080

CMD ["/app/avitoTech"]

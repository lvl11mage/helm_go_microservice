FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum order.go ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o order-service .


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/order-service .

ENV GOMEMLIMIT=512MiB

EXPOSE 8081

CMD ["./order-service"]

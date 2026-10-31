FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum shipping.go ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o shipping-service .


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/shipping-service .

ENV GOMEMLIMIT=512MiB

EXPOSE 8084

CMD ["./shipping-service"]

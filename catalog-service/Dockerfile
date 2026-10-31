FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum catalog.go ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o catalog-service .


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/catalog-service .

ENV GOMEMLIMIT=512MiB

EXPOSE 8086

CMD ["./catalog-service"]

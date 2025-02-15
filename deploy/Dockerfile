FROM golang:1.23.6 as builder

WORKDIR /app

COPY ../go.mod .
COPY ../go.sum .
RUN go mod download

COPY .. .

RUN go build -o bin/avito_shop cmd/avito_shop/main.go

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/bin/avito_shop .

CMD ["./avito_shop"]

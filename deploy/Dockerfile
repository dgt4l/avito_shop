FROM golang:1.23.6-alpine3.20 as builder

WORKDIR /app

COPY ../go.mod .
COPY ../go.sum .
RUN go mod download

COPY .. .

RUN go build -o bin/avito_shop cmd/avito_shop/main.go

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder /app/bin/avito_shop .

EXPOSE 8080

CMD ["./avito_shop"]

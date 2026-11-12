FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM debian:stable-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main ./main

# El backend original necesita este archivo obligatoriamente.
COPY app.env ./app.env

EXPOSE 8080

CMD ["./main"]
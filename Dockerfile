FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY .env ./
COPY migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o song-library cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/song-library .
COPY --from=builder /app/.env ./
COPY --from=builder /app/migrations ./migrations

CMD [ "./song-library" ]
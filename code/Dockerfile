FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /out/server /app/server
COPY migrations /app/migrations
COPY openapi /app/openapi

EXPOSE 8080

ENTRYPOINT ["/app/server"]

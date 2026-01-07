FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary dari root main.go (Clean Architecture entrypoint)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/app ./main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/app /app/app
# COPY .env /app/.env  # optional: bake env, or pass via --env/--env-file

USER app
EXPOSE 8080
ENTRYPOINT ["/app/app"]

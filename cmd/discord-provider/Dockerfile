# Step 1 - compile code binary
FROM golang:1.24.1-alpine AS builder

LABEL maintainer="Club CEDILLE <cedille@etsmtl.net>"

RUN apk add --no-cache --update ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o ./discord_provider ./cmd/discord-provider/main.go

# Step 2 - import necessary files to run program.
FROM alpine

COPY --from=builder /app/discord_provider /discord_provider
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV DISCORD_WEBHOOK_URL=$DISCORD_WEBHOOK_URL

ENTRYPOINT ["/discord_provider"]

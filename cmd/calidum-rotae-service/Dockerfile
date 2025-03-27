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
    go build -o ./calidum_rotae_service ./cmd/calidum-rotae-service/main.go

# Step 2 - import necessary files to run program.
FROM alpine

COPY --from=builder /app/calidum_rotae_service /calidum_rotae_service
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV CALIDUM_ROTAE_SERVICE_API_KEY=$CALIDUM_ROTAE_SERVICE_API_KEY

# Keep this value empty if you don't want to use TLS when dialing the microservices.
# Otherwise, pass the path to the certificate file.
ENV CERTIFICATE_FILE_PATH=$CERTIFICATE_FILE_PATH

ENTRYPOINT ["/calidum_rotae_service"]

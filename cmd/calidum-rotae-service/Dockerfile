# Step 1 - compile code binary
FROM golang:1.21.1-alpine AS builder

LABEL maintainer="Club CEDILLE <cedille@etsmtl.net>"

RUN apk add --no-cache --update ca-certificates

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o ./calidum_rotae_service ./cmd/calidum-rotae-service/main.go

# Add user & group
RUN addgroup -S cedille-group && \
    adduser -S cedille-user -G cedille-group

# Step 2 - import necessary files to run program.
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/calidum_rotae_service /calidum_rotae_service
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER cedille-user

ENTRYPOINT ["/calidum_rotae_service"]
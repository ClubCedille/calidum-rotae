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
    go build -o ./shell_provider ./cmd/shell-provider/main.go

# Step 2 - import necessary files to run program.
FROM alpine

COPY --from=builder /app/shell_provider /shell_provider
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/email_provider"]

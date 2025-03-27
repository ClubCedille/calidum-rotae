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
    go build -o ./email_provider ./cmd/email-provider/main.go

# Step 2 - import necessary files to run program.
FROM alpine

COPY --from=builder /app/email_provider /email_provider
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV EMAIL_FROM_ADDRESS=$EMAIL_FROM_ADDRESS
ENV EMAIL_FROM_NAME=$EMAIL_FROM_NAME
ENV EMAIL_NAME_TO=$EMAIL_NAME_TO
ENV EMAIL_TO_ADDRESS=$EMAIL_TO_ADDRESS
ENV EMAIL_SUBJECT=$EMAIL_SUBJECT
ENV EMAIL_SMTP_API_KEY=$EMAIL_SMTP_API_KEY

ENTRYPOINT ["/email_provider"]

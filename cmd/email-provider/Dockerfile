# Step 1 - compile code binary
FROM golang:1.18.4-alpine AS builder

LABEL maintainer="Club CEDILLE <cedille@etsmtl.net>"

RUN apk add --no-cache --update ca-certificates

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o ./email_provider ./cmd/email-provider/main.go

# Add user & group
RUN addgroup -S cedille-group && \
    adduser -S cedille-user -G cedille-group

ENV FROM_EMAIL_ADDRESS=$FROM_EMAIL_ADDRESS
ENV FROM_NAME=$FROM_NAME
ENV TO_NAME=$TO_NAME
ENV TO_EMAIL_ADDRESS=$TO_EMAIL_ADDRESS
ENV SENDGRID_API_KEY=$SENDGRID_API_KEY
ENV SUBJECT=$SUBJECT

# Step 2 - import necessary files to run program.
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/email_provider /email_provider
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER cedille-user

ENTRYPOINT ["/email_provider"]
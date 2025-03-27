package server

import (
	"context"
	"fmt"
	"log"
	"os"

	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
    "github.com/resend/resend-go/v2"
)

// environment variables
const (
	ENV_MAIL_FROM_ADDRESS      = "EMAIL_FROM_ADDRESS"
	ENV_EMAIL_FROM_NAME        = "EMAIL_FROM_NAME"
	ENV_EMAIL_NAME_TO          = "EMAIL_NAME_TO"
	ENV_EMAIL_TO_ADDRESS       = "EMAIL_TO_ADDRESS"
	ENV_EMAIL_SUBJECT          = "EMAIL_SUBJECT"
	ENV_EMAIL_SMTP_API_KEY     = "EMAIL_SMTP_API_KEY"
)

type Server struct {
	email_provider.EmailProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SendEmail(ctx context.Context, message *email_provider.SendEmailRequest) (*email_provider.SendEmailResponse, error) {
    
	emailFromAddress, exists := os.LookupEnv(ENV_MAIL_FROM_ADDRESS)
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_MAIL_FROM_ADDRESS)
	}

	// emailFromName, exists := os.LookupEnv(ENV_EMAIL_FROM_NAME)
	// if !exists {
	// 	return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_EMAIL_FROM_NAME)
	// }
    
	emailNameTo, exists := os.LookupEnv(ENV_EMAIL_NAME_TO)
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_EMAIL_NAME_TO)
	}

	emailAddressTo, exists := os.LookupEnv(ENV_EMAIL_TO_ADDRESS)
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_EMAIL_TO_ADDRESS)
	}

    emailSMTPApiKey, exists := os.LookupEnv(ENV_EMAIL_SMTP_API_KEY)
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_EMAIL_SMTP_API_KEY)
	}
	emailSubject := os.Getenv(ENV_EMAIL_SUBJECT) // EMAIL_SUBJECT can be null

    client := resend.NewClient(emailSMTPApiKey)

	htmlContent := fmt.Sprintf(
		`
		 <strong> First name: </strong> %s <br> 
		 <strong> Last name: </strong> %s <br>
		 <strong> Email: </strong> %s <br> 
		 <strong> Request details: </strong> %s <br> 
		 <strong> Request service: </strong> %s
		`,
		message.Sender.FirstName, message.Sender.LastName, message.Sender.Email,
		message.RequestDetails, message.RequestService)
    
    params := &resend.SendEmailRequest{
        To:      []string{emailAddressTo},
        From:    emailFromAddress,
        Text:    "hello " + emailNameTo,
        Html:    htmlContent,
        Subject: emailSubject,
        ReplyTo: "noreply@cedille.club",
    }

    sent, err := client.Emails.Send(params)

    if err != nil {
        log.Printf("Email didnt send ")
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error sending email: %s", err)
	}
    
	log.Printf("Email sent id: ", sent)
	return &email_provider.SendEmailResponse{}, nil
}

package server

import (
	"context"
	"fmt"
	"log"
	"os"

	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// environment variables
const (
	ENV_MAIL_FROM_ADDRESS      = "EMAIL_FROM_ADDRESS"
	ENV_EMAIL_FROM_NAME        = "EMAIL_FROM_NAME"
	ENV_EMAIL_NAME_TO          = "EMAIL_NAME_TO"
	ENV_EMAIL_TO_ADDRESS       = "EMAIL_TO_ADDRESS"
	ENV_EMAIL_SUBJECT          = "EMAIL_SUBJECT"
	ENV_EMAIL_SENDGRID_API_KEY = "EMAIL_SENDGRID_API_KEY"
)

type Server struct {
	email_provider.EmailProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SendEmail(ctx context.Context, message *email_provider.SendEmailRequest) (*email_provider.SendEmailResponse, error) {
	m := mail.NewV3Mail()
	emailFromAddress, exists := os.LookupEnv(ENV_MAIL_FROM_ADDRESS)
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_MAIL_FROM_ADDRESS)
	}

	emailFromName, exists := os.LookupEnv(ENV_EMAIL_FROM_NAME)
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_EMAIL_FROM_NAME)
	}
	e := mail.NewEmail(emailFromName, emailFromAddress)
	m.SetFrom(e)

	emailNameTo, exists := os.LookupEnv(ENV_EMAIL_NAME_TO)
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_EMAIL_NAME_TO)
	}

	emailAddressTo, exists := os.LookupEnv(ENV_EMAIL_TO_ADDRESS)
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_EMAIL_TO_ADDRESS)
	}

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		// TODO Send email to more addresses
		mail.NewEmail(emailNameTo, emailAddressTo),
	}
	p.AddTos(tos...)
	p.Subject = os.Getenv(ENV_EMAIL_SUBJECT) // EMAIL_SUBJECT can be null
	m.AddPersonalizations(p)

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
	c := mail.NewContent("text/html", htmlContent)
	m.AddContent(c)

	mailSettings := mail.NewMailSettings()
	bypassBounceManagement := mail.NewSetting(true)
	mailSettings.SetBypassBounceManagement(bypassBounceManagement)
	bypassSpamManagement := mail.NewSetting(true)
	mailSettings.SetBypassSpamManagement(bypassSpamManagement)
	bypassUnsubscribeManagement := mail.NewSetting(true)
	mailSettings.SetBypassUnsubscribeManagement(bypassUnsubscribeManagement)

	footerSetting := mail.NewFooterSetting()
	footerSetting.SetText("footer")
	footerSetting.SetEnable(true)
	footerSetting.SetHTML("<html><body><br><br>Club CEDILLE, ETS</body></html>")
	mailSettings.SetFooter(footerSetting)

	spamCheckSetting := mail.NewSpamCheckSetting()
	spamCheckSetting.SetEnable(true)
	spamCheckSetting.SetSpamThreshold(1)
	spamCheckSetting.SetPostToURL("https://spamcatcher.sendgrid.com")
	mailSettings.SetSpamCheckSettings(spamCheckSetting)
	m.SetMailSettings(mailSettings)

	replyToEmail := mail.NewEmail(emailFromName, emailFromAddress)
	m.SetReplyTo(replyToEmail)

	emailSendgridApiKey, exists := os.LookupEnv(ENV_EMAIL_SENDGRID_API_KEY)
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error getting the env var %s", ENV_EMAIL_SENDGRID_API_KEY)
	}
	request := sendgrid.GetRequest(emailSendgridApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	_, err := sendgrid.API(request)
	if err != nil {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("error sending email: %s", err)
	}

	log.Printf("Email sent")
	return &email_provider.SendEmailResponse{}, nil
}

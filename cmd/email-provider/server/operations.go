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

type Server struct {
	email_provider.EmailProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SendMessage(ctx context.Context, message *email_provider.SendEmailRequest) (*email_provider.SendEmailResponse, error) {
	m := mail.NewV3Mail()
	emailFromAddress, exists := os.LookupEnv("EMAIL_FROM_ADDRESS")
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("the env var EMAIL_FROM_ADDRESS is not set to a value")
	}

	emailFromName, exists := os.LookupEnv("EMAIL_FROM_NAME")
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("the env var EMAIL_FROM_NAME is not set to a value")
	}
	e := mail.NewEmail(emailFromName, emailFromAddress)
	m.SetFrom(e)

	emailNameTo, exists := os.LookupEnv("EMAIL_NAME_TO")
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("the env var EMAIL_NAME_TO is not set to a value")
	}

	emailAddressTo, exists := os.LookupEnv("EMAIL_TO_ADDRESS")
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("the env var EMAIL_TO_ADDRESS is not set to a value")
	}

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		// TODO Send email to more addresses
		mail.NewEmail(emailNameTo, emailAddressTo),
	}
	p.AddTos(tos...)
	p.Subject = os.Getenv("EMAIL_SUBJECT") // EMAIL_SUBJECT can be null
	m.AddPersonalizations(p)

	htmlContent := fmt.Sprintf(
		`
		 <strong> First name: </strong> %s <br> 
		 <strong> Last name: </strong> %s <br>
		 <strong> Email: </strong> %s <br> 
		 <strong> Phone: </strong> %s <br><br> 
		 <strong> Request details: </strong> %s <br> 
		 <strong> Request service: </strong> %s
		`,
		message.Sender.FirstName, message.Sender.LastName, message.Sender.Email,
		message.Sender.PhoneNumber, message.RequestDetails, message.RequestService)
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

	emailSendgridApiKey, exists := os.LookupEnv("EMAIL_SENDGRID_API_KEY")
	if !exists {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("the env var EMAIL_SENDGRID_API_KEY is not set to a value")
	}
	request := sendgrid.GetRequest(emailSendgridApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		return &email_provider.SendEmailResponse{}, fmt.Errorf("failed to send email: %s", err)
	}

	log.Printf("Email sent: %s", response.Body)
	return &email_provider.SendEmailResponse{}, nil
}

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
	address := os.Getenv("FROM_EMAIL_ADDRESS")
	name := os.Getenv("FROM_NAME")
	e := mail.NewEmail(name, address)
	m.SetFrom(e)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(os.Getenv("TO_NAME"), os.Getenv("TO_EMAIL_ADDRESS")),
	}
	p.AddTos(tos...)
	p.Subject = os.Getenv("SUBJECT")
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
		message.Sender.FirstName, message.Sender.LastName, message.Sender.Email, message.Sender.PhoneNumber, message.RequestDetails, message.RequestService)
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

	replyToEmail := mail.NewEmail(os.Getenv("FROM_NAME"), os.Getenv("FROM_EMAIL_ADDRESS"))
	m.SetReplyTo(replyToEmail)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
	}

	log.Printf("Received message content from client: %s", message.Sender)
	return &email_provider.SendEmailResponse{}, nil
}

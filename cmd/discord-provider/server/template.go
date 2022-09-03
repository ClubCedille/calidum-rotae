package server

import "fmt"

// Necessary struct to send a embedded discord message
type discordMessage struct {
	Username  string     `json:"username,omitempty"`
	Content   string     `json:"content,omitempty"`
	Embeddeds []embedded `json:"embeds,omitempty"`
}

type embedded struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Footer      footer `json:"footer,omitempty"`
}

type footer struct {
	Text string `json:"text,omitempty"`
}

type discordSenderInformation struct {
	firstName      string
	lastName       string
	email          string
	phoneNumber    string
	requestDetails string
	requestService string
}

func (senderInformation discordSenderInformation) String() string {
	discordName := fmt.Sprintf("**Name:** %s %s", senderInformation.firstName, senderInformation.lastName)
	discordEmail := fmt.Sprintf("**Email:** %s", senderInformation.email)
	discordPhone := fmt.Sprintf("**Phone:** %s", senderInformation.phoneNumber)
	discordRequestDetails := fmt.Sprintf("**Request details:** %s", senderInformation.requestDetails)
	discordRequestServices := fmt.Sprintf("**Request services:** %s", senderInformation.requestService)
	return fmt.Sprintf("%s\n%s\n%s\n\n%s\n%s", discordName, discordEmail, discordPhone, discordRequestDetails, discordRequestServices)
}

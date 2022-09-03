package server

// Necessary struct to send a embedded discord message
type Message struct {
	Username  string     `json:"username,omitempty"`
	Content   string     `json:"content,omitempty"`
	Embeddeds []Embedded `json:"embeds,omitempty"`
}

type Embedded struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Footer      Footer `json:"footer,omitempty"`
}

type Footer struct {
	Text string `json:"text,omitempty"`
}

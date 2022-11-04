package main

type MailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	FromMail string
	FromName string
}

type MailContent struct {
	To          string `json:"to" bson:"to"`
	ToName      string `json:"to_name" bson:"to_name"`
	Subject     string `json:"subject" bson:"subject"`
	Body        []byte `json:"body" bson:"body"`
	MimeVersion string `json:"mime_version" bson:"mime_version"`
	ContentType string `json:"content_type" bson:"content_type"`
}

func NewMailConfig(host, port, username, password, fromMail, fromName string) *MailConfig {
	return &MailConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		FromName: fromName,
		FromMail: fromMail,
	}
}

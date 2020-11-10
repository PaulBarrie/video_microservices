package config

import "os"

// SmtpCli contains details to connect to the Smtp server
type SmtpCli struct {
	Host     string
	Port     string
	Username string
	Password string
	Addr     string
	Email    string
	Name     string
}

//InitSmtp init Smtp client
func (a *App) InitSmtp() {
	a.Smtp = &SmtpCli{
		Host:     os.Getenv("SMTP_CONTAINER"),
		Port:     "25",
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PWD"),
		Email:    os.Getenv("SMTP_EMAIL"),
		Name:     os.Getenv("SMTP_NAME"),
	}
	(*a.Smtp).Addr = (*a.Smtp).Host + ":" + (*a.Smtp).Port
}

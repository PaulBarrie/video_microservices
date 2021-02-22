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

//InitSMTP init Smtp client
func (a *App) InitSMTP() {
	a.Smtp = &SmtpCli{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("PASSWORD"),
		Email:    os.Getenv("SMTP_EMAIL"),
		Name:     os.Getenv("SMTP_NAME"),
	}
	(*a.Smtp).Addr = (*a.Smtp).Host + ":" + (*a.Smtp).Port
}

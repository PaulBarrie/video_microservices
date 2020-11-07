package config

import "os"

// PostfixCli contains details to connect to the Postfix server
type PostfixCli struct {
	Host     string
	Port     string
	Username string
	Password string
	Addr     string
	Email    string
	Name     string
}

//InitPostfix init Postfix client
func (a *App) InitPostfix() {
	a.Postfix = &PostfixCli{
		Host:     os.Getenv("POSTFIX_HOST"),
		Port:     os.Getenv("POSTFIX_PORT"),
		Username: os.Getenv("POSTFIX_USERNAME"),
		Password: os.Getenv("POSTFIX_PWD"),
		Email:    os.Getenv("POSTFIX_EMAIL"),
		Name:     os.Getenv("POSTFIX_NAME"),
	}
	(*a.Postfix).Addr = (*a.Postfix).Host + ":" + (*a.Postfix).Port
}

package models

import (
	"time"
)

type Token struct {
	Id         int       `json:"id"`
	Code       string    `json:"code"`
	Expired_at time.Time `json:"expired_at"`
	User_id    int       `json:"user_id"`
}

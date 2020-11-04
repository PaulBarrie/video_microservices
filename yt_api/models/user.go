package models

import (
	"time"
)

type User struct {
	Id         int       `json:"id"`
	Username   string    `json:"username"`
	Pseudo     string    `json:"pseudo"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Created_at time.Time `json:"created_at"`
}

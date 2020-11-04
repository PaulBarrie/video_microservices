package models

import (
	"time"
)

type Video struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Duration   float64 `json:"duration"`
	User_id    int       `json:"user_id"`
	Source     string    `json:"source"`
	Created_at time.Time `json:"created_at"`
	View       int       `json:"view"`
	Enabled    int       `json:"enabled"`
}

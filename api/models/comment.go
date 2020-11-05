package models

type Comment struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	User_id  int    `json:"user_id"`
	Video_id int    `json:"video_id"`
}

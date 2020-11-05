package utils

import "models"

type RespUser struct {
	Message string      `json:"message"`
	Data    models.User `json:"data"`
}

type RespToken struct {
	Message string       `json:"message"`
	Data    models.Token `json:"data"`
}

type Pager struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type RespUserPaginated struct {
	Message string        `json:"message"`
	Data    []models.User `json:"data"`
	Pager   Pager         `json:"pager"`
}

type RespVideo struct {
	Message string       `json:"message"`
	Data    models.Video `json:"data"`
}

type RespVideoPaginated struct {
	Message string         `json:"message"`
	Data    []models.Video `json:"data"`
	Pager   Pager          `json:"pager"`
}

type VideoInfo struct {
	Location string
	Size     int64
}
type RespComment struct {
    Message string         `json:"message"`
    Data    models.Comment `json:"data"`
}

type RespCommentPaginated struct {
    Message string           `json:"message"`
    Data    []models.Comment `json:"data"`
    Pager   Pager            `json:"pager"`
}

type VideoDetails struct {
	Duration float64 `json:"duration"`
	Quality int64 `json:"quality"`
}



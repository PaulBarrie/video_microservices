package utils

import "models"

//RespUser struct resp for user creation
type RespUser struct {
	Message string      `json:"message"`
	Data    models.User `json:"data"`
}

//RespToken struct resp for auth
type RespToken struct {
	Message string       `json:"message"`
	Data    models.Token `json:"data"`
}

//Pager struct for pagination
type Pager struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

//RespUserPaginated user with pagination
type RespUserPaginated struct {
	Message string        `json:"message"`
	Data    []models.User `json:"data"`
	Pager   Pager         `json:"pager"`
}

// RespVideo video resp struct
type RespVideo struct {
	Message string       `json:"message"`
	Data    models.Video `json:"data"`
}

// RespVideoPaginated video resp structpaginated
type RespVideoPaginated struct {
	Message string         `json:"message"`
	Data    []models.Video `json:"data"`
	Pager   Pager          `json:"pager"`
}

//VideoInfo vide details struct
type VideoInfo struct {
	Location string
	Size     int64
}

//RespComment struct of comment response
type RespComment struct {
	Message string         `json:"message"`
	Data    models.Comment `json:"data"`
}

// RespCommentPaginated struct of comment response paginated
type RespCommentPaginated struct {
	Message string           `json:"message"`
	Data    []models.Comment `json:"data"`
	Pager   Pager            `json:"pager"`
}

// VideoDetails struct
type VideoDetails struct {
	Duration float64 `json:"duration"`
	Quality  int64   `json:"quality"`
}

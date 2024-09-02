package models

type GetRedirectsRes struct {
	Page  int        `json:"page"`
	Data  []Redirect `json:"data"`
	Count int64      `json:"count"`
}

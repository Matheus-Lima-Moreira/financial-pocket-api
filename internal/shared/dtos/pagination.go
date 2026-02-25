package dtos

type PaginationDTO struct {
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	Total    int `json:"total"`
	Next     int `json:"next"`
	Previous int `json:"previous"`
}

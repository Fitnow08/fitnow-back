package domain

type ErrorResponse struct {
	Description string `json:"description"`
	Key         string `json:"key"`
}

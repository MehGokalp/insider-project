package dto

type SendMessageRequest struct {
	To      string `json:"to" validate:"required"`
	Content string `json:"content" validate:"required,max=160"`
}

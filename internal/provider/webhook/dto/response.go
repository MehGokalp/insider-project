package dto

type SendMessageResponse struct {
	ID     string `json:"messageId" validate:"required"`
	Status string `json:"message" validate:"required"`
}

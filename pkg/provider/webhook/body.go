package webhook

type SendMessageRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

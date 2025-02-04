package list

type Message struct {
	ID      uint   `json:"id"`
	To      string `json:"to"`
	Content string `json:"content"`
	Sent    bool   `json:"sent"`
	SentAt  string `json:"sent_at"`
}

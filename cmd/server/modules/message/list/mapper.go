package list

import (
	"github.com/mehgokalp/insider-project/pkg/database"
	"time"
)

func mapMessageList(messages []database.Message) []Message {
	if messages == nil {
		return []Message{}
	}

	var r []Message

	for _, m := range messages {
		r = append(
			r, Message{
				ID:      m.ID,
				To:      m.To,
				Content: m.Content,
				Sent:    m.Sent,
				SentAt:  m.SentAt.Format(time.RFC3339),
			},
		)
	}

	return r
}

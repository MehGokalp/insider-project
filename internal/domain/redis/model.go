package redis

import "time"

type ModelInterface interface {
	Duration() time.Duration
}

type Message struct {
	ID   string `json:"id"`
	Time string `json:"time"`
}

func (m Message) Duration() time.Duration {
	return time.Hour * 24
}

type MessageEngineRunningStatus struct {
	Consume bool `json:"consume"`
}

func (s MessageEngineRunningStatus) Duration() time.Duration {
	return 0
}

package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

type ActiveQueue struct {
	Id         ksuid.KSUID `json:"id"`
	QueueName  string      `json:"queue_name"`
	PollExpiry time.Time   `json:"poll_expiry"`
	Payload    string      `json:"payload"`
}

type IdleQueue struct {
	Items []string `json:"items"`
	// add other info
}

type EnqueuePayload struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

type QueueData struct {
	ActiveQueue map[ksuid.KSUID]ActiveQueue `json:"active_queue"`
	Queue       map[string]IdleQueue        `json:"queue"`
}

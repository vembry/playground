package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

type ActiveQueue struct {
	Id         ksuid.KSUID `json:"id"`
	QueueName  string      `json:"queue_name"`
	PollExpiry time.Time   `json:"poll_expiry"`
	Queue      *Queue      `json:"queue"`
}

type IdleQueue struct {
	Items []*Queue `json:"items"`
	// add other info
	// ...
}

type EnqueuePayload struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

type Queue struct {
	Payload string
	// add other info
	// ...
}

type QueueData struct {
	ActiveQueueCount int64 `json:"active_queue"`
	IdleQueueCount   int64 `json:"idle_queue"`
}

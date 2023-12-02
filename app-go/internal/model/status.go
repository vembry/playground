package model

type Status string

const (
	StatusPending   Status = "pending"
	StatusFailed    Status = "failed"
	StatusCompleted Status = "completed"
)

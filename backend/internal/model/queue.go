package model

import "time"

type QueueState struct {
	ID            int
	CurrentLetter *string
	CurrentNumber *int
	CurrentQueue  string
	UpdatedAt     time.Time
	IssuedAt      *time.Time `json:"-"`
}

func (QueueState) TableName() string {
	return "queue_state"
}

type QueueHistory struct {
	ID          int64
	QueueNumber string
	CreatedAt   time.Time
}

func (QueueHistory) TableName() string {
	return "queue_history"
}

type QueueResponse struct {
	CurrentQueue string     `json:"current_queue,omitempty"`
	QueueNumber  string     `json:"queue_number,omitempty"`
	IssuedAt     *time.Time `json:"issued_at,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

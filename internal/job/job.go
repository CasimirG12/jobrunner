package job

import (
	"slices"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

var allowedTransitions = map[Status][]Status{
	StatusPending: {
		StatusRunning,
	},
	StatusRunning: {
		StatusCompleted,
		StatusFailed,
		StatusPending,
	},
	StatusFailed: {
		StatusPending,
	},
	StatusCompleted: {},
}

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusRunning, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

func (from Status) CanTransitionTo(to Status) bool {
	return slices.Contains(allowedTransitions[from], to)
}

type Kind string

const (
	KindWordCount     Kind = "word_count"
	KindCalculateHash Kind = "calculate_hash"
)

func (k Kind) IsValid() bool {
	switch k {
	case KindWordCount, KindCalculateHash:
		return true
	default:
		return false
	}
}

type Job struct {
	ID           string     `json:"id"`
	Kind         Kind       `json:"kind"`
	Status       Status     `json:"status"`
	Input        string     `json:"input"`
	Attempts     int        `json:"attempts"`
	Result       string     `json:"result"`
	ErrorMessage string     `json:"error_message"`
	CreatedAt    time.Time  `json:"created_at"`
	StartedAt    *time.Time `json:"started_at"`
	EndedAt      *time.Time `json:"ended_at"`
}

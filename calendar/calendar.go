package calendar

import "time"

type Calendar struct {
	ID         string    `json:"id"`
	ResourceID string    `json:"resource_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// Events     []Event   `json:"events"`
}

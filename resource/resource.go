package resource

import "time"

type Resource struct {
	ID        string        `json:"id" gorm:"default:uuid_generate_v4()"`
	OwnID     string        `json:"own_id"`
	Setup     time.Duration `json:"setup"`
	Cleanup   time.Duration `json:"cleanup"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

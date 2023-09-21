package resource

import "time"

type CreateResourceDTO struct {
	OwnID   string        `json:"own_id"`
	Setup   time.Duration `json:"setup"`
	Cleanup time.Duration `json:"cleanup"`
}

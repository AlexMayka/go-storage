package domain

import (
	"time"
)

type Company struct {
	ID          string
	Path        string
	Name        string
	Description string
	CreatedAt   time.Time
}

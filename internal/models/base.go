package models

import "time"

type TimestampedModel struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

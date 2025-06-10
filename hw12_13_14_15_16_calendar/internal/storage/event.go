package storage

import "time"

// Event структура для хранения данных о событии.
type Event struct {
	ID             string        `db:"id" json:"id,omitempty"`
	Title          string        `db:"title" json:"title,omitempty"`
	StartsAt       *time.Time    `db:"starts_at" json:"starts_at,omitempty"`
	EndsAt         *time.Time    `db:"ends_at" json:"ends_at,omitempty"`
	Description    string        `db:"description" json:"description,omitempty"`
	UserID         string        `db:"user_id" json:"user_id,omitempty"`
	NotifyInterval time.Duration `db:"notify_interval" json:"notify_interval,omitempty"`
	Processed      *bool         `db:"processed" json:"processed,omitempty"`
}

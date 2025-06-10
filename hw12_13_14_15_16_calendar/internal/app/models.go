package app

import "time"

// EventNotification структура уведомления о событии.
type EventNotification struct {
	EventID    string    `json:"event_id,omitempty"`
	EventTitle string    `json:"event_title,omitempty"`
	EventDate  time.Time `json:"event_date,omitempty"`
	UserID     string    `json:"user_id,omitempty"`
}

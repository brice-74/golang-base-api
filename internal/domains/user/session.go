package user

import "time"

type Session struct {
	ID            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeactivatedAt time.Time
	IP            string
	Name          string
	Location      string
	UserID        string
}

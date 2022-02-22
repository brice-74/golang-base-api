package user

import "time"

type Session struct {
	ID            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeactivatedAt time.Time
	IP            string
	Agent         string
	UserID        string
}

func (s *Session) IsActive() bool {
	return s.DeactivatedAt.After(time.Now())
}

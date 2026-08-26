package models

import "time"

type PasswordChangeLog struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"index" json:"user_id"`
	ChangedAt time.Time `json:"changed_at"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

func (PasswordChangeLog) TableName() string {
	return "password_change_logs"
}

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type User struct {
	ID              string              `gorm:"primaryKey" json:"clerk_id"`
	Email           string              `gorm:"uniqueIndex" json:"email"`
	FirstName       string              `json:"first_name"`
	LastName        string              `json:"last_name"`
	ProfileImageUrl string              `json:"profile_image_url"`
	PhoneNumber     string              `json:"phone_number"`
	Role            string              `gorm:"default:'user'" json:"role"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	ClerkSyncedAt   *time.Time          `json:"clerk_synced_at"`
	ClerkMetadata   json.RawMessage     `gorm:"type:jsonb" json:"clerk_metadata,omitempty"`

	// Relationships
	SocialLinks []UserSocialLink `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"social_links"`
}

type UserSocialLink struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"index" json:"user_id"`
	Platform  string    `json:"platform"`
	URL       string    `json:"url"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Implement driver.Valuer and sql.Scanner for User to work with gorm
func (u User) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *User) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &u)
}

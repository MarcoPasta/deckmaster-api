package models

import (
	"time"

	"gorm.io/gorm"
)

type DigimonCard struct {
	Id        string `gorm:"primaryKey"`
	Name      string
	CardType  string
	Color     []string `gorm:"serializer:json"`
	Cost      int8
	Dp        int16
	DvCost    []string `gorm:"serializer:json"`
	Form      string
	Attribute string
	Species   []string `gorm:"serializer:json"` // "type" is a primitve, so we have to use species here
	Text      []string `gorm:"serializer:json"`
	Set       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type DigimonCardAltArt struct {
	gorm.Model
	CardID string `gorm:"foreignKey:Id"`
	Image  string
}
type DigimonDecks struct {
	gorm.Model
	Name  string
	Cards []DigimonCard `gorm:"many2many:decks_cards"`
}

// etc..

package models

import "gorm.io/gorm"

type DigimonCard struct {
	gorm.Model
	Name  string
	Color string
}

type DigimonDecks struct {
	gorm.Model
	Name  string
	Cards []DigimonCard `gorm:"many2many:decks_cards"`
}

// etc..

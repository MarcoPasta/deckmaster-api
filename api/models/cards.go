package models

import "gorm.io/gorm"

type DigimonCard struct {
	gorm.Model
	ID    uint `gorm:"primaryKey"`
	Name  string
	Color string
}

type DigimonDecks struct {
	gorm.Model
	Name  string
	cards []DigimonCard
}

// etc..

package models

import (
	"gorm.io/gorm"
)

type DigimonCard struct {
	gorm.Model
	CardId    string
	Name      string
	CardType  string
	Color     []string `gorm:"serializer:json"`
	Cost      int8
	Dp        int16
	DvCost    []string `gorm:"serializer:json"`
	Form      string
	Attribute string
	Species   []string `gorm:"serializer:json"` // "type" is a primitve, so we have to use species here
	CardText  []string `gorm:"serializer:json"`
	Set       string
}

type DigimonDecks struct {
	gorm.Model
	Name  string
	Cards []DigimonCard `gorm:"many2many:decks_cards"`
}

// etc..

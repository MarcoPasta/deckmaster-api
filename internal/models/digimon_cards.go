package models

import (
	"gorm.io/gorm"
)

type DigimonCard struct {
	gorm.Model
	CardNumber string
	Name       string
	CardType   string
	Color      []string `gorm:"serializer:json" collection_format:"multi"`
	Cost       int8
	Dp         int16
	DvCost     []string `gorm:"serializer:json"`
	Form       string
	Attribute  string
	Species    []string `gorm:"serializer:json"` // "type" is a primitve, so we have to use species here
	Text       []string `gorm:"serializer:json"`
	Set        string
}

type DigimonDecks struct {
	gorm.Model
	UserId    string             `json:"userId"`
	Name      string             `json:"name"`
	CoverCard string             `json:"cover"`
	Entries   []DigimonDeckEntry `gorm:"foreignKey:DeckId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"entries"`
}

type DigimonDeckEntry struct {
	gorm.Model
	DeckId uint        `json:"deckId"`
	CardId uint        `json:"cardId"`
	Card   DigimonCard `gorm:"foreignKey:CardId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"card"`
}

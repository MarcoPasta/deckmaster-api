package models

import "gorm.io/gorm"

type NarutoCard struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Cost        int      `json:"cost"`
	Power       int      `json:"power"`
	Keyword     string   `json:"keyword"`
	Group       string   `json:"group"`
	Effect      string   `json:"effect"`
	EffectTypes []string `gorm:"serializer:json" json:"effect_types"`
	Rarity      string   `json:"rarity"`
	CardNumber  string   `json:"card_number"`
	Set         string   `json:"set"`
	Tags        []string `gorm:"serializer:json" json:"tags"`
}

type NarutoDecks struct {
	gorm.Model
	UserId  string            `json:"userId"`
	Name    string            `json:"name"`
	Entries []NarutoDeckEntry `gorm:"foreignKey:DeckId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"entries"`
}

type NarutoDeckEntry struct {
	gorm.Model
	DeckId uint       `json:"deckId"`
	CardId uint       `json:"cardId"`
	Card   NarutoCard `gorm:"foreignKey:CardId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"card"`
}

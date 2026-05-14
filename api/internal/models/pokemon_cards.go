package models

import (
	"gorm.io/gorm"
)

type PokemonCard struct {
	gorm.Model
	CardType    string
	Phase       string
	Name        string
	Kp          int16
	EnergyType  string
	Ability     string
	Attacks     []Attack `gorm:"serializer:json"`
	Weakness    string
	Resistance  string
	RetreatCost int8
	CardNumber  string
	Set         string
	Code        string
	Language    string `json:"language" db:"language"`
}

type Attack struct {
	Cost   []string `json:"cost"`
	Name   string   `json:"name"`
	Effect string   `json:"effect"`
	Damage string   `json:"damage"`
}

type PokemonDecks struct {
	gorm.Model
	UserId  string             `json:"userId"`
	Name    string             `json:"name"`
	Entries []PokemonDeckEntry `gorm:"foreignKey:DeckId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"entries"`
}

type PokemonDeckEntry struct {
	gorm.Model
	DeckId uint        `json:"deckId"`
	CardId uint        `json:"cardId"`
	Card   PokemonCard `gorm:"foreignKey:CardId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"card"`
}

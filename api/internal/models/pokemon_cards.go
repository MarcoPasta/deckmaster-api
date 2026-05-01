package models

import "gorm.io/gorm"

type PokemonCard struct {
	gorm.Model
	DeckId      uint `gorm:"index"` // Foreign key to PokemonDecks
	CardType    string
	Phase       string
	Name        string
	Kp          int16
	EnergyType  string
	Ability     string
	Attacks     []Attack `gorm:"serializer:json"`
	AttacksEN   []Attack `gorm:"serializer:json"`
	Weakness    string
	Resistance  string
	RetreatCost int8
	CardNumber  string
	Set         string
}

type Attack struct {
	Cost   []string `json:"cost"`
	Name   string   `json:"name"`
	Effect string   `json:"effect"`
	Damage string   `json:"damage"`
}

type PokemonDecks struct {
	gorm.Model
	UserId string
	Name   string
	Cards  []PokemonCard `gorm:"foreignKey:DeckId"`
}

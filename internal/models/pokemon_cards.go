package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type PokemonCard struct {
	gorm.Model
	CardType    string
	TrainerType string
	Phase       string
	Name        string
	Kp          int16
	EnergyType  string
	Ability     Ability  `gorm:"serializer:json"`
	Attacks     []Attack `gorm:"serializer:json"`
	Weakness    string
	Resistance  string
	RetreatCost int8
	CardNumber  string
	Set         string
	Code        string
	Tags        []string `gorm:"serializer:json"`
	Language    string   `json:"language" db:"language"`
}

type Ability struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (a *Ability) UnmarshalJSON(data []byte) error {
	// Handle the case where Ability is a plain string in the JSON
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		a.Name = s
		return nil
	}
	// Otherwise, decode as a struct
	type AbilityAlias Ability
	var alias AbilityAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*a = Ability(alias)
	return nil
}

type Attack struct {
	Cost   []string `json:"cost"`
	Name   string   `json:"name"`
	Effect string   `json:"effect"`
	Damage string   `json:"damage"`
}

type PokemonDecks struct {
	gorm.Model
	UserId    string             `json:"userId"`
	Name      string             `json:"name"`
	CoverCard string             `json:"cover"`
	Entries   []PokemonDeckEntry `gorm:"foreignKey:DeckId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"entries"`
}

type PokemonDeckEntry struct {
	gorm.Model
	DeckId uint        `json:"deckId"`
	CardId uint        `json:"cardId"`
	Card   PokemonCard `gorm:"foreignKey:CardId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"card"`
}

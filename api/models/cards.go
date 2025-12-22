package models

import "gorm.io/gorm"

type DigimonCard struct {
	gorm.Model
	ID    uint `gorm:"primaryKey"`
	Name  string
	Color string
}

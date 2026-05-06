package repository

import (
	"context"

	"github.com/MarcoPasta/deckmaster/api/internal/models"
	"gorm.io/gorm"
)

func FindAllNarutoCards(ctx context.Context) ([]models.NarutoCard, error) {
	return gorm.G[models.NarutoCard](DB).Find(ctx)
}

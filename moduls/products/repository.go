package products

import (
	"context"
	"tz/models"
)

type Repository interface {
	GetAll(context.Context) ([]models.Product, error)
	Get(context.Context, int) (models.Product, error)
	Post(context.Context, models.Product) (int, error)
	Delete(context.Context, int) error
	Update(context.Context, models.Product) error
	Find(context.Context, string) (models.Product, error)
}

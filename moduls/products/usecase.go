package products

import (
	"context"
	"tz/models"
)

type Usecase interface {
	GetAll(context.Context) ([]models.Product, error)
	Post(context.Context, models.Product) (models.Product, error)
	Delete(context.Context, int) error
}

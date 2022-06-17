package products

import (
	"context"
	"tz/models"
)

type Repository interface {
	GetAll(context.Context) ([]models.Product, error)
	Post(context.Context)
	Update(context.Context)
	Delete(context.Context)
}

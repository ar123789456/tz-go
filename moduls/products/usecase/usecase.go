package usecase

import (
	"context"
	"tz/models"
	"tz/moduls/products"
)

type Usecase struct {
	repository products.Repository
}

func NewUsecase(repository products.Repository) *Usecase {
	return &Usecase{
		repository: repository,
	}
}

func (uc *Usecase) GetAll(c context.Context) ([]models.Product, error) {
	return uc.repository.GetAll(c)
}

package usecase

import (
	"context"
	"tz/models"
	"tz/moduls/products"
)

type Usecase struct {
	repository products.Repository
}

func NewUsecase(repository products.Repository) Usecase {
	return Usecase{
		repository: repository,
	}
}

func (uc *Usecase) GetAll(c context.Context) ([]models.Product, error) {
	return uc.repository.GetAll(c)
}

func (uc *Usecase) Get(c context.Context, id int) (models.Product, error) {
	return uc.repository.Get(c, id)
}

func (uc *Usecase) Post(c context.Context, p models.Product) (int, error) {
	if p.Price < 0 {
		return 0, products.ErrInvalidPrice
	}
	return uc.repository.Post(c, p)
}

func (uc *Usecase) Delete(c context.Context, i int) error {
	return uc.repository.Delete(c, i)
}

func (uc *Usecase) Update(c context.Context, p models.Product) error {
	if p.Price < 0 {
		return products.ErrInvalidPrice
	}
	return nil
}

func (uc *Usecase) Find(context.Context, string) (models.Product, error) {
	return models.Product{}, nil
}

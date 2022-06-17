package products

import "errors"

var (
	ErrInvalidPrice = errors.New("price must be positive number")
	ErrInvalidName  = errors.New("product with this name already exists")
	ErrTempNotFound = errors.New("template not found")
)

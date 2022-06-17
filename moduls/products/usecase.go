package products

type Usecase interface {
	GetAll()
	Post()
	Update()
	Delete()
}

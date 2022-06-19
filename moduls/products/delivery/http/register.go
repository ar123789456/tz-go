package http

import (
	"tz/moduls/products/repository"
	"tz/moduls/products/usecase"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func RegisterProduct(r *gin.Engine, db *bolt.DB) {
	rep := repository.NewRepository(db)
	u := usecase.NewUsecase(&rep)
	h := NewHandler(&u)

	endpoints := r.Group("")
	{
		endpoints.GET("/products", h.Base)
		endpoints.GET("/products/add", h.Add)
		endpoints.GET("/product/edit/:id", h.Edit)
		endpoints.POST("/q/product-search-by-name", h.Find)
		endpoints.POST("/cmd/add-product", h.Post)
		endpoints.POST("/cmd/edit-product", h.Update)
		endpoints.POST("/cmd/delete-product", h.Delete)
	}
}

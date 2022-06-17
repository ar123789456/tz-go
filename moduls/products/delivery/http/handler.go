package http

import (
	"html/template"
	"net/http"
	"tz/moduls/products"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase products.Usecase
}

func NewHandler(usecase products.Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Base(c *gin.Context) {
	temp, err := template.ParseFiles("base.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, products.ErrTempNotFound)
		return
	}
	data, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if temp.Execute(c.Writer, data) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *Handler) Add(c *gin.Context) {
	temp, err := template.ParseFiles("add.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, products.ErrTempNotFound)
		return
	}
	if temp.Execute(c.Writer, nil) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *Handler) Edit(c *gin.Context) {
	temp, err := template.ParseFiles("edit.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, products.ErrTempNotFound)
		return
	}
	if temp.Execute(c.Writer, nil) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

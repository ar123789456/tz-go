package http

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"tz/models"
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

type allErr struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func (h *Handler) Base(c *gin.Context) {
	temp, err := template.ParseFiles("templates/base.html")
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
	temp, err := template.ParseFiles("templates/add.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, products.ErrTempNotFound)
		return
	}
	if temp.Execute(c.Writer, nil) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *Handler) Edit(c *gin.Context) {
	temp, err := template.ParseFiles("templates/edit.html")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, products.ErrTempNotFound)
		return
	}
	id := c.Param("id")
	intid, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	p, err := h.usecase.Get(c.Request.Context(), intid)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if temp.Execute(c.Writer, p) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

type find struct {
	Name string `json:"searchName"`
}

type out_Product struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type outputFind struct {
	Product *out_Product `json:"product"`
}

func (h *Handler) Find(c *gin.Context) {
	inp := new(find)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := h.usecase.Find(c.Request.Context(), inp.Name)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, &outputFind{Product: toOutProduct(res)})
}

type inputCreate struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type createSuccess struct {
	Status string `json:"status"`
	Id     int    `json:"productId"`
}

func (h *Handler) Post(c *gin.Context) {
	inp := new(inputCreate)
	if err := c.BindJSON(inp); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  "bad request",
		})
		return
	}
	pr, err := strconv.Atoi(inp.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  "invalid price",
		})
		return
	}
	id, err := h.usecase.Post(c.Request.Context(), models.Product{
		Name:  inp.Name,
		Price: pr,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &createSuccess{
		Status: "success",
		Id:     id,
	})
}

type inputUpdate struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

type updateDeleteSuccess struct {
	Status string `json:"status"`
}

func (h *Handler) Update(c *gin.Context) {
	inp := new(inputUpdate)
	if err := c.BindJSON(inp); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  "bad request",
		})
		return
	}
	id, err := strconv.Atoi(inp.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  "bad request",
		})
		return
	}
	pr, err := strconv.Atoi(inp.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  "bad request",
		})
		return
	}

	err = h.usecase.Update(c.Request.Context(), models.Product{
		Id:    id,
		Name:  inp.Name,
		Price: pr,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &updateDeleteSuccess{
		Status: "success",
	})
}

type inputDelete struct {
	Id string `json:"id"`
}

func (h *Handler) Delete(c *gin.Context) {
	inp := new(inputDelete)
	if err := c.BindJSON(inp); err != nil {
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  "bad request",
		})
		return
	}
	id, err := strconv.Atoi(inp.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  "id is not a number",
		})
	}
	err = h.usecase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &allErr{
			Status: "failure",
			Error:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &updateDeleteSuccess{
		Status: "success",
	})
}

func toOutProduct(p models.Product) *out_Product {
	return &out_Product{
		Id:    p.Id,
		Name:  p.Name,
		Price: p.Price,
	}
}

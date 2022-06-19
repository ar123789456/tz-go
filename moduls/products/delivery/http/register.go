package http

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func RegisterProduct(r *gin.RouterGroup, db *bolt.DB)

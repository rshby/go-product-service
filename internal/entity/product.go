package entity

import (
	"context"
	"github.com/gin-gonic/gin"
)

type ProductController interface {
	GetByID(c *gin.Context)
}

type ProductService interface {
	GetByID(ctx context.Context, id int) (string, error)
}

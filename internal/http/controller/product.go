package controller

import (
	"github.com/gin-gonic/gin"
	otel "go-product-service/drivers/tracing"
	"go-product-service/internal/entity"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"strconv"
)

type productController struct {
	productService entity.ProductService
}

// NewProductController is function to create new instance of productController. it implements from interface ProductController
func NewProductController(productService entity.ProductService) entity.ProductController {
	return &productController{
		productService: productService,
	}
}

// GetByID is a method to handle request get product by id
func (p *productController) GetByID(c *gin.Context) {
	ctx, span := otel.Start(c)
	defer span.End()

	traceID := ctx.Value("traceID").(string)

	// get id from parms
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"traceID": traceID,
		})
		return
	}

	// call method GetByID in service
	product, err := p.productService.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
			"traceID": traceID,
		})
		return
	}

	// success
	c.JSON(http.StatusOK, gin.H{
		"message": "success get product",
		"traceID": traceID,
		"data":    product,
	})
}

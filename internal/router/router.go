package router

import (
	"github.com/gin-gonic/gin"
	"go-product-service/internal/entity"
	"go-product-service/internal/http/controller"
	"go-product-service/internal/service"
)

type router struct {
	app               *gin.RouterGroup
	productController entity.ProductController
}

// Route is function to create instance router. then register all endpoints API
func Route(app *gin.RouterGroup) {
	r := &router{
		app: app,
	}

	r.Injection()
	r.ApiV1()
}

func (r *router) Injection() {
	// create service
	productService := service.NewProductService()

	// create controller
	productController := controller.NewProductController(productService)

	r.productController = productController
}

// ApiV1 is method to register all v1 endpoints
func (r *router) ApiV1() {
	v1Group := r.app.Group("/v1")
	{
		productV1Group := v1Group.Group("/product")
		{
			productV1Group.GET("/:id", r.productController.GetByID)
		}
	}
}

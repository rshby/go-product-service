package service

import (
	"context"
	"fmt"
	otel "go-product-service/drivers/tracing"
	"go-product-service/internal/entity"
)

type productService struct {
}

// NewProductService is function to create new instance productService. it implements interface ProductService
func NewProductService() entity.ProductService {
	return &productService{}
}

func (p *productService) GetByID(ctx context.Context, id int) (string, error) {
	ctx, span := otel.Start(ctx)
	defer span.End()

	response := fmt.Sprintf("product with ID %d", id)
	return response, nil
}

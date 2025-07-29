# Creating APIs

This guide demonstrates how to create new APIs in the Go backend template following the clean architecture pattern.

## 📋 Before You Start

Before implementing any API, you should:

1. **Plan Your API Design** - Define your API requirements:
   - Request/response formats
   - Validation rules  
   - HTTP status codes
   - Authentication requirements
   - Error handling patterns

2. **Set Up Database Schema First** - Always create PostgreSQL database migrations before implementing the repository layer:
   - Create migration files using timestamp naming convention (e.g., `20250128120300_create_products.up.sql`)
   - Use PostgreSQL-specific features and data types
   - Define proper indexes, constraints, and relationships
   - Test migrations with both `up` and `down` scripts
   - Follow consistent naming patterns for tables and columns

3. **Understand the Architecture** - This guide follows the clean architecture pattern with clear separation of concerns:
   - **Models** (Data structures)
   - **Repository** (Data access layer)
   - **Service** (Business logic layer) 
   - **Handler** (API layer)

4. **Follow Existing Patterns** - Use existing implementations in the project as reference examples for consistency.

## 🔗 Related Documentation
- Review your project's existing API implementations for patterns and conventions
- Check authentication middleware and error handling implementations
- Study existing models, repositories, and services for consistency

## Quick Start Example: Product API

Let's create a complete Product API to demonstrate the process.

> **💡 Tip**: Before implementing, review existing APIs in your project to understand the expected data structures and validation rules. Study patterns used in authentication, user management, or other existing endpoints.

### 1. Define the Model (Data Layer)

The model should match your planned data structure. Ensure field names, types, and validation rules align with your API requirements.

Create the Product model:

```go
// internal/models/product.go
package models

import (
	"time"
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=2,max=100"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price" validate:"required,min=0"`
	SKU         string    `json:"sku" db:"sku" validate:"required"`
	Stock       int       `json:"stock" db:"stock" validate:"min=0"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

### 2. Create Database Migration (PostgreSQL Schema)

**Before implementing the repository layer**, create the PostgreSQL database schema. This ensures your data access layer has the proper database structure to work with.

Create migration files using timestamp naming convention in your migrations directory:

```sql
-- migrations/20250128120300_create_products.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    sku VARCHAR(50) UNIQUE NOT NULL,
    stock INTEGER DEFAULT 0 CHECK (stock >= 0),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- PostgreSQL-specific indexes for optimal performance
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_name ON products USING gin(to_tsvector('english', name));
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_created_at ON products(created_at DESC);

-- Add trigger for automatic updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_products_updated_at 
    BEFORE UPDATE ON products 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- migrations/20250128120300_create_products.down.sql
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS products;
```

> **💡 PostgreSQL Best Practice**: The migration uses PostgreSQL-specific features like `uuid-ossp` extension, `gin` indexes for text search, and triggers for automatic timestamp updates. Always test both `up` and `down` migrations before implementing the repository layer.

### 3. Create Repository Interface and Implementation (Data Layer)

The repository interface should support all operations you need for your API endpoints (GET, POST, PUT, DELETE).

Define repository interface:

```go
// internal/data/repositories/interfaces.go
package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/yantology/golang_template/internal/models"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetBySKU(ctx context.Context, sku string) (*models.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, limit, offset int) ([]*models.Product, error)
}
```

Implement repository:

```go
// internal/data/repositories/product_repository.go
package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	
	"github.com/google/uuid"
	"github.com/yantology/golang_template/internal/models"
	"github.com/yantology/golang_template/pkg/errors"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	
	query := `
		INSERT INTO products (id, name, description, price, sku, stock, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at`
	
	err := r.db.QueryRowContext(ctx, query,
		product.ID, product.Name, product.Description, product.Price,
		product.SKU, product.Stock, product.IsActive,
		product.CreatedAt, product.UpdatedAt,
	).Scan(&product.CreatedAt, &product.UpdatedAt)
	
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	
	return product, nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, sku, stock, is_active, created_at, updated_at
		FROM products
		WHERE id = $1`
	
	product := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price,
		&product.SKU, &product.Stock, &product.IsActive,
		&product.CreatedAt, &product.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("product not found")
		}
		return nil, errors.NewDatabaseError(err)
	}
	
	return product, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, sku, stock, is_active, created_at, updated_at
		FROM products
		WHERE sku = $1`
	
	product := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price,
		&product.SKU, &product.Stock, &product.IsActive,
		&product.CreatedAt, &product.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("product not found")
		}
		return nil, errors.NewDatabaseError(err)
	}
	
	return product, nil
}

func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	query := `
		SELECT id, name, description, price, sku, stock, is_active, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	defer rows.Close()
	
	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.SKU, &product.Stock, &product.IsActive,
			&product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, errors.NewDatabaseError(err)
		}
		products = append(products, product)
	}
	
	if err = rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	
	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	product.UpdatedAt = time.Now()
	
	query := `
		UPDATE products
		SET name = $2, description = $3, price = $4, sku = $5, stock = $6, is_active = $7, updated_at = $8
		WHERE id = $1
		RETURNING updated_at`
	
	err := r.db.QueryRowContext(ctx, query,
		product.ID, product.Name, product.Description, product.Price,
		product.SKU, product.Stock, product.IsActive, product.UpdatedAt,
	).Scan(&product.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("product not found")
		}
		return nil, errors.NewDatabaseError(err)
	}
	
	return product, nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewDatabaseError(err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return errors.NewNotFoundError("product not found")
	}
	
	return nil
}

func (r *productRepository) Search(ctx context.Context, query string, limit, offset int) ([]*models.Product, error) {
	searchPattern := "%" + query + "%"
	sqlQuery := `
		SELECT id, name, description, price, sku, stock, is_active, created_at, updated_at
		FROM products
		WHERE name ILIKE $1 OR description ILIKE $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := r.db.QueryContext(ctx, sqlQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	defer rows.Close()
	
	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.SKU, &product.Stock, &product.IsActive,
			&product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, errors.NewDatabaseError(err)
		}
		products = append(products, product)
	}
	
	if err = rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	
	return products, nil
}
```

### 4. Create Service Interface and Implementation (Business Layer)

The service layer implements business logic and validation rules for your API. Define validation requirements, error handling, and business rules based on your API requirements.

Define service interface:

```go
// internal/business/services/product_service.go
package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/yantology/golang_template/internal/models"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*models.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProducts(ctx context.Context, req *GetProductsRequest) (*GetProductsResponse, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req *UpdateProductRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	SearchProducts(ctx context.Context, req *SearchProductsRequest) (*SearchProductsResponse, error)
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	SKU         string  `json:"sku" validate:"required"`
	Stock       int     `json:"stock" validate:"min=0"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,min=0"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

type GetProductsRequest struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

type GetProductsResponse struct {
	Products []*models.Product `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

type SearchProductsRequest struct {
	Query string `json:"query" validate:"required,min=1"`
	Page  int    `json:"page" validate:"min=1"`
	Limit int    `json:"limit" validate:"min=1,max=100"`
}

type SearchProductsResponse struct {
	Products []*models.Product `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Query    string            `json:"query"`
}
```

Implement service:

```go
// internal/business/services/product_service_impl.go
package services

import (
	"context"
	"time"
	
	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
	
	"github.com/yantology/golang_template/internal/models"
	"github.com/yantology/golang_template/internal/data/repositories"
	"github.com/yantology/golang_template/internal/pkg/logger"
	"github.com/yantology/golang_template/pkg/errors"
)

type productService struct {
	productRepo repositories.ProductRepository
	validator   *validator.Validate
	logger      logger.Logger
}

func NewProductService(
	productRepo repositories.ProductRepository,
	validator *validator.Validate,
	logger logger.Logger,
) ProductService {
	return &productService{
		productRepo: productRepo,
		validator:   validator,
		logger:      logger,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*models.Product, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		s.logger.WithField("request", req).Error("Validation failed for create product")
		return nil, errors.NewValidationError("validation failed")
	}

	// Check if SKU already exists
	existingProduct, err := s.productRepo.GetBySKU(ctx, req.SKU)
	if err == nil && existingProduct != nil {
		s.logger.WithField("sku", req.SKU).Warn("Attempt to create product with duplicate SKU")
		return nil, errors.NewConflictError("product with this SKU already exists")
	}

	// Create product
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
		Stock:       req.Stock,
		IsActive:    true,
	}

	createdProduct, err := s.productRepo.Create(ctx, product)
	if err != nil {
		s.logger.WithField("product", product).Error("Failed to create product")
		return nil, err
	}

	s.logger.WithField("product_id", createdProduct.ID).Info("Product created successfully")
	return createdProduct, nil
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithField("product_id", id).Error("Failed to get product")
		return nil, err
	}
	return product, nil
}

func (s *productService) GetProducts(ctx context.Context, req *GetProductsRequest) (*GetProductsResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError("validation failed")
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	products, err := s.productRepo.GetAll(ctx, req.Limit, offset)
	if err != nil {
		s.logger.Error("Failed to get products")
		return nil, err
	}

	// In a real implementation, you'd also get the total count
	// For simplicity, we're setting it to the length of products
	total := int64(len(products))

	return &GetProductsResponse{
		Products: products,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req *UpdateProductRequest) (*models.Product, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError("validation failed")
	}

	// Get existing product
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	updatedProduct, err := s.productRepo.Update(ctx, product)
	if err != nil {
		s.logger.WithField("product_id", id).Error("Failed to update product")
		return nil, err
	}

	s.logger.WithField("product_id", id).Info("Product updated successfully")
	return updatedProduct, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	// Check if product exists
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		s.logger.WithField("product_id", id).Error("Failed to delete product")
		return err
	}

	s.logger.WithField("product_id", id).Info("Product deleted successfully")
	return nil
}

func (s *productService) SearchProducts(ctx context.Context, req *SearchProductsRequest) (*SearchProductsResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, errors.NewValidationError("validation failed")
	}

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	products, err := s.productRepo.Search(ctx, req.Query, req.Limit, offset)
	if err != nil {
		s.logger.WithField("query", req.Query).Error("Failed to search products")
		return nil, err
	}

	total := int64(len(products))

	return &SearchProductsResponse{
		Products: products,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
		Query:    req.Query,
	}, nil
}
```

### 5. Create Handler (API Layer)

The handler layer implements the HTTP interface for your API endpoints. Define HTTP methods, status codes, and response formats based on your API requirements.

```go
// internal/api/handlers/product_handler.go
package handlers

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"github.com/yantology/golang_template/internal/business/services"
	"github.com/yantology/golang_template/internal/pkg/logger"
	"github.com/yantology/golang_template/pkg/response"
	"github.com/yantology/golang_template/pkg/errors"
)

type ProductHandler struct {
	productService services.ProductService
	logger         logger.Logger
}

func NewProductHandler(productService services.ProductService, logger logger.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req services.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithField("error", err).Error("Failed to bind JSON")
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithField("error", err).Error("Failed to create product")
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(c, appErr.GetStatusCode(), "Failed to create product", err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to create product", err.Error())
		}
		return
	}

	response.Success(c, http.StatusCreated, "Product created successfully", product)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", "Product ID must be a valid UUID")
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		h.logger.WithField("product_id", id).Error("Failed to get product")
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(c, appErr.GetStatusCode(), "Product not found", err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to get product", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Product retrieved successfully", product)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &services.GetProductsRequest{
		Page:  page,
		Limit: limit,
	}

	products, err := h.productService.GetProducts(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to get products")
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(c, appErr.GetStatusCode(), "Failed to get products", err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to get products", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Products retrieved successfully", products)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", "Product ID must be a valid UUID")
		return
	}

	var req services.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithField("error", err).Error("Failed to bind JSON")
		response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.WithField("product_id", id).Error("Failed to update product")
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(c, appErr.GetStatusCode(), "Failed to update product", err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to update product", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Product updated successfully", product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID", "Product ID must be a valid UUID")
		return
	}

	err = h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		h.logger.WithField("product_id", id).Error("Failed to delete product")
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(c, appErr.GetStatusCode(), "Failed to delete product", err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to delete product", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}

func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.Error(c, http.StatusBadRequest, "Missing search query", "Query parameter 'q' is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	req := &services.SearchProductsRequest{
		Query: query,
		Page:  page,
		Limit: limit,
	}

	products, err := h.productService.SearchProducts(c.Request.Context(), req)
	if err != nil {
		h.logger.WithField("query", query).Error("Failed to search products")
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(c, appErr.GetStatusCode(), "Failed to search products", err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, "Failed to search products", err.Error())
		}
		return
	}

	response.Success(c, http.StatusOK, "Products search completed", products)
}
```

### 6. Register Routes

Register routes that match your planned URL patterns and HTTP methods.

```go
// internal/api/routes/product_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yantology/golang_template/internal/api/handlers"
)

func RegisterProductRoutes(router *gin.RouterGroup, productHandler *handlers.ProductHandler) {
	products := router.Group("/products")
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.GetProducts)
		products.GET("/search", productHandler.SearchProducts)
		products.GET("/:id", productHandler.GetProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}
}
```

Update main routes file:

```go
// internal/api/routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yantology/golang_template/internal/api/handlers"
	"github.com/yantology/golang_template/internal/api/middleware"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler, // Add this
	authMiddleware middleware.AuthMiddleware,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		RegisterAuthRoutes(v1, authHandler)

		// Protected routes
		protected := v1.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			RegisterUserRoutes(protected, userHandler)
			RegisterProductRoutes(protected, productHandler) // Add this
		}
	}
}
```

### 7. Update Dependency Injection

Add to your DI container:

```go
// internal/container/container.go
func (c *Container) GetProductRepository() repositories.ProductRepository {
	return repositories.NewProductRepository(c.db)
}

func (c *Container) GetProductService() services.ProductService {
	return services.NewProductService(
		c.GetProductRepository(),
		c.validator,
		c.logger,
	)
}

func (c *Container) GetProductHandler() *handlers.ProductHandler {
	return handlers.NewProductHandler(
		c.GetProductService(),
		c.logger,
	)
}
```

Update route setup in main:

```go
// cmd/main.go
productHandler := container.GetProductHandler()

routes.SetupRoutes(
	router,
	authHandler,
	userHandler,
	productHandler, // Add this
	authMiddleware,
)
```


## API Testing

Test your new API using your planned request formats and expected responses:

> **🧪 Testing Guide**: Create comprehensive test scenarios for each endpoint to validate your implementation works as expected.

```bash
# Create a product
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "sku": "LAPTOP-001",
    "stock": 10
  }'

# Get all products
curl -X GET "http://localhost:8080/api/v1/products?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get a specific product
curl -X GET http://localhost:8080/api/v1/products/PRODUCT_UUID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Search products
curl -X GET "http://localhost:8080/api/v1/products/search?q=laptop&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Update a product
curl -X PUT http://localhost:8080/api/v1/products/PRODUCT_UUID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Updated Laptop",
    "price": 899.99,
    "stock": 15
  }'

# Delete a product
curl -X DELETE http://localhost:8080/api/v1/products/PRODUCT_UUID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Summary

This example demonstrates the complete flow of creating an API in the Go backend template:

1. **Model Definition**: Define the data structure with validation
2. **Database Migration**: Create PostgreSQL schema with proper indexes and constraints
3. **Repository Layer**: Handle data persistence with interface and implementation
4. **Service Layer**: Implement business logic with validation and error handling
5. **Handler Layer**: Handle HTTP requests and responses
6. **Route Registration**: Define API endpoints
7. **Dependency Injection**: Wire everything together

The implementation uses the existing project packages:
- **Logger**: `internal/pkg/logger.Logger` interface for structured logging
- **Errors**: `pkg/errors.AppError` for consistent error handling with HTTP status codes
- **Response**: `pkg/response.Success` and `pkg/response.Error` for standardized API responses

Follow this pattern for any new API you want to create in your application.

## 🔄 Development Workflow

When implementing any API endpoint:

1. **Plan Your API** - Define your API requirements and design
2. **Define the Model** - Create the data structure that matches your API design
3. **Create Database Schema** - Set up PostgreSQL migrations before implementing data access
4. **Implement Bottom-Up** - Follow the layers: Model → Database Migration → Repository → Service → Handler → Routes
5. **Validate Implementation** - Ensure your implementation matches your planned request/response formats
6. **Test Thoroughly** - Create comprehensive test cases for all endpoints
7. **Update Documentation** - Keep your documentation current with implementation changes

## 📚 Cross-Reference Guide

| Implementation Layer | Planning Considerations |
|---------------------|-------------------------|
| **Models** | Data structures based on your API requirements |
| **Database Migration** | Data relationships and constraints for your use case |
| **Repository** | Database operations needed for your endpoints |
| **Service** | Business rules and validation logic |
| **Handler** | Request/response formats and HTTP status codes |
| **Routes** | Endpoint URLs and HTTP methods |
| **Testing** | Test cases covering all scenarios |

For specific examples, study existing implementations in your project to maintain consistency across your codebase.


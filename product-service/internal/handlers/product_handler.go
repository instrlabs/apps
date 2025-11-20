package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/instrlabs/product-service/internal/helpers"
	"github.com/instrlabs/product-service/internal/services"
	"github.com/instrlabs/product-service/internal/validators"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	productService services.ProductServiceInterface
	validator      *validators.RequestValidator
	responseHelper *helpers.ResponseHelper
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(productService services.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		validator:      validators.NewRequestValidator(),
		responseHelper: helpers.NewResponseHelper(),
	}
}

// ListProducts handles GET /products requests
func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	// Validate request parameters
	req, err := h.validator.ValidateProductListRequest(c)
	if err != nil {
		return h.responseHelper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request parameters", err.Error())
	}

	// Call service layer
	result, err := h.productService.ListProducts(req.Type, req.Page, req.Limit)
	if err != nil {
		return h.responseHelper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch products", err.Error())
	}

	// Send successful response
	return h.responseHelper.SuccessResponseWithPagination(c, "Success",
		fiber.Map{"products": result.Products}, result.Pagination)
}

// ListProductsByType handles GET /:type requests
func (h *ProductHandler) ListProductsByType(c *fiber.Ctx) error {
	// Get type from route parameter
	productType := c.Params("type")

	// Validate request parameters
	req, err := h.validator.ValidateProductListRequestWithType(c, productType)
	if err != nil {
		return h.responseHelper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request parameters", err.Error())
	}

	// Call service layer
	result, err := h.productService.ListProducts(req.Type, req.Page, req.Limit)
	if err != nil {
		return h.responseHelper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch products", err.Error())
	}

	// Send successful response
	return h.responseHelper.SuccessResponseWithPagination(c, "Success",
		fiber.Map{"products": result.Products}, result.Pagination)
}

// GetProductByID handles GET /:type/:id requests
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	// Validate request parameters
	req, err := h.validator.ValidateGetProductRequest(c)
	if err != nil {
		return h.responseHelper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request parameters", err.Error())
	}

	// Lookup by product key with type filter from URL
	product, err := h.productService.GetProductByKey(req.Identifier, req.LookupType)
	if err != nil {
		if err.Error() == "product not found" {
			return h.responseHelper.ErrorResponse(c, fiber.StatusNotFound, "Product not found",
				"Product with the specified identifier not found")
		}
		return h.responseHelper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch product", err.Error())
	}

	// Send successful response
	return h.responseHelper.SuccessResponse(c, "Success", fiber.Map{"product": product})
}

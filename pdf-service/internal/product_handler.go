package internal

import (
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productRepo *ProductRepository
}

func NewProductHandler(productRepo *ProductRepository) *ProductHandler {
	return &ProductHandler{
		productRepo: productRepo,
	}
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	products, err := h.productRepo.List("pdf")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch products",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
		"errors":  nil,
		"data":    products,
	})
}

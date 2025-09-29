package internal

import (
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	repo *ProductRepository
}

func NewProductHandler(repo *ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	products, err := h.repo.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  nil,
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Products retrieved successfully",
		"errors":  nil,
		"data":    products,
	})
}

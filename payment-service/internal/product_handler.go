package internal

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	repo *ProductRepository
}

func NewProductHandler(repo *ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var p Product
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  nil,
			"data":    nil,
		})
	}
	// validate required fields: name, key, price, active, isFree
	if p.Name == "" || p.Key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing required fields: name and key are required",
			"errors":  nil,
			"data":    nil,
		})
	}
	// Set userId from authenticated user if available (override any client-sent value)
	if userID, ok := c.Locals("UserID").(string); ok && userID != "" {
		p.UserID = userID
	}
	if err := h.repo.CreateProduct(c.Context(), &p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  nil,
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"errors":  nil,
		"data":    p,
	})
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	onlyActive := c.QueryBool("active", false)
	products, err := h.repo.ListProducts(c.Context(), onlyActive)
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

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	var product *Product

	if _, err := primitive.ObjectIDFromHex(id); err == nil {
		product, err = h.repo.GetProductByID(c.Context(), id)
	} else {
		product, err = h.repo.GetProductByKey(c.Context(), id)
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
			"errors":  nil,
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product retrieved successfully",
		"errors":  nil,
		"data":    product,
	})
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  nil,
			"data":    nil,
		})
	}
	allowed := map[string]bool{
		"name":        true,
		"key":         true,
		"price":       true,
		"description": true,
		"image":       true,
		"productType": true,
		"active":      true,
		"isFree":      true,
	}
	updateFields := bson.M{}
	for k, v := range body {
		if allowed[k] {
			updateFields[k] = v
		}
	}
	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No valid fields to update",
			"errors":  nil,
			"data":    nil,
		})
	}
	if err := h.repo.UpdateProduct(c.Context(), id, updateFields); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  nil,
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
		"errors":  nil,
		"data":    nil,
	})
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.DeleteProduct(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  nil,
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
		"errors":  nil,
		"data":    nil,
	})
}

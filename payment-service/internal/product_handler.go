package internal

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := h.repo.CreateProduct(c.Context(), &p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(p)
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	onlyActive := c.QueryBool("active", false)
	products, err := h.repo.ListProducts(c.Context(), onlyActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.repo.GetProductByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}
	return c.JSON(product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	allowed := map[string]bool{
		"name":        true,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no valid fields to update"})
	}
	if err := h.repo.UpdateProduct(c.Context(), id, updateFields); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "updated"})
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.DeleteProduct(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

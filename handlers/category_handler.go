package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gomed/database"
	"gomed/models"
)

func GetAllCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to load categories",
		})
	}
	return c.JSON(categories)
}
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gomed/database"
	"gomed/models"
)

func CreateCart(c *fiber.Ctx) error {
    cart := models.Cart{}
    if err := c.BodyParser(&cart); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    database.DB.Create(&cart)
    return c.Status(201).JSON(cart)
}

func CreateOrder(c *fiber.Ctx) error {
    order := models.Order{}
    if err := c.BodyParser(&order); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    // คำนวณยอดรวม
    var cart models.Cart
    if err := database.DB.First(&cart, order.CartID).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Cart not found"})
    }

    database.DB.Create(&order)
    return c.Status(201).JSON(order)
}

func AddOrderItem(c *fiber.Ctx) error {
    item := models.OrderItem{}
    if err := c.BodyParser(&item); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    var medicine models.Medicine
    if err := database.DB.First(&medicine, item.MedicineID).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Medicine not found"})
    }

    item.Price = medicine.Price
    database.DB.Create(&item)
    return c.Status(201).JSON(item)
}
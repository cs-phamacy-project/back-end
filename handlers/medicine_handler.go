package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gomed/database"
	"gomed/models"
)

// CreateMedicine - เพิ่มยาใหม่ลงในระบบ
func CreateMedicine(c *fiber.Ctx) error {
    var medicine models.Medicine

    // ✅ ใช้ JSON Bind แบบถูกต้อง
    if err := c.BodyParser(&medicine); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    // ✅ ตรวจสอบค่า PharmacyID
    if medicine.PharmacyID == 0 {
        return c.Status(400).JSON(fiber.Map{"error": "PharmacyID is required"})
    }

    // ✅ เช็คว่ามี Pharmacy นี้อยู่หรือไม่
    var pharmacy models.Pharmacy
    if err := database.DB.First(&pharmacy, medicine.PharmacyID).Error; err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Pharmacy not found"})
    }

    // ✅ สร้างข้อมูลใหม่
    if err := database.DB.Create(&medicine).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create medicine"})
    }

    return c.Status(201).JSON(medicine)
}

func GetMedicineByID(c *fiber.Ctx) error {
	medicineID := c.Params("id") // รับค่า ID จาก URL

	var medicine models.Medicine
	result := database.DB.First(&medicine, medicineID)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Medicine not found"})
	}

	return c.JSON(medicine)
}


package handlers

import (
	"fmt"
	"gomed/database"
	"gomed/models"

	"github.com/gofiber/fiber/v2"
)

func GetPharmacies(c *fiber.Ctx) error {
    var pharmacies []models.Pharmacy

    result := database.DB.
        Preload("Medicines").  
        Find(&pharmacies)

    if result.Error != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to fetch pharmacies",
        })
    }

    return c.JSON(pharmacies)
}


func GetMedicinesByPharmacy(c *fiber.Ctx) error {
	pharmacyID := c.Params("id") // รับ Pharmacy ID จาก URL

	var medicines []models.Medicine
	result := database.DB.Where("pharmacy_id = ?", pharmacyID).Find(&medicines)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch medicines"})
	}

	return c.JSON(medicines)
}

func GetPharmacyByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var pharmacy models.Pharmacy

    result := database.DB.Preload("Medicines").First(&pharmacy, id)
    if result.Error != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Pharmacy not found"})
    }

    return c.JSON(pharmacy)
}

func ApprovePharmacy(c *fiber.Ctx) error {
    id := c.Params("id") // รับ ID ของ pharmacy จาก URL

    var pharmacy models.Pharmacy
    if err := database.DB.First(&pharmacy, id).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Pharmacy not found"})
    }

    // ✅ อัปเดต Status เป็น "approve"
    pharmacy.Status = "approve"
    if err := database.DB.Save(&pharmacy).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update status"})
    }

    return c.JSON(fiber.Map{"message": "Pharmacy approved successfully", "pharmacy": pharmacy})
}

func GetUnapprovedPharmacies(c *fiber.Ctx) error {
	fmt.Println("🔍 Called GetUnapprovedPharmacies") // <-- เพิ่มตรงนี้

    var pharmacies []models.Pharmacy

    result := database.DB.
        Preload("Medicines"). 
        Where("status = ?", "unapprove"). 
        Find(&pharmacies)

    if result.Error != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to fetch pharmacies",
        })
    }

	return c.JSON(pharmacies)
}


	
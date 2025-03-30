package handlers

import (

	"gomed/database"
	"gomed/models"
	"strconv"


	"github.com/gofiber/fiber/v2"
)

func CreateMedicine(c *fiber.Ctx) error {
    // รับค่าทีละ field ด้วย FormValue แทน BodyParser
    productName := c.FormValue("product_name")
    description := c.FormValue("description")
    categoryID, _ := strconv.Atoi(c.FormValue("category_id"))
    pharmacyID, _ := strconv.Atoi(c.FormValue("pharmacy_id"))
    price, _ := strconv.ParseFloat(c.FormValue("price"), 64)
    stock, _ := strconv.Atoi(c.FormValue("stock"))
    expiredDate := c.FormValue("expired_date")
    fda := c.FormValue("fda")
    status := c.FormValue("status")

    // รับไฟล์ภาพ
    file, err := c.FormFile("image")
    imagePath := ""
    if err == nil && file != nil {
        // เซฟไฟล์ไปในโฟลเดอร์ uploads/
        imagePath = "uploads/" + file.Filename
        c.SaveFile(file, imagePath)
    }

    // เช็ก Pharmacy ID
    if pharmacyID == 0 {
        return c.Status(400).JSON(fiber.Map{"error": "PharmacyID is required"})
    }

    // เช็กว่า pharmacy มีอยู่หรือไม่
    var pharmacy models.Pharmacy
    if err := database.DB.First(&pharmacy, pharmacyID).Error; err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Pharmacy not found"})
    }

    // สร้าง medicine
    medicine := models.Medicine{
        ProductName: productName,
        Description: description,
        CategoryID:  uint(categoryID),
        PharmacyID:  uint(pharmacyID),
        Price:       price,
        Stock:       stock,
        ExpiredDate: expiredDate,
        Fda:         fda,
        Status:      status,
        Image:       imagePath,
    }

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

func UpdateMedicine(c *fiber.Ctx) error {
    id := c.Params("id")

    var medicine models.Medicine
    if err := database.DB.First(&medicine, id).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Medicine not found"})
    }

    var updatedData models.Medicine
    if err := c.BodyParser(&updatedData); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON body"})
    }

    // Assign fields you allow to update
    medicine.ProductName = updatedData.ProductName
    medicine.Description = updatedData.Description
    medicine.Fda = updatedData.Fda
    medicine.ExpiredDate = updatedData.ExpiredDate
    medicine.Status = updatedData.Status
    medicine.Price = updatedData.Price
    medicine.Stock = updatedData.Stock
    medicine.CategoryID = updatedData.CategoryID

    // ไม่ต้องจัดการ image ตรงนี้ ถ้าไม่ได้อัพโหลด
    // ถ้าอยากอัพก็อาจต้องเพิ่ม logic upload image แบบ base64 หรือ URL ก็ได้

    if err := database.DB.Save(&medicine).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update medicine"})
    }

    return c.JSON(medicine)
}


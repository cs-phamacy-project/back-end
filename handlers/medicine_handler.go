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

    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid form data"})
    }

    data := form.Value

    medicine.ProductName = data["product_name"][0]
    medicine.Description = data["description"][0]
    medicine.Fda = data["fda"][0]
    medicine.ExpiredDate = data["expired_date"][0]
    medicine.Status = data["status"][0]

    price, _ := strconv.ParseFloat(data["price"][0], 64)
    stock, _ := strconv.Atoi(data["stock"][0])
    categoryID, _ := strconv.Atoi(data["category_id"][0])

    medicine.Price = price
    medicine.Stock = stock
    medicine.CategoryID = uint(categoryID)

    files := form.File["image"]
    if len(files) > 0 {
        file := files[0]
        // สำหรับเทสต์: สมมติใช้ชื่อไฟล์เก็บไว้เฉย ๆ
        medicine.Image = file.Filename

        // ❌ ข้ามการ save ไฟล์
        // path := fmt.Sprintf("./uploads/%s", file.Filename)
        // c.SaveFile(file, path)
    }

    if err := database.DB.Save(&medicine).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update medicine"})
    }

    return c.JSON(medicine)
}


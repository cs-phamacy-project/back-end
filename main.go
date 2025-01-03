package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// RegistrationRequest โครงสร้างข้อมูลสำหรับการสมัครสมาชิก
type RegistrationRequest struct {
	Email             string `json:"email"`
	ConfirmEmail      string `json:"confirmEmail"`
	Password          string `json:"password"`
	ConfirmPassword   string `json:"confirmPassword"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Gender            string `json:"gender"`
	DateOfBirth       string `json:"dateOfBirth"`
	Phone             string `json:"phone"`
	ChronicDisease    string `json:"chronicDisease"`
	MedicationAllergy string `json:"medicationAllergy"`
}

func main() {
	app := fiber.New()

	// เปิดใช้ CORS เพื่อให้เชื่อมต่อกับ Frontend ได้
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // URL ของ Frontend
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// เส้นทางสำหรับสมัครสมาชิก
	app.Post("/register", registerUser)

	// รันเซิร์ฟเวอร์
	log.Fatal(app.Listen(":3001"))
}

func registerUser(c *fiber.Ctx) error {
	var data RegistrationRequest

	// แปลง JSON จาก Body Request
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// ตรวจสอบข้อมูลพื้นฐาน
	if data.Email != data.ConfirmEmail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Emails do not match",
		})
	}

	if data.Password != data.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully!",
		"data":    data,
	})
}

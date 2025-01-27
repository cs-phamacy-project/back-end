package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// RegistrationRequest โครงสร้างข้อมูลสำหรับการสมัครสมาชิก
type RegistrationRequest struct {
	Email             string            `json:"email"`
	ConfirmEmail      string            `json:"confirmEmail"`
	Password          string            `json:"password"`
	ConfirmPassword   string            `json:"confirmPassword"`
	FirstName         string            `json:"firstName"`
	LastName          string            `json:"lastName"`
	Gender            string            `json:"gender"`
	Phone             string            `json:"phone"`
	Address           Address           `json:"address"`
	HealthInformation HealthInformation `json:"healthInformation"`
}

// Address โครงสร้างข้อมูลสำหรับที่อยู่
type Address struct {
	SubDistrict string `json:"subDistrict"`
	District    string `json:"district"`
	Province    string `json:"province"`
	ZipCode     string `json:"zipCode"`
}

// HealthInformation โครงสร้างข้อมูลสำหรับข้อมูลสุขภาพ
type HealthInformation struct {
	ChronicDisease    string `json:"chronicDisease"`
	MedicationAllergy string `json:"medicationAllergy"`
}

// User โครงสร้างข้อมูลที่บันทึกลงในฐานข้อมูล
type User struct {
	ID                uint   `gorm:"primaryKey"`
	Email             string `gorm:"unique;not null"`
	Password          string `gorm:"not null"`
	FirstName         string `gorm:"not null"`
	LastName          string `gorm:"not null"`
	Gender            string `gorm:"not null"`
	Phone             string `gorm:"not null"`
	SubDistrict       string
	District          string
	Province          string
	ZipCode           string
	ChronicDisease    string
	MedicationAllergy string
}

// initDatabase ฟังก์ชันสำหรับเชื่อมต่อฐานข้อมูล
func initDatabase() {
	var err error
	dsn := "host=postgres user=postgres password=01022546zazakoeiei dbname=gomedDB port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// สร้างตารางในฐานข้อมูลอัตโนมัติ
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

// registerUser ฟังก์ชันสำหรับจัดการการลงทะเบียน
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

	// เข้ารหัสรหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// บันทึกข้อมูลลงฐานข้อมูล
	user := User{
		Email:             data.Email,
		Password:          string(hashedPassword),
		FirstName:         data.FirstName,
		LastName:          data.LastName,
		Gender:            data.Gender,
		Phone:             data.Phone,
		SubDistrict:       data.Address.SubDistrict,
		District:          data.Address.District,
		Province:          data.Address.Province,
		ZipCode:           data.Address.ZipCode,
		ChronicDisease:    data.HealthInformation.ChronicDisease,
		MedicationAllergy: data.HealthInformation.MedicationAllergy,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save user to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully!",
		"user":    user,
	})
}

// loginUser ฟังก์ชันสำหรับตรวจสอบข้อมูลการเข้าสู่ระบบ
func loginUser(c *fiber.Ctx) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// แปลง JSON จาก Body Request
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// ค้นหาผู้ใช้ในฐานข้อมูล
	var user User
	if err := db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
	})
}

func main() {
	// เริ่มการเชื่อมต่อกับฐานข้อมูล
	initDatabase()

	app := fiber.New()

	// เปิดใช้ CORS เพื่อให้เชื่อมต่อกับ Frontend ได้
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // URL ของ Frontend
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// เส้นทางสำหรับสมัครสมาชิก
	app.Post("/register", registerUser)

	// เส้นทางสำหรับเข้าสู่ระบบ
	app.Post("/login", loginUser)

	// รันเซิร์ฟเวอร์
	log.Fatal(app.Listen(":3001"))
}

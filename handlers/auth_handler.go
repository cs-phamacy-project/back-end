package handlers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gomed/database"
	"gomed/models"
)

type RegistrationRequest struct {
	Email             string `json:"email"`
	ConfirmEmail      string `json:"confirmEmail"`
	Password          string `json:"password"`
	ConfirmPassword   string `json:"confirmPassword"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Phone             string `json:"phone"`
	AddressDescription string `json:"addressDescription"`
	SubDistrict        string `json:"subDistrict"`
	District           string `json:"district"`
	Province           string `json:"province"`
	ZipCode            string `json:"zipCode"`
	ChronicDisease    string `json:"chronicDisease"`
	MedicationAllergy string `json:"medicationAllergy"`
}

type RegistrationPharmacyRequest struct {
	Email           string `json:"email"`
	ConfirmEmail    string `json:"confirmEmail"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Phone           string `json:"phone"`
	Certificate     string `json:"certificate"`
	StoreImg        string `json:"storeImg"`
	AddressDescription string `json:"addressDescription"`
	SubDistrict        string `json:"subDistrict"`
	District           string `json:"district"`
	Province           string `json:"province"`
	ZipCode            string `json:"zipCode"`
}

func RegisterUser(c *fiber.Ctx) error {
	var data RegistrationRequest

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if data.Email != data.ConfirmEmail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Emails do not match"})
	}

	if data.Password != data.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Passwords do not match"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := models.User{
		Email:             data.Email,
		Password:          string(hashedPassword),
		FirstName:         data.FirstName,
		LastName:          data.LastName,
		Phone:             data.Phone,
		SubDistrict:       data.SubDistrict,
		District:          data.District,
		Province:          data.Province,
		ZipCode:           data.ZipCode,
		ChronicDisease:    data.ChronicDisease,
		MedicationAllergy: data.MedicationAllergy,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save user to database"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User registered successfully!", "user": user})
}

func RegisterPharmacy(c *fiber.Ctx) error {
    var data RegistrationPharmacyRequest
    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // ตรวจสอบไฟล์ certificate
    certificateFile, err := c.FormFile("certificate")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certificate file is required"})
    }

    // ตรวจสอบไฟล์ storeImg
    storeImgFile, err := c.FormFile("storeImg")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Store image file is required"})
    }

    // ตรวจสอบประเภทไฟล์ (optional) ถ้าต้องการ
    if certificateFile != nil {
        if certificateFile.Size == 0 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Certificate file is empty"})
        }
    }

    if storeImgFile != nil {
        if storeImgFile.Size == 0 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Store image file is empty"})
        }
    }

    // หลีกเลี่ยงการบันทึกไฟล์
    // แค่ตรวจสอบข้อมูลการสมัคร
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
    }

    var existing models.Pharmacy
    if err := database.DB.Where("email = ?", data.Email).First(&existing).Error; err == nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already registered"})
    }

    pharmacy := models.Pharmacy{
        Email:              data.Email,
        Password:           string(hashedPassword),
        FirstName:          data.FirstName,
        LastName:           data.LastName,
        Phone:              data.Phone,
        Certificate:        certificateFile.Filename, // บันทึกชื่อไฟล์ (แต่ไม่บันทึกไฟล์จริง)
        StoreImg:           storeImgFile.Filename,     // บันทึกชื่อไฟล์ (แต่ไม่บันทึกไฟล์จริง)
        AddressDescription: data.AddressDescription,
        SubDistrict:        data.SubDistrict,
        District:           data.District,
        Province:           data.Province,
        ZipCode:            data.ZipCode,
    }

    // เก็บข้อมูลในฐานข้อมูล
    if err := database.DB.Create(&pharmacy).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save pharmacy to database"})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Pharmacy registered successfully!",
        "pharmacy": pharmacy,
    })
}


func LoginUser(c *fiber.Ctx) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var user models.User
	if err := database.DB.Where("email = ?", data.Email).First(&user).Error; err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Login successful",
			"user":    user,
			"role":    "user",
			"user_id": user.ID,
		})
	}

	if err := database.DB.Where("email = ?", data.Email).First(&user).Error; err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Login successful",
			"user":    user,
			"role":    "admin",
		})
	}



	var pharmacy models.Pharmacy
	if err := database.DB.Where("email = ?", data.Email).First(&pharmacy).Error; err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(pharmacy.Password), []byte(data.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":  "Login successful",
			"user":     pharmacy,
			"role":     "staff",
			"pharmacy_id": pharmacy.ID,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email not found"})
}


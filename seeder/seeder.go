package seeder

import (
    "errors" 
    "gorm.io/gorm" 
	"gomed/database"
	"gomed/models"
)

func SeedCategories() {
    categories := []models.Category{
        {Name: "ยา"},
        {Name: "อุปกรณ์ทางการแพทย์"},		
        {Name: "เวชสำอาง"},
		{Name: "อาหารเสริม"},
    }

    for _, cat := range categories {
        var existing models.Category
        err := database.DB.Where("name = ?", cat.Name).First(&existing).Error
        if errors.Is(err, gorm.ErrRecordNotFound) {
            database.DB.Create(&cat)
        }
    }
}

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"

// 	"github.com/bxcodec/faker/v4"
// 	"gomed/database"
// )

// // SeedUsers สร้างข้อมูลปลอมสำหรับ Users
// func SeedUsers() {
// 	for i := 0; i < 10; i++ {
// 		user := User{
// 			Email:             faker.Email(),
// 			Password:          "password123",
// 			FirstName:         faker.FirstName(),
// 			LastName:          faker.LastName(),
// 			Phone:             faker.Phonenumber(),
// 			AddressDescription: faker.Word(),
// 			SubDistrict:       faker.Word(),
// 			District:          faker.Word(),
// 			Province:          faker.Word(),
// 			ZipCode:           "12110",
// 			ChronicDisease:    faker.Word(),
// 			MedicationAllergy: faker.Word(),
// 		}
// 		database.DB.Create(&user)
// 	}
// 	fmt.Println("✅ Seeded 10 Users")
// }

// // SeedPharmacies สร้างข้อมูลปลอมสำหรับ Pharmacies
// func SeedPharmacies() {
// 	for i := 0; i < 10; i++ {
// 		pharmacy := Pharmacy{
// 			Email:              faker.Email(),
// 			Password:           "password123",
// 			FirstName:          faker.FirstName(),
// 			LastName:           faker.LastName(),
// 			Phone:              faker.Phonenumber(),
// 			Certificate:        faker.Word(),
// 			StoreImg:           faker.URL(),
// 			AddressDescription: faker.Word(),
// 			SubDistrict:        faker.Word(),
// 			District:           faker.Word(),
// 			Province:           faker.Word(),
// 			ZipCode:            "12110",
// 			Contact:            faker.Phonenumber(),
// 			Status:             "unapprove",
// 			Role:               "staff",
// 		}
// 		database.DB.Create(&pharmacy)
// 	}
// 	fmt.Println("✅ Seeded 10 Pharmacies")
// }

// // SeedMedicines สร้างข้อมูลปลอมสำหรับ Medicines
// func SeedMedicines() {
// 	var pharmacies []Pharmacy
// 	database.DB.Find(&pharmacies) // ดึง Pharmacy ทั้งหมด
// 	if len(pharmacies) == 0 {
// 		log.Println("⚠️ ไม่มี Pharmacy ใน Database, ข้ามการสร้าง Medicines")
// 		return
// 	}

// 	for i := 0; i < 10; i++ {
// 		medicine := Medicine{
// 			PharmacyID:  pharmacies[rand.Intn(len(pharmacies))].ID, // เลือก Pharmacy แบบสุ่ม
// 			Quantity:    rand.Intn(100) + 1,
// 			Name:        faker.Word(),
// 			Price:       rand.Float64() * 100,
// 			Description: faker.Sentence(),
// 			ExpiredDate: time.Now().AddDate(0, rand.Intn(12), 0).Format("2006-01-02"),
// 			FDA:         faker.Word(),
// 			CategoryID:  rand.Intn(5) + 1, // ให้ CategoryID เป็นค่าระหว่าง 1-5
// 		}
// 		database.DB.Create(&medicine)
// 	}
// 	fmt.Println("✅ Seeded 10 Medicines")
// }

// // SeedAll เรียกใช้ฟังก์ชัน Seed ทั้งหมด
// func SeedAll() {
// 	SeedUsers()
// 	SeedPharmacies()
// 	SeedMedicines()
// }

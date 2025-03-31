package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gomed/database"
	"gomed/models"
	"gomed/routes"
	"gomed/seeder"
	"gomed/handlers"
	
)

func main() {
	// เชื่อมต่อฐานข้อมูล
	database.InitDatabase()

	// AutoMigrate
	models.MigrateDB(database.DB)
	handlers.SetupWebSocketServer()
	seeder.SeedCategories()
	// models.SeedAll()

	app := fiber.New()

	// เปิดใช้ CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// API Routes
	routes.AuthRoutes(app)
	routes.SetupOrderRoutes(app)
	routes.SetupPaymentRoutes(app)
	routes.SetupChatRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend Running")
	})

	// รันเซิร์ฟเวอร์
	log.Fatal(app.Listen(":3001"))
}

package routes

import (
	"github.com/gofiber/fiber/v2"
	"gomed/handlers"
)

func AuthRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/register", handlers.RegisterUser)
	api.Post("/login", handlers.LoginUser)
	api.Post("/pharmacyRegister", handlers.RegisterPharmacy)
	api.Get("/pharmacies", handlers.GetPharmacies)
	app.Get("/api/pharmacies/:id/medicines", handlers.GetMedicinesByPharmacy)
	app.Get("/api/pharmacies/:id", handlers.GetPharmacyByID)
	api.Post("/medicines", handlers.CreateMedicine)
	api.Put("/:id/approve", handlers.ApprovePharmacy)
	api.Post("/carts", handlers.CreateCart)
	api.Post("/orders", handlers.CreateOrder)
	api.Post("/order-items", handlers.AddOrderItem)
	api.Get("/medicines/:id", handlers.GetMedicineByID)
	api.Get("/pharmacies/unapproved", handlers.GetUnapprovedPharmacies)
	api.Get("/categories", handlers.GetAllCategories)
	app.Put("/api/medicines/:id", handlers.UpdateMedicine)

}

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
	api.Get("/medicines/:id", handlers.GetMedicineByID)
	api.Get("/pharmacies/unapproved", handlers.GetUnapprovedPharmacies)
	api.Get("/categories", handlers.GetAllCategories)
	app.Put("/api/medicines/:id", handlers.UpdateMedicine)
	api.Post("/cart", handlers.AddToCart)                     
    api.Get("/cart/:userId", handlers.GetCart)                
    api.Put("/cart/item", handlers.UpdateCartItem)            
    api.Delete("/cart/item/:cartItemId", handlers.RemoveCartItem)
	api.Delete("/cart/clear/:userId", handlers.ClearCartHandler)
	api.Delete("/cart/item/:cartItemId", handlers.RemoveCartItemHandler)
}

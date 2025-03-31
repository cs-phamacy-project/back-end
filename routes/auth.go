package routes

import (
	"github.com/gofiber/fiber/v2"
	"gomed/handlers"
	"github.com/gofiber/websocket/v2"
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



func SetupOrderRoutes(app *fiber.App) {
	// กำหนดเส้นทาง API
	orders := app.Group("/api/orders")
	orders.Post("/", handlers.CreateOrder) // สร้างคำสั่งซื้อใหม่
	orders.Get("/:order_id", handlers.GetOrderDetails) // ดึงข้อมูลคำสั่งซื้อ
  }

  func SetupPaymentRoutes(app *fiber.App) {
	// กำหนดเส้นทาง API
	api := app.Group("/api/payment")
	
	// API สำหรับการชำระเงิน
	api.Post("/charge", handlers.CreatePayment)
	api.Get("/:id", handlers.GetPaymentByID)
}

func SetupChatRoutes(app *fiber.App) {
    // กลุ่มเส้นทางสำหรับผู้ใช้ (ไม่มี middleware)
    userRoutes := app.Group("/api/user")

    // การสนทนาของผู้ใช้
    userRoutes.Get("/conversations", handlers.GetUserConversationsHandler)
    userRoutes.Post("/conversations", handlers.CreateConversationHandler)
    userRoutes.Get("/conversations/:id/messages", handlers.GetMessagesHandler)
    userRoutes.Post("/conversations/:id/messages", handlers.SendMessageHandler)

    // กลุ่มเส้นทางสำหรับร้านขายยา (ไม่มี middleware)
    pharmacyRoutes := app.Group("/api/pharmacy")

    // การสนทนาของร้านขายยา
    pharmacyRoutes.Get("/conversations", handlers.GetPharmacyConversationsHandler)
    pharmacyRoutes.Get("/conversations/:id/messages", handlers.GetMessagesHandler)
    pharmacyRoutes.Post("/conversations/:id/messages", handlers.SendMessageHandler)
    
    // ตั้งค่า WebSocket
    app.Use("/ws", func(c *fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            c.Locals("allowed", true)
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })

    // เส้นทาง WebSocket chat
    app.Get("/ws/chat", websocket.New(handlers.WebSocketHandler))
}
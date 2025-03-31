package handlers

import(
	"gomed/models"
	"gomed/database"
	"github.com/gofiber/fiber/v2"
	"errors"
	"gorm.io/gorm"
)

func CreateOrder(c *fiber.Ctx) error {
	// รับข้อมูล request
	var req struct {
	  UserID          uint     `json:"user_id"`
	  CartItemIDs     []uint   `json:"cart_item_ids"`
	  ShippingAddress string   `json:"shipping_address"`
	  PhoneNumber     string   `json:"phone_number"`
	  PaymentMethod   string   `json:"payment_method"`
	  TotalAmount     float64  `json:"total_amount"`
	}
  
	if err := c.BodyParser(&req); err != nil {
	  return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "รูปแบบข้อมูลไม่ถูกต้อง: " + err.Error(),
	  })
	}


  
	// สร้างคำสั่งซื้อในฐานข้อมูล
	db := database.GetDB()
	order := models.Order{
	  UserID:          req.UserID,
	  PaymentMethod:   req.PaymentMethod,
	  PaymentStatus:   "pending",
	  OrderStatus:     "pending",
	  TotalAmount:     req.TotalAmount,
	  ShippingAddress: req.ShippingAddress,
	  PhoneNumber:     req.PhoneNumber,
	}
	
  
	if err := db.Create(&order).Error; err != nil {
	  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "ไม่สามารถสร้างคำสั่งซื้อได้: " + err.Error(),
	  })
	}
  
	// ย้ายสินค้าจากตะกร้าไปเป็นรายการสั่งซื้อ
	for _, cartItemID := range req.CartItemIDs {
	  var cartItem models.CartItem
	  if err := db.Preload("Medicine").First(&cartItem, cartItemID).Error; err == nil {
		orderItem := models.OrderItem{
		  OrderID:    order.OrderID,
		  MedicineID: cartItem.MedicineID,
		  Quantity:   cartItem.Quantity,
		  TotalPrice:      cartItem.Medicine.Price,
		}
		
		db.Create(&orderItem)
		
		// ลบรายการออกจากตะกร้า
		db.Delete(&cartItem)
	  }
	}
  
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	  "success": true,
	  "message": "สร้างคำสั่งซื้อสำเร็จ",
	  "order_id": order.OrderID,
	})
  }

  // GetOrderDetails retrieves detailed information for a specific order
func GetOrderDetails(c *fiber.Ctx) error {
	// รับ order_id จาก params
	orderID, err := c.ParamsInt("order_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รูปแบบ order ID ไม่ถูกต้อง",
		})
	}

	// ดึงข้อมูลจากฐานข้อมูล
	db := database.GetDB()
	var order models.Order
	
	// ดึงข้อมูล order พร้อมกับ items และ medicine ที่เกี่ยวข้อง
	if err := db.Preload("Items.Medicine").First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ไม่พบข้อมูลคำสั่งซื้อ",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการดึงข้อมูล: " + err.Error(),
		})
	}

	// คำนวณยอดรวมจำนวนสินค้า
	var totalQuantity int
	for _, item := range order.Items {
		totalQuantity += item.Quantity
	}
	order.TotalQuantity = totalQuantity

	// สร้าง response
	orderDetails := fiber.Map{
		"order_id":         order.OrderID,
		"user_id":          order.UserID,
		"order_status":     order.OrderStatus,
		"payment_status":   order.PaymentStatus,
		"payment_method":   order.PaymentMethod,
		"total_amount":     order.TotalAmount,
		"total_quantity":   order.TotalQuantity,
		"shipping_address": order.ShippingAddress,
		"phone_number":     order.PhoneNumber,
		"items":            []fiber.Map{},
	}

	// เพิ่มข้อมูลรายการสินค้า
	for _, item := range order.Items {
		orderDetails["items"] = append(orderDetails["items"].([]fiber.Map), fiber.Map{
			"order_item_id": item.OrderItemID,
			"medicine_id":   item.MedicineID,
			"medicine_name": item.Medicine.ProductName,
			"quantity":      item.Quantity,
			"price":         item.Medicine.Price,
			"total_price":   item.TotalPrice,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    orderDetails,
	})
}
package handlers


import (
	"strconv"
	"gomed/database"
	"gomed/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AddToCartRequest รับข้อมูล request สำหรับเพิ่มสินค้าลงตะกร้า
type AddToCartRequest struct {
	UserID     uint `json:"userId" binding:"required"`
	PharmacyID uint `json:"pharmacyId" binding:"required"`
	MedicineID uint `json:"medicineId" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required,min=1"`
}

// AddToCartResponse ส่งข้อมูล response กลับ
type AddToCartResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	CartID  uint        `json:"cartId,omitempty"`
	Item    *models.CartItem `json:"item,omitempty"`
}

// AddToCart เพิ่มสินค้าลงในตะกร้า
func AddToCart(c *fiber.Ctx) error {
	var request AddToCartRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(AddToCartResponse{
			Success: false,
			Message: "ข้อมูลไม่ถูกต้อง: " + err.Error(),
		})
	}

	db := database.GetDB()

	// ตรวจสอบว่ามีตะกร้าของผู้ใช้และร้านขายยานี้อยู่แล้วหรือไม่
	var cart models.Cart
	result := db.Where("user_id = ? AND pharmacy_id = ? AND status = ?", 
		request.UserID, request.PharmacyID, "active").First(&cart)

	// ถ้ายังไม่มีตะกร้า ให้สร้างตะกร้าใหม่
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			cart = models.Cart{
				UserID:     request.UserID,
				PharmacyID: request.PharmacyID,
				Status:     "active",
			}
			if err := db.Create(&cart).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(AddToCartResponse{
					Success: false,
					Message: "ไม่สามารถสร้างตะกร้าใหม่ได้: " + err.Error(),
				})
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(AddToCartResponse{
				Success: false,
				Message: "เกิดข้อผิดพลาดในการค้นหาตะกร้า: " + result.Error.Error(),
			})
		}
	}

	// ตรวจสอบว่ามีสินค้านี้ในตะกร้าแล้วหรือไม่
	var cartItem models.CartItem
	result = db.Where("cart_id = ? AND medicine_id = ?", cart.CartID, request.MedicineID).First(&cartItem)

	// ถ้ามีสินค้านี้ในตะกร้าแล้ว ให้อัปเดตจำนวน
	if result.Error == nil {
		cartItem.Quantity += request.Quantity
		if err := db.Save(&cartItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(AddToCartResponse{
				Success: false,
				Message: "ไม่สามารถอัปเดตจำนวนสินค้าในตะกร้าได้: " + err.Error(),
			})
		}
	} else if result.Error == gorm.ErrRecordNotFound {
		// ถ้ายังไม่มีสินค้านี้ในตะกร้า ให้เพิ่มสินค้าใหม่
		cartItem = models.CartItem{
			CartID:     cart.CartID,
			MedicineID: request.MedicineID,
			Quantity:   request.Quantity,
		}
		if err := db.Create(&cartItem).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(AddToCartResponse{
				Success: false,
				Message: "ไม่สามารถเพิ่มสินค้าลงตะกร้าได้: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(AddToCartResponse{
			Success: false,
			Message: "เกิดข้อผิดพลาดในการค้นหาสินค้าในตะกร้า: " + result.Error.Error(),
		})
	}

	// ดึงข้อมูลสินค้าใหม่
	if err := db.Preload("Medicine").First(&cartItem, cartItem.CartItemID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(AddToCartResponse{
			Success: false,
			Message: "เพิ่มสินค้าสำเร็จแล้ว แต่พบข้อผิดพลาดในการดึงข้อมูล",
		})
	}

	return c.Status(fiber.StatusOK).JSON(AddToCartResponse{
		Success: true,
		Message: "เพิ่มสินค้าลงตะกร้าเรียบร้อยแล้ว",
		CartID:  cart.CartID,
		Item:    &cartItem,
	})
}

// ฟังก์ชั่นเพิ่มเติมที่อาจเป็นประโยชน์

// GetCart ดึงข้อมูลตะกร้าสินค้าของผู้ใช้
func GetCart(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	db := database.GetDB()
	var cart models.Cart
	result := db.Where("user_id = ? AND status = ?", userID, "active").
		Preload("CartItems").
		Preload("CartItems.Medicine").
		First(&cart)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"message": "ยังไม่มีตะกร้าสินค้า",
				"cart":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการดึงข้อมูลตะกร้า: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "ดึงข้อมูลตะกร้าสำเร็จ",
		"cart":    cart,
	})
}

// UpdateCartItem อัปเดตจำนวนสินค้าในตะกร้า
func UpdateCartItem(c *fiber.Ctx) error {
	var request struct {
		CartItemID uint `json:"cartItemId"`
		Quantity   int  `json:"quantity" validate:"min=1"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ถูกต้อง: " + err.Error(),
		})
	}

	// ตรวจสอบข้อมูล
	if request.CartItemID == 0 || request.Quantity < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ถูกต้อง: CartItemID และ Quantity ต้องระบุ และ Quantity ต้องมากกว่า 0",
		})
	}

	db := database.GetDB()
	var cartItem models.CartItem
	if err := db.First(&cartItem, request.CartItemID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบรายการในตะกร้า",
		})
	}

	cartItem.Quantity = request.Quantity
	if err := db.Save(&cartItem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอัปเดตจำนวนสินค้าได้: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "อัปเดตจำนวนสินค้าเรียบร้อย",
		"item":    cartItem,
	})
}

// RemoveCartItem ลบสินค้าออกจากตะกร้า
func RemoveCartItem(c *fiber.Ctx) error {
	cartItemID, err := strconv.ParseUint(c.Params("cartItemId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสรายการไม่ถูกต้อง",
		})
	}

	db := database.GetDB()
	result := db.Delete(&models.CartItem{}, cartItemID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถลบสินค้าออกจากตะกร้าได้: " + result.Error.Error(),
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบรายการในตะกร้า",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "ลบสินค้าออกจากตะกร้าเรียบร้อย",
	})
}

// ClearCartHandler สำหรับล้างตะกร้าทั้งหมดและรีเซ็ต pharmacy_id
func ClearCartHandler(c *fiber.Ctx) error {
    db := database.GetDB()
    
    userID, err := strconv.Atoi(c.Params("userId"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid user ID",
        })
    }
    
    // ค้นหาตะกร้าของผู้ใช้
    var cart models.Cart
    result := db.Where("user_id = ? AND status = ?", userID, "active").First(&cart)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "Cart not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Error finding cart",
        })
    }
    
    // เริ่ม transaction
    tx := db.Begin()
    if tx.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to start transaction",
        })
    }
    
    // ลบทุกรายการสินค้าในตะกร้า
    if err := tx.Where("cart_id = ?", cart.CartID).Delete(&models.CartItem{}).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Error clearing cart items",
        })
    }
    
    // รีเซ็ต pharmacy_id
    if err := tx.Model(&models.Cart{}).Where("cart_id = ?", cart.CartID).Update("pharmacy_id", nil).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Error resetting pharmacy ID",
        })
    }
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to commit transaction",
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Cart cleared successfully",
    })
}

// RemoveCartItemHandler สำหรับลบสินค้าจากตะกร้า
func RemoveCartItemHandler(c *fiber.Ctx) error {
    db := database.GetDB()
    
    cartItemID, err := strconv.Atoi(c.Params("cartItemId"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid cart item ID",
        })
    }
    
    // ค้นหารายการสินค้าในตะกร้า
    var cartItem models.CartItem
    result := db.First(&cartItem, cartItemID)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "Cart item not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Error finding cart item",
        })
    }
    
    // เริ่ม transaction
    tx := db.Begin()
    if tx.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to start transaction",
        })
    }
    
    // ลบรายการสินค้า
    if err := tx.Delete(&cartItem).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Error removing cart item",
        })
    }
    
    // ตรวจสอบว่ายังมีสินค้าเหลือในตะกร้าหรือไม่
    var count int64
    if err := tx.Model(&models.CartItem{}).Where("cart_id = ?", cartItem.CartID).Count(&count).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Error checking remaining items",
        })
    }
    
    // ถ้าไม่มีสินค้าเหลือ ให้รีเซ็ต pharmacy_id
    if count == 0 {
        if err := tx.Model(&models.Cart{}).Where("cart_id = ?", cartItem.CartID).Update("pharmacy_id", nil).Error; err != nil {
            tx.Rollback()
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Error resetting pharmacy ID",
            })
        }
    }
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to commit transaction",
        })
    }
    
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Item removed from cart",
    })
}
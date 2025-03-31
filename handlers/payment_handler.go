package handlers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

// OmiseChargeRequest รับข้อมูลจาก client
type OmiseChargeRequest struct {
	Token       string                 `json:"token"`
	Amount      int64                  `json:"amount"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CreatePayment จัดการการชำระเงินผ่าน Omise
func CreatePayment(c *fiber.Ctx) error {
	// สร้าง Omise client
	client, err := omise.NewClient(os.Getenv("OMISE_PUBLIC_KEY"), os.Getenv("OMISE_SECRET_KEY"))
	if err != nil {
		log.Printf("Error creating Omise client: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถเชื่อมต่อกับ Omise ได้",
		})
	}

	// รับข้อมูลจาก request
	var req OmiseChargeRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รูปแบบข้อมูลไม่ถูกต้อง",
		})
	}

	// ตรวจสอบข้อมูลที่จำเป็น
	if req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบ token",
		})
	}

	if req.Amount < 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "จำนวนเงินต้องมากกว่า 100 สตางค์ (1 บาท)",
		})
	}

	// ตั้งค่า default description ถ้าไม่มี
	if req.Description == "" {
		req.Description = "ชำระสินค้า"
	}

	// บันทึก log สำหรับการดีบัก
	log.Printf("Creating charge with token: %s", req.Token)
	log.Printf("Amount: %d", req.Amount)

	// สร้าง charge
	charge := &omise.Charge{}
	createCharge := &operations.CreateCharge{
		Amount:      req.Amount,
		Currency:    "thb",
		Card:        req.Token,
		Description: req.Description,
	}

	// ถ้ามี metadata ให้เพิ่มลงไป
	if req.Metadata != nil {
		createCharge.Metadata = req.Metadata
	} else {
		// สร้าง default metadata ด้วยเวลาปัจจุบัน
		createCharge.Metadata = map[string]interface{}{
			"order_id": fmt.Sprintf("order_%d", time.Now().Unix()),
		}
	}

	// ดำเนินการ create charge
	createErr := client.Do(charge, createCharge)
	if createErr != nil {
		log.Printf("Error creating charge: %v", createErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการชำระเงิน",
			"error":   createErr.Error(),
		})
	}

	log.Printf("Charge created: %s", charge.ID)
	log.Printf("Charge status: %s", charge.Status)

	// ตรวจสอบสถานะการชำระเงิน
	if charge.Status == "successful" || charge.Status == "pending" {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "การชำระเงินสำเร็จ",
			"charge": fiber.Map{
				"id":             charge.ID,
				"amount":         float64(charge.Amount) / 100, // แปลงกลับเป็นบาท
				"status":         charge.Status,
				"paid":           charge.Paid,
				"transaction_id": charge.Transaction,
			},
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("การชำระเงินล้มเหลว: %v", charge.FailureMessage),
			"charge": fiber.Map{
				"id":              charge.ID,
				"status":          charge.Status,
				"failure_code":    charge.FailureCode,
				"failure_message": charge.FailureMessage,
			},
		})
	}
}

// GetPaymentByID ดึงข้อมูลการชำระเงิน
func GetPaymentByID(c *fiber.Ctx) error {
	// สร้าง Omise client
	client, err := omise.NewClient(os.Getenv("OMISE_PUBLIC_KEY"), os.Getenv("OMISE_SECRET_KEY"))
	if err != nil {
		log.Printf("Error creating Omise client: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถเชื่อมต่อกับ Omise ได้",
		})
	}
	
	// ดึง ID จาก URL params
	chargeID := c.Params("id")
	if chargeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบ ID",
		})
	}
	
	// ดึงข้อมูล charge
	charge := &omise.Charge{}
	retrieveCharge := &operations.RetrieveCharge{
		ChargeID: chargeID,
	}
	
	if err := client.Do(charge, retrieveCharge); err != nil {
		log.Printf("Error retrieving charge: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบข้อมูลการชำระเงิน",
			"error":   err.Error(),
		})
	}
	
	return c.JSON(fiber.Map{
		"success": true,
		"charge": fiber.Map{
			"id":              charge.ID,
			"amount":          float64(charge.Amount) / 100,
			"status":          charge.Status,
			"paid":            charge.Paid,
			"transaction_id":  charge.Transaction,
			"failure_code":    charge.FailureCode,
			"failure_message": charge.FailureMessage,
		},
	})
}
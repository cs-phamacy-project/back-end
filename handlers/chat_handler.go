// handlers/chat_handler.go
package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"gomed/database"
	"gomed/models"
)

// GetUserConversationsHandler ดึงรายการแชทของผู้ใช้
func GetUserConversationsHandler(c *fiber.Ctx) error {
	// ดึง user_id จาก query parameter แทนที่จะดึงจาก c.Locals
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing user_id parameter",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user_id parameter",
		})
	}

	var conversations []models.Conversation
	result := database.DB.
		Preload("Pharmacy").
		Where("user_id = ?", userID).
		Find(&conversations)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching conversations",
		})
	}

	type ConversationResponse struct {
		ID           uint      `json:"id"`
		PharmacyID   uint      `json:"pharmacy_id"`
		StoreName    string    `json:"store_name"`
		LastMessage  string    `json:"last_message"`
		LastSent     time.Time `json:"last_sent"`
		CreatedAt    time.Time `json:"created_at"`
	}

	var response []ConversationResponse

	for _, conversation := range conversations {
		// ดึงข้อความล่าสุด
		var lastMessage models.Message
		database.DB.Where("conversation_id = ?", conversation.ID).
			Order("created_at DESC").
			Limit(1).
			Find(&lastMessage)

		response = append(response, ConversationResponse{
			ID:          conversation.ID,
			PharmacyID:  conversation.PharmacyID,
			StoreName:   conversation.Pharmacy.StoreName,
			LastMessage: lastMessage.Content,
			LastSent:    lastMessage.CreatedAt,
			CreatedAt:   conversation.CreatedAt,
		})
	}

	return c.JSON(response)
}

// GetPharmacyConversationsHandler ดึงรายการแชทของร้านขายยา
func GetPharmacyConversationsHandler(c *fiber.Ctx) error {
	// ดึง pharmacy_id จาก query parameter แทนที่จะดึงจาก c.Locals
	pharmacyIDStr := c.Query("pharmacy_id")
	if pharmacyIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing pharmacy_id parameter",
		})
	}

	pharmacyID, err := strconv.Atoi(pharmacyIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pharmacy_id parameter",
		})
	}

	var conversations []models.Conversation
	result := database.DB.
		Preload("User").
		Where("pharmacy_id = ?", pharmacyID).
		Find(&conversations)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching conversations",
		})
	}

	type ConversationResponse struct {
		ID           uint      `json:"id"`
		UserID       uint      `json:"user_id"`
		FirstName    string    `json:"first_name"`
		LastName     string    `json:"last_name"`
		LastMessage  string    `json:"last_message"`
		LastSent     time.Time `json:"last_sent"`
		CreatedAt    time.Time `json:"created_at"`
	}

	var response []ConversationResponse

	for _, conversation := range conversations {
		// ดึงข้อความล่าสุด
		var lastMessage models.Message
		database.DB.Where("conversation_id = ?", conversation.ID).
			Order("created_at DESC").
			Limit(1).
			Find(&lastMessage)

		response = append(response, ConversationResponse{
			ID:          conversation.ID,
			UserID:      conversation.UserID,
			FirstName:   conversation.User.FirstName,
			LastName:    conversation.User.LastName,
			LastMessage: lastMessage.Content,
			LastSent:    lastMessage.CreatedAt,
			CreatedAt:   conversation.CreatedAt,
		})
	}

	return c.JSON(response)
}

// GetMessagesHandler ดึงข้อความในการสนทนา
func GetMessagesHandler(c *fiber.Ctx) error {
	conversationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid conversation ID",
		})
	}

	// ตรวจสอบสิทธิ์ด้วย query parameter
	userIDStr := c.Query("user_id")
	pharmacyIDStr := c.Query("pharmacy_id")

	if userIDStr != "" {
		// ตรวจสอบว่าเป็นผู้ใช้จริง
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user_id parameter",
			})
		}

		var conversation models.Conversation
		result := database.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation)
		if result.Error != nil || result.RowsAffected == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to access this conversation",
			})
		}
	} else if pharmacyIDStr != "" {
		// ตรวจสอบว่าเป็นร้านขายยาจริง
		pharmacyID, err := strconv.Atoi(pharmacyIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid pharmacy_id parameter",
			})
		}

		var conversation models.Conversation
		result := database.DB.Where("id = ? AND pharmacy_id = ?", conversationID, pharmacyID).First(&conversation)
		if result.Error != nil || result.RowsAffected == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to access this conversation",
			})
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing user_id or pharmacy_id parameter",
		})
	}

	// ดึงข้อความในการสนทนา
	var messages []models.Message
	result := database.DB.Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching messages",
		})
	}

	return c.JSON(messages)
}

// CreateConversationHandler สร้างการสนทนาใหม่ (สำหรับผู้ใช้)
func CreateConversationHandler(c *fiber.Ctx) error {
	// ดึง user_id จาก query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing user_id parameter",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user_id parameter",
		})
	}

	type CreateConversationRequest struct {
		PharmacyID uint   `json:"pharmacy_id"`
		Message    string `json:"message"`
	}

	var request CreateConversationRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// ตรวจสอบว่ามีร้านขายยาที่ระบุหรือไม่
	var pharmacy models.Pharmacy
	if result := database.DB.First(&pharmacy, request.PharmacyID); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pharmacy not found",
		})
	}

	// ตรวจสอบว่ามีการสนทนาระหว่างผู้ใช้และร้านขายยานี้อยู่แล้วหรือไม่
	var existingConversation models.Conversation
	result := database.DB.Where("user_id = ? AND pharmacy_id = ?", userID, request.PharmacyID).First(&existingConversation)

	var conversationID uint
	if result.Error == gorm.ErrRecordNotFound {
		// สร้างการสนทนาใหม่
		conversation := models.Conversation{
			UserID:     uint(userID),
			PharmacyID: request.PharmacyID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if result := database.DB.Create(&conversation); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create conversation",
			})
		}

		conversationID = conversation.ID
	} else if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check existing conversations",
		})
	} else {
		// ใช้การสนทนาที่มีอยู่แล้ว
		conversationID = existingConversation.ID
		// อัปเดตเวลาล่าสุด
		database.DB.Model(&existingConversation).Update("updated_at", time.Now())
	}

	// สร้างข้อความ
	message := models.Message{
		ConversationID: conversationID,
		SenderID:       uint(userID),
		SenderType:     "user",
		Content:        request.Message,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if result := database.DB.Create(&message); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create message",
		})
	}

	return c.JSON(fiber.Map{
		"conversation_id": conversationID,
		"message":         message,
	})
}

// SendMessageHandler ส่งข้อความ
func SendMessageHandler(c *fiber.Ctx) error {
	conversationID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid conversation ID",
		})
	}

	// ตรวจสอบสิทธิ์ด้วย query parameter
	userIDStr := c.Query("user_id")
	pharmacyIDStr := c.Query("pharmacy_id")

	var senderID uint
	var senderType string

	if userIDStr != "" {
		// ตรวจสอบว่าเป็นผู้ใช้จริง
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user_id parameter",
			})
		}

		var conversation models.Conversation
		result := database.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation)
		if result.Error != nil || result.RowsAffected == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to access this conversation",
			})
		}
		senderID = uint(userID)
		senderType = "user"
	} else if pharmacyIDStr != "" {
		// ตรวจสอบว่าเป็นร้านขายยาจริง
		pharmacyID, err := strconv.Atoi(pharmacyIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid pharmacy_id parameter",
			})
		}

		var conversation models.Conversation
		result := database.DB.Where("id = ? AND pharmacy_id = ?", conversationID, pharmacyID).First(&conversation)
		if result.Error != nil || result.RowsAffected == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to access this conversation",
			})
		}
		senderID = uint(pharmacyID)
		senderType = "pharmacy"
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing user_id or pharmacy_id parameter",
		})
	}

	type SendMessageRequest struct {
		Content string `json:"content"`
	}

	var request SendMessageRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message content cannot be empty",
		})
	}

	// สร้างข้อความ
	message := models.Message{
		ConversationID: uint(conversationID),
		SenderID:       senderID,
		SenderType:     senderType,
		Content:        request.Content,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if result := database.DB.Create(&message); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create message",
		})
	}

	// อัปเดตเวลาล่าสุดของการสนทนา
	database.DB.Model(&models.Conversation{}).Where("id = ?", conversationID).Update("updated_at", time.Now())

	return c.JSON(message)
}
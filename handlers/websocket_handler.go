// handlers/websocket_handler.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"

	"gomed/database"
	"gomed/models"
)

// Client เป็นโครงสร้างข้อมูลของ WebSocket client
type Client struct {
	ID        string
	Conn      *websocket.Conn
	UserID    uint
	PharmacyID uint
	UserType  string // "user" หรือ "pharmacy"
	RoomID    uint   // conversation ID
	mu        sync.Mutex
}

// ข้อความที่ส่งผ่าน WebSocket
type ChatMessage struct {
	ConversationID uint   `json:"conversation_id"`
	SenderID       uint   `json:"sender_id"`
	SenderType     string `json:"sender_type"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at"`
	MessageID      uint   `json:"message_id"`
}

// บันทึกการเชื่อมต่อทั้งหมด
var (
	clients    = make(map[string]*Client)
	register   = make(chan *Client)
	unregister = make(chan *Client)
	broadcast  = make(chan ChatMessage)
	rooms      = make(map[uint]map[string]bool) // conversation ID -> map[clientID]bool
	mutex      = &sync.Mutex{}
)

// เริ่มต้น WebSocket server
func RunWebSocketServer() {
	for {
		select {
		case client := <-register:
			mutex.Lock()
			// เพิ่ม client ในห้องสนทนา
			if _, ok := rooms[client.RoomID]; !ok {
				rooms[client.RoomID] = make(map[string]bool)
			}
			rooms[client.RoomID][client.ID] = true
			clients[client.ID] = client
			mutex.Unlock()

		case client := <-unregister:
			mutex.Lock()
			// ลบ client จากห้องสนทนา
			if _, ok := rooms[client.RoomID]; ok {
				delete(rooms[client.RoomID], client.ID)
				if len(rooms[client.RoomID]) == 0 {
					delete(rooms, client.RoomID)
				}
			}
			delete(clients, client.ID)
			client.Conn.Close()
			mutex.Unlock()

		case message := <-broadcast:
			// แปลงข้อความเป็น JSON
			data, err := json.Marshal(message)
			if err != nil {
				log.Println("Error marshalling message:", err)
				continue
			}

			// ส่งข้อความไปยังทุก client ในห้องสนทนา
			mutex.Lock()
			if roomClients, ok := rooms[message.ConversationID]; ok {
				for clientID := range roomClients {
					client, exists := clients[clientID]
					if !exists {
						continue
					}
					
					client.mu.Lock()
					err := client.Conn.WriteMessage(websocket.TextMessage, data)
					client.mu.Unlock()
					
					if err != nil {
						log.Println("Error writing message to client:", err)
						client.Conn.Close()
						unregister <- client
					}
				}
			}
			mutex.Unlock()
		}
	}
}

// WebSocketHandler จัดการการเชื่อมต่อ WebSocket
func WebSocketHandler(c *websocket.Conn) {
	// ตรวจสอบ query parameters
	userID := c.Query("user_id")
	pharmacyID := c.Query("pharmacy_id")
	conversationID := c.Query("conversation_id")
	
	// ไม่ใช้ token ในการตรวจสอบเนื่องจากไม่มี middleware
	// token := c.Query("token")

	var clientUserID uint = 0
	var clientPharmacyID uint = 0
	var clientUserType string
	var clientRoomID uint = 0

	// ตรวจสอบว่าเป็น user หรือ pharmacy
	if userID != "" {
		clientUserID = parseUint(userID)
		clientUserType = "user"
		
		// ตรวจสอบว่า user มีอยู่จริง
		var user models.User
		if result := database.DB.First(&user, clientUserID); result.Error != nil {
			log.Println("User not found:", clientUserID)
			c.Close()
			return
		}
	} else if pharmacyID != "" {
		clientPharmacyID = parseUint(pharmacyID)
		clientUserType = "pharmacy"
		
		// ตรวจสอบว่า pharmacy มีอยู่จริง
		var pharmacy models.Pharmacy
		if result := database.DB.First(&pharmacy, clientPharmacyID); result.Error != nil {
			log.Println("Pharmacy not found:", clientPharmacyID)
			c.Close()
			return
		}
	} else {
		log.Println("Missing user_id or pharmacy_id parameter")
		c.Close()
		return
	}

	// ตรวจสอบสิทธิ์ในการเข้าถึงห้องสนทนา
	clientRoomID = parseUint(conversationID)
	if clientRoomID == 0 {
		log.Println("Invalid conversation_id parameter")
		c.Close()
		return
	}
	
	// ตรวจสอบว่าการสนทนามีอยู่จริงและผู้ใช้มีสิทธิ์เข้าถึง
	var conversation models.Conversation
	var result *gorm.DB
	
	if clientUserType == "user" {
		result = database.DB.Where("id = ? AND user_id = ?", clientRoomID, clientUserID).First(&conversation)
	} else {
		result = database.DB.Where("id = ? AND pharmacy_id = ?", clientRoomID, clientPharmacyID).First(&conversation)
	}
	
	if result.Error != nil || result.RowsAffected == 0 {
		log.Println("Conversation not found or access denied")
		c.Close()
		return
	}

	// สร้าง client ใหม่
	client := &Client{
		ID:        fmt.Sprintf("%s-%d-%d", clientUserType, time.Now().UnixNano(), clientRoomID),
		Conn:      c,
		UserID:    clientUserID,
		PharmacyID: clientPharmacyID,
		UserType:  clientUserType,
		RoomID:    clientRoomID,
	}

	// ลงทะเบียน client
	register <- client

	// อ่านข้อความจาก WebSocket
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			unregister <- client
			break
		}

		// แปลงข้อความเป็นโครงสร้างข้อมูล
		var chatMsg ChatMessage
		if err := json.Unmarshal(msg, &chatMsg); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// ตรวจสอบว่า client มีสิทธิ์ส่งข้อความในนามของผู้ส่งหรือไม่
		if (chatMsg.SenderType == "user" && clientUserID != chatMsg.SenderID) ||
			(chatMsg.SenderType == "pharmacy" && clientPharmacyID != chatMsg.SenderID) {
			log.Println("Unauthorized message sender")
			continue
		}

		// ตรวจสอบว่า client มีสิทธิ์ส่งข้อความในห้องสนทนานี้หรือไม่
		if chatMsg.ConversationID != clientRoomID {
			log.Println("Unauthorized conversation access")
			continue
		}

		// บันทึกข้อความลงฐานข้อมูล
		message := models.Message{
			ConversationID: chatMsg.ConversationID,
			SenderID:       chatMsg.SenderID,
			SenderType:     chatMsg.SenderType,
			Content:        chatMsg.Content,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		result := database.DB.Create(&message)
		if result.Error != nil {
			log.Println("Error saving message:", result.Error)
			continue
		}

		// อัปเดตเวลาล่าสุดของการสนทนา
		database.DB.Model(&models.Conversation{}).
			Where("id = ?", chatMsg.ConversationID).
			Update("updated_at", time.Now())

		// ส่งข้อความที่บันทึกแล้วกลับไปยัง client พร้อม ID
		chatMsg.MessageID = message.ID
		chatMsg.CreatedAt = message.CreatedAt.Format(time.RFC3339)

		// ส่งข้อความต่อไปยัง client อื่น
		broadcast <- chatMsg
	}
}

// SetupWebSocketServer เริ่มต้น WebSocket server
func SetupWebSocketServer() {
	go RunWebSocketServer()
}

// Helper function to parse uint
func parseUint(s string) uint {
	var result uint
	fmt.Sscanf(s, "%d", &result)
	return result
}
package models

type Cart struct {
    CartID     uint   `gorm:"primaryKey"`
    UserID     uint   `gorm:"not null"`
    Status     string `gorm:"default:'active'"`
    Orders     []Order `gorm:"foreignKey:CartID"`
}

type Order struct {
    OrderID       uint    `gorm:"primaryKey"`
    CartID        uint    `gorm:"not null"`
    UserID        uint    `gorm:"not null"`
    TotalQuantity int     `gorm:"not null"`
    TotalAmount   float64 `gorm:"not null"`
    OrderStatus   string  `gorm:"default:'pending'"`
    Items         []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
    OrderItemID uint    `gorm:"primaryKey"`
    OrderID     uint    `gorm:"not null"`
    MedicineID  uint      `gorm:"not null" json:"medicine_id"`
    Medicine    Medicine  `gorm:"foreignKey:MedicineID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"medicine"`
    Quantity    int     `gorm:"not null;check:quantity > 0"`
    Price       float64 `gorm:"not null"`
}
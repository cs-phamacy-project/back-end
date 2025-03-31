package models

type Cart struct {
    CartID     uint `gorm:"primaryKey;column:cart_id"`
    UserID     uint `gorm:"column:user_id"`
    PharmacyID uint `gorm:"column:pharmacy_id"` // ใช้ uint แทน *uint
    Status     string
    CartItems  []CartItem `gorm:"foreignKey:CartID"`
}

type Order struct {
    OrderID       uint    `gorm:"primaryKey"`
    CartID              uint      `gorm:"not null;uniqueIndex"`
    UserID        uint    `gorm:"not null"`
    TotalQuantity int     `gorm:"not null"`
    TotalAmount   float64 `gorm:"not null"`
    OrderStatus   string  `gorm:"default:'pending'"`
    PaymentStatus       string    `gorm:"default:'pending'"`
    Items         []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
    OrderItemID uint    `gorm:"primaryKey"`
    OrderID     uint    `gorm:"not null"`
    MedicineID  uint      `gorm:"not null" json:"medicine_id"`
    Medicine    Medicine  `gorm:"foreignKey:MedicineID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"medicine"`
    Quantity    int     `gorm:"not null;check:quantity > 0"`
    TotalPrice       float64 `gorm:"not null"`
}

type CartItem struct {
    CartItemID  uint      `gorm:"primaryKey"`
    CartID      uint      `gorm:"not null"`
    MedicineID  uint      `gorm:"not null"`
    Medicine    Medicine  `gorm:"foreignKey:MedicineID"`
    Quantity    int       `gorm:"not null;check:quantity > 0"`
}
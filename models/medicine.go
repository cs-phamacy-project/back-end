package models

type Medicine struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    ProductName string    `gorm:"not null" json:"product_name"`
    Description string    `gorm:"not null" json:"description"`
    CategoryID  uint      `gorm:"not null" json:"category_id"`
    Category    Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
    PharmacyID  uint      `gorm:"not null" json:"pharmacy_id"`
    Pharmacy    Pharmacy `gorm:"foreignKey:PharmacyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"pharmacy"`
    Price       float64   `gorm:"not null" json:"price"`
    Stock       int       `gorm:"not null" json:"stock"`
    ExpiredDate string    `gorm:"not null" json:"expired_date"`
    Fda         string    `gorm:"not null" json:"fda"`
    Status      string    `gorm:"not null" json:"status"`
    Image       string    `json:"image"`
}

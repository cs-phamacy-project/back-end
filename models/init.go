package models

import "gorm.io/gorm"

func MigrateDB(db *gorm.DB) {
    db.AutoMigrate(&User{}, &Pharmacy{}, &Medicine{}, &Cart{}, &Order{}, &OrderItem{}, &Category{}, &CartItem{}, &Payment{})
}

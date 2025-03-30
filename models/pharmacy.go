package models

type Pharmacy struct {
	ID                 uint       `gorm:"primaryKey"`
	Email              string     `gorm:"unique;not null"`
	Password           string     `gorm:"not null"`
	FirstName          string     `gorm:"not null"`
	LastName           string     `gorm:"not null"`
	Phone              string     `gorm:"not null"`
	Certificate        string     `gorm:"not null"`
	StoreImg           string     `gorm:"not null"`
	AddressDescription string     `gorm:"not null"`
	SubDistrict        string     `gorm:"not null"`
	District           string     `gorm:"not null"`
	Province           string     `gorm:"not null"`
	ZipCode            string     `gorm:"not null"`
	Status             string     `gorm:"default:'unapprove'"`
	Role               string     `gorm:"default:'staff'"`
	Medicines          []Medicine `gorm:"foreignKey:PharmacyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

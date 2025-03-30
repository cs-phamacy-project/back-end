package models

type User struct {
	ID                 uint   `gorm:"primaryKey"`
	Email              string `gorm:"unique;not null"`
	Password           string `gorm:"not null"`
	FirstName          string `gorm:"not null"`
	LastName           string `gorm:"not null"`
	Phone              string `gorm:"not null"`
	AddressDescription string 
	SubDistrict        string 
	District           string
	Province           string
	ZipCode            string
	ChronicDisease     string
	MedicationAllergy  string 
	Role               string `gorm:"default:'user'"`
}



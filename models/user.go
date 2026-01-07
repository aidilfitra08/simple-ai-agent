package models

// User model represents a user in the database
type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Interests string
}

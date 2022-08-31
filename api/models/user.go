package models

import (
	"time"
)

// User model
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}

// TableName --> Table for User Model
func (User) TableName() string {
	return "users"
}

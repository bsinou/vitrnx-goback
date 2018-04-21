package model

import "github.com/jinzhu/gorm"

// User represents a user in both gorm and json worlds
type User struct {
	gorm.Model
	Name  string `gorm:"not null" json:"name"`
	Email string `gorm:"not null" json:"email"`
	ExtID string `gorm:"" json:"extId"`
	Roles []Role `gorm:"" json:"roles"`
}

type Role struct {
	gorm.Model
	Label string `gorm:"" json:"label"`
	ExtID uint   `gorm:"" json:"extId"`
}

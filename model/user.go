package model

import "github.com/jinzhu/gorm"

// User represents a user in both gorm and json worlds
type User struct {
	gorm.Model
	UserID  string `gorm:"" json:"userId"`
	Name    string `gorm:"not null" json:"name"`
	Email   string `gorm:"not null" json:"email"`
	Address string `gorm:"" json:"address"`
	Roles   []Role `gorm:"" json:"roles"`
}

type Role struct {
	gorm.Model
	Label string `gorm:"" json:"label"`
	ExtID uint   `gorm:"" json:"extId"`
}

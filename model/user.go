package model

import "github.com/jinzhu/gorm"

// User represents a user in both gorm and json worlds
type User struct {
	gorm.Model
	UserID  string `gorm:"" json:"userId"`
	Name    string `gorm:"not null" json:"name"`
	Email   string `gorm:"not null" json:"email"`
	Address string `gorm:"" json:"address"`
	// RolesStr string `gorm:"" json:"rolesStr"`
	// TODO manage this cleanly
	Roles []Role `gorm:"many2many:user_roles" json:"roles"`
}

// Role is used to manage permissions
type Role struct {
	RoleID string `gorm:"primary_key" json:"roleId"`
	Label  string `gorm:"" json:"label"`
}

package model

import "github.com/jinzhu/gorm"

// User represents a user in both gorm and json worlds
type User struct {
	gorm.Model
	UserID  string `gorm:"not null" json:"userId"`
	Name    string `gorm:"" json:"name"`
	Email   string `gorm:"not null" json:"email"`
	Address string `gorm:"" json:"address"`
	Roles   []Role `gorm:"many2many:user_roles" json:"roles"`

	// Ease front end and JSON communication, not persisted
	UserRoles []string               `gorm:"-" json:"userRoles"`
	Meta      map[string]interface{} `gorm:"-" json:"meta"`
}

// Role is used to manage permissions
type Role struct {
	RoleID string `gorm:"primary_key" json:"roleId"`
	Label  string `gorm:"" json:"label"`
}

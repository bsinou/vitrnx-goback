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

// AuthRequest is used by the front end during the auth process.
type AuthRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

// Role is used to manage permissions
type Role struct {
	RoleID string `gorm:"primary_key" json:"roleId"`
	Label  string `gorm:"" json:"label"`
}

// Group might be a user or a role
type Group struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

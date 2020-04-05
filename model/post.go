package model

import "time"

// Post model for the database
type Post struct {
	// Usually included by extending the default model
	// gorm.Model
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`

	// Model specific fields
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	AuthorID    string     `json:"authorId"`
	Tags        string     `json:"tags"`
	Desc        string     `json:"desc"`
	Hero        string     `json:"hero"`
	Thumb       string     `json:"thumb"`
	Body        string     `json:"body"`
	Audience    string     `json:"audience,omitempty"`
	Weight      int        `json:"weight,omitempty"`
	UpdatedBy   string     `json:"updatedBy,omitempty" `

	// Convenience fields to communicate with the front end. Not persisted.
	Author       string `gorm:"-" json:"author"`
	CommentCount int    `gorm:"-" json:"commentCount"`
}

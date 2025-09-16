package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	ID         uint      `json:"id" gorm:"primaryKey"`
	Login      string    `json:"login" gorm:"uniqueIndex;size:255"`
	Password   string    `json:"-"`
	Posts      []Post    `json:"posts,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}

type Post struct {
	gorm.Model `json:"-"`
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title"`
	UserID     uint      `json:"user_id"`
	Comments   []Comment `json:"comments,omitempty"`
}

type Comment struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primaryKey"`
	Title      string `json:"title"`
	PostID     uint   `json:"post_id"`
	UserID     uint   `json:"user_id"`
}

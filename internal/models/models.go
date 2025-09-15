package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login    string `gorm:"uniqueIndex;size:255"`
	Password string
	Posts    []Post
	Comments []Comment
}

type Post struct {
	gorm.Model
	Title    string
	UserID   uint
	Comments []Comment
}

type Comment struct {
	gorm.Model
	Title  string
	PostID uint
	UserID uint
}

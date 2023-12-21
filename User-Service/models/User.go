package models

import (
	"time"

	"gorm.io/gorm"
)

type UserType string

const (
	AdminUser    UserType = "admin"
	TheatreAdmin UserType = "theatre-admin"
	NormalUser   UserType = "user"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	UserID    int64          `gorm:"type:BIGINT(7);column:user_id;primarykey" json:"userId"`
	UserName  string         `gorm:"column:user_name;uniqueKey" json:"userName"`
	Email     string         `gorm:"column:email" json:"email"`
	Password  string         `gorm:"column:password" json:"password"`
	Age       int            `gorm:"column:age" json:"age"`
	UserType  UserType       `gorm:"column:user_type;type:ENUM('admin','theatre-admin','user');default:user"`
	IsActive  bool           `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`
}

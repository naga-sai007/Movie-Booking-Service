package models

import "time"

type PasswordResetRequest struct {
	Email string `json:"email"`
}

type PasswordResetForm struct {
	UserId      int64  `json:"userId"`
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type PasswordReset struct {
	ResetID   int64     `gorm:"type:BIGINT(7);column:reset_id;primarykey" json:"resetId"`
	Email     string    `gorm:"column:email" json:"email"`
	UserID    int64     `gorm:"type:BIGINT(7);column:user_id;" json:"userId"`
	Token     string    `gorm:"column:token" json:"token"`
	Expiry    time.Time `gorm:"column:expiry_time" json:"expiry"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

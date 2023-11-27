package models

import "time"

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest - структура для запроса сброса пароля
type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=3"`
	Token       string `json:"token" binding:"required"`
}

type PasswordResetRequest struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

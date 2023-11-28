package controllers

import (
	"fmt"
	"net/http"
	"net/smtp"
	"time"
	"wtd/initializers"
	"wtd/models"
	"wtd/pkg/vars"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ForgotPassword(c *gin.Context) {
	var request models.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	if err := initializers.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		fmt.Println("Error generating token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	resetRequest := models.PasswordResetRequest{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 1),
	}

	tx := initializers.DB.Begin()
	if err := tx.Create(&resetRequest).Error; err != nil {
		fmt.Println("Error creating password reset request:", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create password reset request"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Println("Error committing transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create password reset request"})
		return
	}

	fmt.Println("Password reset request created:", resetRequest)

	// Генерация URL для сброса пароля
	resetURL := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", resetRequest.Token)

	// Отправка письма
	err = sendPasswordResetEmail(user.Email, resetURL)
	if err != nil {
		fmt.Println("Error sending email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send password reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent successfully"})
}

// Send Reset Password Link To User
func sendPasswordResetEmail(to, resetURL string) error {
	auth := smtp.PlainAuth(
		"",
		vars.POST_NAME,
		vars.POST_PASS,
		vars.POST_SERVER,
	)

	msg := fmt.Sprintf("Subject: Восстановление пароля\r\nContent-Type: text/html; charset=UTF-8\r\n\r\nПерейдите по <a href='%s'>ссылке</a> чтобы установить новый пароль.", resetURL)

	// fmt.Println("SMTP Server:", vars.POST_SERVER)
	// fmt.Println("SMTP Port:", vars.POST_PORT)
	// fmt.Println("SMTP Name:", vars.POST_NAME)
	// fmt.Println("SMTP Password:", vars.POST_PASS)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", vars.POST_SERVER, vars.POST_PORT),
		auth,
		vars.POST_NAME,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("Failed to send reset link: %s", err)
	}

	return nil
}

func ResetPassword(c *gin.Context) {
	var request models.ResetPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resetRequest := models.PasswordResetRequest{}
	fmt.Println("Token in request:", request.Token)
	if err := initializers.DB.Where("token = ? AND expires_at > ?", request.Token, time.Now()).First(&resetRequest).Error; err != nil {
		fmt.Println("Error finding reset request:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired reset token"})
		return
	}
	fmt.Println("Reset request found:", resetRequest)

	user := models.User{}
	if err := initializers.DB.Where("id = ?", resetRequest.UserID).First(&user).Error; err != nil {
		fmt.Println("Error finding user:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating hashed password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate hashed password"})
		return
	}

	user.Password = string(hashedPassword)

	if err := initializers.DB.Save(&user).Error; err != nil {
		fmt.Println("Error updating user password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user password"})
		return
	}

	// TODO: Optionally, invalidate the used token

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

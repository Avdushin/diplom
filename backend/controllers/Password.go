package controllers

import (
	"fmt"
	"net/http"
	"time"
	"wtd/initializers"
	"wtd/models"
	"wtd/pkg/vars"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mail.v2"
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

//@ Prod SMTP server
// func sendPasswordResetEmail(to, resetURL string) error {
// 	m := mail.NewMessage()
// 	m.SetHeader("From", vars.POST_NAME)
// 	m.SetHeader("To", to)
// 	m.SetHeader("Subject", "Password Reset")
// 	m.SetBody("text/html", fmt.Sprintf("Click <a href='%s'>here</a> to reset your password.", resetURL))

// 	fmt.Println("Port: ", vars.POST_PORT)
// 	if vars.POST_PORT == "" {
// 		return errors.New("POST_PORT is empty")
// 	}

// 	port, err := strconv.Atoi(vars.POST_PORT)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("Port: ", port)

// 	d := mail.NewDialer(vars.POST_SERVER, port, vars.POST_NAME, vars.POST_PASS)

// 	if err := d.DialAndSend(m); err != nil {
// 		return err
// 	}

// 	return nil
// }

// mailHog localhost http test
func sendPasswordResetEmail(to, resetURL string) error {
	m := mail.NewMessage()
	m.SetHeader("From", vars.POST_NAME)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Password Reset")
	m.SetBody("text/html", fmt.Sprintf("Click <a href='%s'>here</a> to reset your password.", resetURL))

	d := mail.NewDialer("localhost", 1025, vars.POST_NAME, vars.POST_PASS)

	if err := d.DialAndSend(m); err != nil {
		return err
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

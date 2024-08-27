package infrastructure

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

// Email and server configuration
const (
    SmtpHost      = "smtp.gmail.com"          // Correct SMTP server for Gmail
    SmtpPort      = 465                       // Port for SMTPS (SSL/TLS)
    EmailFrom     = "yordanoslegesse15@gmail.com" // Your Gmail address
    EmailPassword = "bcewmdllhervddxu"        // Your app-specific password
    ServerHost    = "http://localhost:8080"   // Change to your domain in production
    TokenTTlL      = time.Hour                 // Token Time-To-Live
)

// Generates a secure random token
func GenerateResetToken() (string, error) {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

// Sends the password reset email
func SendResetEmail(to, token string) error {
    resetLink := fmt.Sprintf("%s/api/v0/auth/reset?token=%s", ServerHost, token)
    body := fmt.Sprintf(`
        Hi,

        You requested a password reset. Click the link below to reset your password:

        %s

        If you did not request this, please ignore this email.
    `, resetLink)

    m := gomail.NewMessage()
    m.SetHeader("From", fmt.Sprintf("%s <%s>", "Eskalate G5 Blog Project", EmailFrom))
    m.SetHeader("To", to)
    m.SetHeader("Subject", "Password Reset")
    m.SetBody("text/plain", body)

    d := gomail.NewDialer(SmtpHost, SmtpPort, EmailFrom, EmailPassword)

    if err := d.DialAndSend(m); err != nil {
        return err
    }

    return nil
}

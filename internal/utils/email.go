package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func SendVerificationEmail(to, token, appURL, smtpHost, smtpPort, smtpUsername, smtpPassword string) error {
	if smtpHost == "" {
		return fmt.Errorf("SMTP not configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify your email for ErgrAcer")
	
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", appURL, token)
	body := fmt.Sprintf("Please click the following link to verify your email: %s", verifyURL)
	m.SetBody("text/plain", body)

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(smtpHost, port, smtpUsername, smtpPassword)
	return d.DialAndSend(m)
}
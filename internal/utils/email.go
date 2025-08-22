package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v5"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func SendVerificationEmail(to, token, appURL, domain, apiKey, fromEmail, fromName string) error {
	if domain == "" || apiKey == "" {
		return fmt.Errorf("Mailgun not configured")
	}

	mg := mailgun.NewMailgun(apiKey)

	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", appURL, token)

	subject := "Verify your email for ErgRacer"
	textBody := fmt.Sprintf(`
Hi there!

Welcome to ErgRacer! Please verify your email address by clicking the link below:

%s

If you didn't create an account with ErgRacer, you can safely ignore this email.

Thanks,
The ErgRacer Team
`, verifyURL)

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Verify your email</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button { background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block; margin: 20px 0; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; font-size: 14px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Welcome to ErgRacer!</h2>
        <p>Thank you for signing up. Please verify your email address by clicking the button below:</p>
        <a href="%s" class="button">Verify Email Address</a>
        <p>If the button doesn't work, you can copy and paste this link into your browser:</p>
        <p>%s</p>
        <div class="footer">
            <p>If you didn't create an account with ErgRacer, you can safely ignore this email.</p>
            <p>Thanks,<br>The ErgRacer Team</p>
        </div>
    </div>
</body>
</html>
`, verifyURL, verifyURL)

	message := mailgun.NewMessage(
		domain,
		fmt.Sprintf("%s <%s>", fromName, fromEmail),
		subject,
		textBody,
		to,
	)
	message.SetHTML(htmlBody)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := mg.Send(ctx, message)
	return err
}

package utils

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

func SendPasswordResetEmail(resetCode string, email string) error {
	// TODO: Construct domain dynamically
	resetUrl := fmt.Sprintf("http://localhost:3000/execs/password/reset/%s", resetCode)
	message := fmt.Sprintf("Forgot your password? Reset it using the following link:\n%s\n\nIf you didn't request a password reset, please ignore this email.", resetUrl)

	mailMessage := mail.NewMessage()
	mailMessage.SetHeader("From", "schooladmin@school.com")
	mailMessage.SetHeader("To", email)
	mailMessage.SetHeader("Subject", "Your password reset link")
	mailMessage.SetBody("text/plain", message)

	mailDialer := mail.NewDialer("localhost", 1025, "", "")
	err := mailDialer.DialAndSend(mailMessage)

	if err != nil {
		return ErrorHandler(err, "Couldn't send email")
	}
	return nil
}

package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/mailgun/mailgun-go/v4"
)

var mg *mailgun.MailgunImpl

// InitMailgun initializes the Mailgun client
func InitMailgun() {
	mg = mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
}

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() string {
	otp := make([]byte, 6)
	_, err := rand.Read(otp)
	if err != nil {
		return ""
	}
	for i := range otp {
		otp[i] = '0' + otp[i]%10
	}
	return string(otp)
}

// SendOTP sends an OTP to the specified email address
func SendOTP(ctx context.Context, recipient, otp string) error {
	sender := "noreply@yourdomain.com"
	subject := "Your Login OTP"
	body := fmt.Sprintf("Your OTP is: %s", otp)

	message := mg.NewMessage(sender, subject, body, recipient)

	_, _, err := mg.Send(ctx, message)
	return err
}

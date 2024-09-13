package msg

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/db"
)

type NotificationService interface {
	SendWarning(userGUID string, userClientIP string) error
	GetEmailAddressFromGUID(userGUID string) (string, error)
}

type EmailNotificationService struct {
	emailDataBase db.TokenDB
}

func (emailService *EmailNotificationService) SendWarning(userGUID string, userClientIP string) error {
	fmt.Println("Sending email...")
	to, err := emailService.GetEmailAddressFromGUID(userGUID)
	if err != nil {
		return err
	}
	subject := "Suspicious sign in from different IP"
	body := "If that was you, then feel free to ignore that email.\n" +
		"We have detected a new sign in from new IP address: " +
		userClientIP + ".\n"
	sendEmailSMTP(to, subject, body)
	return nil
}

func (emailService *EmailNotificationService) GetEmailAddressFromGUID(userGUID string) (string, error) {
	email, err := emailService.emailDataBase.GetEmailAddressFromGUID(userGUID)
	if err != nil {
		return "", errors.New("failed to find user email")
	}
	return email, nil
}

func sendEmailSMTP(to string, subject string, body string) error {
	smtpHost := os.Getenv("GOAUTH_BACKDEV_SMTP_HOST")
	smtpPort := "587"
	senderEmail := os.Getenv("GOAUTH_BACKDEV_EMAIL_USERNAME")
	password := os.Getenv("GOAUTH_BACKDEV_EMAIL_PASSWORD")

	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)

	// Format the email headers and body
	message := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" + // Empty line between headers and body
			body + "\r\n")

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

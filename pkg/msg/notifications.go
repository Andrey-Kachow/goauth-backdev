package msg

import (
	"fmt"
	"net/smtp"
	"os"
)

type NotificationService interface {
	SendWarning(userEmail string, userClientIP string) error
}

type EmailNotificationService struct {
}

func (emailService *EmailNotificationService) SendWarning(userEmail string, userClientIP string) error {
	fmt.Println("Sending email...")
	subject := "Suspicious sign in from different IP"
	body := "If that was you, then feel free to ignore that email.\n" +
		"We have detected a new sign in from new IP address: " +
		userClientIP + ".\n"
	err := sendEmailSMTP(userEmail, subject, body)
	if err != nil {
		fmt.Printf("Failed to send email to %s", userEmail)
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Sent email to %s", userEmail)
	}
	return nil
}

func sendEmailSMTP(to string, subject string, body string) error {
	smtpHost := os.Getenv("GOAUTH_BACKDEV_SMTP_HOST")
	smtpPort := "587"
	senderEmail := os.Getenv("GOAUTH_BACKDEV_EMAIL_USERNAME")
	password := os.Getenv("GOAUTH_BACKDEV_EMAIL_PASSWORD")

	auth := smtp.PlainAuth("", senderEmail, password, smtpHost)

	message := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

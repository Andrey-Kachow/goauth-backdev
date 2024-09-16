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
	host        string
	port        string
	senderEmail string
	password    string
}

func (emailService *EmailNotificationService) SendWarning(userEmail string, userClientIP string) error {
	fmt.Println("Sending email...")
	subject := "Suspicious sign in from different IP"
	body := "If that was you, then feel free to ignore that email.\n" +
		"We have detected a new sign in from new IP address: " +
		userClientIP + ".\n"
	err := emailService.sendEmailSMTP(userEmail, subject, body)
	if err != nil {
		fmt.Printf("Failed to send email to %s\n", userEmail)
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Sent email to %s\n", userEmail)
	}
	return nil
}

func (emailService *EmailNotificationService) sendEmailSMTP(recipientEmail string, subject string, body string) error {
	auth := smtp.PlainAuth("", emailService.senderEmail, emailService.password, emailService.host)

	message := []byte(
		"To: " + recipientEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n")

	err := smtp.SendMail(emailService.host+":"+emailService.port, auth, emailService.senderEmail, []string{recipientEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

type DummyNotificationService struct {
}

func (dummyService *DummyNotificationService) SendWarning(userEmail string, userClientIP string) error {
	fmt.Println("No warning email has been sent via dummy notification service")
	return nil
}

func ProvideNotificationService() NotificationService {
	requiredEnvVars := []string{
		"GOAUTH_BACKDEV_SMTP_HOST",
		"GOAUTH_BACKDEV_EMAIL_USERNAME",
		"GOAUTH_BACKDEV_EMAIL_PASSWORD",
	}
	for _, requiredEnvVar := range requiredEnvVars {
		envVar := os.Getenv(requiredEnvVar)
		if envVar == "" || envVar == "CHANGEME" {
			fmt.Printf("Environment variable %s is not set. Email notifications disabled\n", requiredEnvVar)
			return &DummyNotificationService{}
		}
	}
	return &EmailNotificationService{
		host:        os.Getenv("GOAUTH_BACKDEV_SMTP_HOST"),
		port:        "587",
		senderEmail: os.Getenv("GOAUTH_BACKDEV_EMAIL_USERNAME"),
		password:    os.Getenv("GOAUTH_BACKDEV_EMAIL_PASSWORD"),
	}
}

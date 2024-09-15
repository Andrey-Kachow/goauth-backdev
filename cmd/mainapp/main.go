package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/api"
	"github.com/Andrey-Kachow/goauth-backdev/pkg/msg"
)

func sampleClientHandler(writer http.ResponseWriter, request *http.Request) {
	htmlFile := "cmd/mainapp/sampleclient.html"
	content, err := os.ReadFile(htmlFile)
	if err != nil {
		http.Error(writer, "Unable to read HTML file", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "text/html")
	writer.Write(content)
}

func main() {
	requiredEnvVars := []string{
		"GOAUTH_BACKDEV_SMTP_HOST",
		"GOAUTH_BACKDEV_EMAIL_USERNAME",
		"GOAUTH_BACKDEV_EMAIL_PASSWORD",
	}
	for _, v := range requiredEnvVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is not set. Fail.\n", v)
		}
	}

	var service msg.EmailNotificationService = msg.EmailNotificationService{}
	service.SendWarning(os.Getenv("GOAUTH_BACKDEV_EMAIL_USERNAME"), "123.23.23.0")

	fmt.Println("Starting the app")
	http.HandleFunc("/", sampleClientHandler)
	http.HandleFunc("/api/access", api.AccessHandler)
	http.HandleFunc("/api/refresh", api.RefreshHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/api"
)

func handler(writer http.ResponseWriter, request *http.Request) {
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
	fmt.Println("Starting the app")
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/access", api.AccessHandler)
	http.HandleFunc("/api/refresh", api.RefreshHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

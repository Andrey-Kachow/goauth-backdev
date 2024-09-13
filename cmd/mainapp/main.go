package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/api"
)

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hi there, I love %s!", request.URL.Path[1:])
}

func main() {
	fmt.Println("Starting the app")
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/login", api.AccessHandler)
	http.HandleFunc("/api/refresh", api.RefreshHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

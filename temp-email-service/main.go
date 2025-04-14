package main

import (
	"bounceshield/temp-email-service/api"
	"fmt"
	"log"
	"net/http"
)

func main() {
	api.SetupRoutes()

	port := ":8082"
	fmt.Println("🚀 Temp Email Service running at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

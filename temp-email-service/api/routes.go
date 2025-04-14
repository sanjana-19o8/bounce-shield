package api

import "net/http"

func SetupRoutes() {
	http.HandleFunc("/api/create-temp-email", CreateTempEmailHandler)
	http.HandleFunc("/api/check-temp-inbox", CheckTempInboxHandler)
}

package api

import (
	"bounceshield/temp-email-service/services"
	"encoding/json"
	"net/http"
)

func CreateTempEmailHandler(w http.ResponseWriter, r *http.Request) {
	tempEmail := services.GenerateTempEmail()
	json.NewEncoder(w).Encode(tempEmail)
}

func CheckTempInboxHandler(w http.ResponseWriter, r *http.Request) {
	tempEmail := r.URL.Query().Get("tempEmail")
	emails, valid := services.CheckInbox(tempEmail)
	if !valid {
		http.Error(w, "Temp email not found or expired", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(emails)
}

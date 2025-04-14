package services

import (
	"fmt"
	"time"
)

type TempEmail struct {
	Address   string
	CreatedAt time.Time
}

var TempEmailInbox = make(map[string][]string)

func GenerateTempEmail() TempEmail {
	address := fmt.Sprintf("tempuser%d@tempemail.com", time.Now().UnixNano())
	return TempEmail{Address: address, CreatedAt: time.Now()}
}

func CheckInbox(tempEmail string) ([]string, bool) {
	emails, exists := TempEmailInbox[tempEmail]
	if !exists {
		return nil, false
	}
	return emails, true
}

func SimulateIncomingEmails() {
	for {
		time.Sleep(5 * time.Second)
		email := fmt.Sprintf("Welcome to Temp Email service at %s", time.Now())
		TempEmailInbox["tempuser1234@tempemail.com"] = append(TempEmailInbox["tempuser1234@tempemail.com"], email)
	}
}

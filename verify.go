package main

import (
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

func isValidEmailFormat(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func getDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func getMXRecords(domain string) ([]*net.MX, error) {
	return net.LookupMX(domain)
}

func verifySMTP(mxHost, email string) bool {
	conn, err := net.DialTimeout("tcp", mxHost+":25", 10*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, mxHost)
	if err != nil {
		return false
	}
	defer client.Quit()

	if err := client.Hello("localhost"); err != nil {
		return false
	}
	if err := client.Mail("test@yourdomain.com"); err != nil {
		return false
	}
	if err := client.Rcpt(email); err != nil {
		return false
	}

	return true
}

func verifyEmail(email string) string {
	if !isValidEmailFormat(email) {
		return "Invalid format"
	}

	domain := getDomain(email)
	mxRecords, err := getMXRecords(domain)
	if err != nil || len(mxRecords) == 0 {
		return "No MX records"
	}

	for _, mx := range mxRecords {
		if verifySMTP(mx.Host, email) {
			return "Valid"
		}
	}

	return "Undeliverable or blocked"
}

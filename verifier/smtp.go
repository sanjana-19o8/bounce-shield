package verifier

import (
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

// Result holds detailed classification info
type VerificationResult struct {
	Email   string `json:"email"`
	Status  string `json:"status"`  // "valid", "invalid", or "unknown"
	Reason  string `json:"reason"`  // "invalid syntax", "domain not found", etc.
	Details string `json:"details,omitempty"`
}

var disposableDomains = map[string]bool{
	"mailinator.com":  true,
	"10minutemail.com": true,
	"tempmail.com":    true,
	"guerrillamail.com": true,
	"yopmail.com":     true,
	// Add more from a real list
}

// 1. Validate email syntax
func isValidEmailSyntax(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// 2. Check disposable domain
func isDisposable(email string) bool {
	domain := strings.ToLower(strings.Split(email, "@")[1])
	return disposableDomains[domain]
}

// 3. DNS MX record check
func hasMX(domain string) (bool, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return false, err
	}
	return len(mxRecords) > 0, nil
}

// 4-5. SMTP Check
func smtpVerify(email string) (string, error) {
	domain := strings.Split(email, "@")[1]
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return "domain not found", err
	}

	mxHost := mxRecords[0].Host
	conn, err := net.DialTimeout("tcp", mxHost+":25", 10*time.Second)
	if err != nil {
		return "smtp connection failed", err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, mxHost)
	if err != nil {
		return "smtp handshake failed", err
	}
	defer client.Quit()

	client.Hello("bounce-shield.io")
	_ = client.Mail("verify@bounce-shield.io")
	if err := client.Rcpt(email); err != nil {
		if strings.Contains(err.Error(), "550") {
			return "mailbox not found", err
		}
		return "smtp rejected", err
	}

	return "ok", nil
}

// Master function
func VerifyEmail(email string) VerificationResult {
	if !isValidEmailSyntax(email) {
		return VerificationResult{
			Email:  email,
			Status: "invalid",
			Reason: "invalid syntax",
		}
	}

	if isDisposable(email) {
		return VerificationResult{
			Email:  email,
			Status: "invalid",
			Reason: "disposable domain",
		}
	}

	domain := strings.Split(email, "@")[1]
	hasMx, err := hasMX(domain)
	if err != nil || !hasMx {
		return VerificationResult{
			Email:  email,
			Status: "invalid",
			Reason: "domain not found",
			Details: err.Error(),
		}
	}

	smtpResult, smtpErr := smtpVerify(email)
	if smtpResult == "ok" {
		return VerificationResult{
			Email:  email,
			Status: "valid",
			Reason: "smtp mailbox accepted",
		}
	}

	return VerificationResult{
		Email:  email,
		Status: "invalid",
		Reason: smtpResult,
		Details: smtpErr.Error(),
	}
}

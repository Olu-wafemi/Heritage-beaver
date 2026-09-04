package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// Mailer sends transactional email via Gmail SMTP when configured,
// otherwise via the Resend HTTP API. When neither is configured it
// logs the link so local development works without delivering anything.
type Mailer struct {
	apiKey   string
	from     string
	smtpHost string
	smtpPort string
	smtpUser string
	smtpPass string
	http     *http.Client
}

func New(apiKey, from string) *Mailer { return NewWithSMTP(apiKey, from, "", "", "", "") }

func NewWithSMTP(apiKey, from, smtpHost, smtpPort, smtpUser, smtpPass string) *Mailer {
	if from == "" {
		from = "Hearthside <onboarding@resend.dev>"
	}
	if smtpPort == "" {
		smtpPort = "587"
	}
	return &Mailer{
		apiKey:   apiKey,
		from:     from,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		smtpUser: smtpUser,
		smtpPass: smtpPass,
		http:     &http.Client{Timeout: 15 * time.Second},
	}
}

func (m *Mailer) Enabled() bool { return m.apiKey != "" || m.smtpUser != "" }

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func (m *Mailer) SendVerificationEmail(ctx context.Context, to, link string) error {
	html := fmt.Sprintf(`<div style="font-family:sans-serif;max-width:480px">`+
		`<h2>Welcome to Hearthside</h2>`+
		`<p>Confirm your email to begin the telling — your family's stories are waiting.</p>`+
		`<p><a href="%s" style="background:#b4430f;color:#fff;padding:12px 24px;border-radius:999px;text-decoration:none">Confirm my email</a></p>`+
		`<p style="color:#777;font-size:13px">Or paste this link: %s</p>`+
		`<p style="color:#777;font-size:13px">The link expires in 24 hours.</p></div>`, link, link)

	if m.smtpHost != "" && m.smtpUser != "" && m.smtpPass != "" {
		if err := m.sendViaSMTP(to, "Confirm your Hearthside account", html); err == nil {
			return nil
		} else {
			log.Printf("[mail] smtp send failed for %s: %v — falling back", to, err)
		}
	}

	if m.apiKey == "" && m.smtpUser == "" {
		log.Printf("[dev] verification link for %s: %s", to, link)
		return nil
	}

	// Prefer SMTP already tried; otherwise Resend.
	if m.apiKey == "" {
		log.Printf("[dev] verification link for %s: %s", to, link)
		return nil
	}

	body, err := json.Marshal(resendRequest{
		From:    m.from,
		To:      []string{to},
		Subject: "Confirm your Hearthside account",
		HTML:    html,
	})
	if err != nil {
		return fmt.Errorf("encode email: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build email request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.apiKey)

	res, err := m.http.Do(req)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var errBody struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(res.Body).Decode(&errBody)
		if errBody.Message == "" {
			errBody.Message = res.Status
		}
		if os.Getenv("APP_ENV") != "production" {
			log.Printf("[dev] resend failed for %s, use link directly: %s", to, link)
		}
		return fmt.Errorf("resend returned %d: %s", res.StatusCode, errBody.Message)
	}
	return nil
}

func (m *Mailer) sendViaSMTP(to, subject, html string) error {
	fromAddr := m.from
	if idx := strings.LastIndex(fromAddr, "<"); idx != -1 {
		end := strings.LastIndex(fromAddr, ">")
		if end > idx {
			fromAddr = strings.TrimSpace(fromAddr[idx+1 : end])
		}
	}
	fromAddr = strings.TrimSpace(fromAddr)
	if fromAddr == "" {
		fromAddr = m.smtpUser
	}

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n", m.from, to, subject)
	msg := []byte(headers + html)

	addr := fmt.Sprintf("%s:%s", m.smtpHost, m.smtpPort)
	auth := smtp.PlainAuth("", m.smtpUser, m.smtpPass, m.smtpHost)
	return smtp.SendMail(addr, auth, fromAddr, []string{to}, msg)
}

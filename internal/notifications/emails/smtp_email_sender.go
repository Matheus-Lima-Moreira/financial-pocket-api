package emails

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"net/smtp"
	neturl "net/url"
	"strings"
)

//go:embed templates/layout/base-template.html
var baseTemplateHTML string

//go:embed templates/mails/verify-email.template.html
var verifyEmailHTML string

//go:embed templates/mails/reset-password.template.html
var resetPasswordHTML string

//go:embed images/logo.png
var logoPNG []byte

const (
	emailPrimaryColor = "#2563eb"
	companyName       = "Financial Pocket"
)

type SMTPEmailSender struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSMTPEmailSender(host, port, user, pass, from string) EmailSender {
	return &SMTPEmailSender{
		host: host,
		port: port,
		user: user,
		pass: pass,
		from: from,
	}
}

func (s *SMTPEmailSender) SendVerifyEmail(_ context.Context, to, name, verifyURL string) error {
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	address := fmt.Sprintf("%s:%s", s.host, s.port)

	subject := "Subject: Verify your email\r\n"
	headers := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	receiverName := strings.TrimSpace(name)
	if receiverName == "" {
		receiverName = "there"
	}

	body := buildVerifyEmailHTML(receiverName, verifyURL)
	message := []byte(subject + headers + body)

	return smtp.SendMail(address, auth, s.from, []string{to}, message)
}

func (s *SMTPEmailSender) SendResetPasswordEmail(_ context.Context, to, name, resetPasswordURL string) error {
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	address := fmt.Sprintf("%s:%s", s.host, s.port)

	subject := "Subject: Reset your password\r\n"
	headers := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	receiverName := strings.TrimSpace(name)
	if receiverName == "" {
		receiverName = "there"
	}

	body := buildResetPasswordHTML(receiverName, resetPasswordURL)
	message := []byte(subject + headers + body)

	return smtp.SendMail(address, auth, s.from, []string{to}, message)
}

func buildVerifyEmailHTML(name, verifyURL string) string {
	content := verifyEmailHTML
	content = replacePlaceholders(content, map[string]string{
		"{{NAME}}":          name,
		"{{URL}}":           verifyURL,
		"{{PRIMARY_COLOR}}": emailPrimaryColor,
		"{{COMPANY_NAME}}":  companyName,
	})

	logoTag := `<img src="data:image/png;base64,` + base64.StdEncoding.EncodeToString(logoPNG) + `" alt="` + companyName + `" style="width: 250px; height: auto;" />`

	html := baseTemplateHTML
	html = replacePlaceholders(html, map[string]string{
		"{{CONTENT}}":         content,
		"{{LOGO}}":            logoTag,
		"{{PRIMARY_COLOR}}":   emailPrimaryColor,
		"{{COMPANY_NAME}}":    companyName,
		"{{GLOBAL_SITE_URL}}": resolveGlobalSiteURL(verifyURL),
	})

	return html
}

func buildResetPasswordHTML(name, resetPasswordURL string) string {
	content := resetPasswordHTML
	content = replacePlaceholders(content, map[string]string{
		"{{NAME}}":          name,
		"{{URL}}":           resetPasswordURL,
		"{{PRIMARY_COLOR}}": emailPrimaryColor,
		"{{COMPANY_NAME}}":  companyName,
	})

	logoTag := `<img src="data:image/png;base64,` + base64.StdEncoding.EncodeToString(logoPNG) + `" alt="` + companyName + `" style="width: 250px; height: auto;" />`

	html := baseTemplateHTML
	html = replacePlaceholders(html, map[string]string{
		"{{CONTENT}}":         content,
		"{{LOGO}}":            logoTag,
		"{{PRIMARY_COLOR}}":   emailPrimaryColor,
		"{{COMPANY_NAME}}":    companyName,
		"{{GLOBAL_SITE_URL}}": resolveGlobalSiteURL(resetPasswordURL),
	})

	return html
}

func replacePlaceholders(html string, values map[string]string) string {
	result := html
	for key, value := range values {
		result = strings.ReplaceAll(result, key, value)
	}
	return result
}

func resolveGlobalSiteURL(verifyURL string) string {
	parsedURL, err := neturl.Parse(verifyURL)
	if err != nil {
		return ""
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ""
	}
	return parsedURL.Scheme + "://" + parsedURL.Host
}

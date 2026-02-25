package messages

import goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

var en = []*goi18n.Message{
	{ID: "user.created", Other: "User created successfully"},
	{ID: "user.listed", Other: "Users listed successfully"},
	{ID: "user.details", Other: "User details"},
	{ID: "auth.login_success", Other: "Login successful"},
	{ID: "auth.refresh_success", Other: "Token refreshed successfully"},
	{ID: "auth.verify_email_sent", Other: "Verification email sent successfully"},
	{ID: "auth.verify_email_success", Other: "Email verified successfully"},
	{ID: "auth.verify_email_send_failed", Other: "Failed to send verification email"},
	{ID: "auth.email_not_verified", Other: "Email not verified"},
	{ID: "error.missing_body", Other: "Missing request body"},
	{ID: "error.invalid_json", Other: "Invalid JSON"},
	{ID: "error.validation", Other: "Validation error"},
	{ID: "error.internal", Other: "Internal server error"},
}

func GetENMessages() []*goi18n.Message {
	return en
}

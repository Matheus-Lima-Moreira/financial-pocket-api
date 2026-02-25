package auth

import "context"

type EmailSender interface {
	SendVerifyEmail(ctx context.Context, to, name, verifyURL string) error
}

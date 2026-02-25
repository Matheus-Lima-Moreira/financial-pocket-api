package emails

import "context"

type EmailSender interface {
	SendVerifyEmail(ctx context.Context, to, name, verifyURL string) error
	SendResetPasswordEmail(ctx context.Context, to, name, resetPasswordURL string) error
}

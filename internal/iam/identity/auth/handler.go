package auth

import (
	"net/http"
	"time"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/dtos"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service     *Service
	rateLimiter *AuthRateLimiter
}

func NewHandler(service *Service, rateLimiter *AuthRateLimiter) *Handler {
	return &Handler{
		service:     service,
		rateLimiter: rateLimiter,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var request RegisterRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	if !h.handleRateLimit(c, AuthRateLimitRegister, request.Email) {
		return
	}

	if err := h.service.Register(c.Request.Context(), RegisterInputDTO{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	},
	); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dtos.ReplyDTO{
		Message: "auth.verify_email_sent",
		Data:    nil,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var request LoginRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	if !h.handleRateLimit(c, AuthRateLimitLogin, request.Email) {
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Data:    tokens,
		Message: "auth.login_success",
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var request RefreshRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	tokens, err := h.service.RefreshToken(c.Request.Context(), request.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Data:    tokens,
		Message: "auth.refresh_success",
	})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var request ResetPasswordRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), request.Token, request.NewPassword); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Message: "auth.reset_password_success",
		Data:    nil,
	})
}

func (h *Handler) ResendVerificationEmail(c *gin.Context) {
	var request ResendVerificationEmailRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	if !h.handleRateLimit(c, AuthRateLimitResendVerificationEmail, request.Email) {
		return
	}

	if err := h.service.SendVerificationEmail(c.Request.Context(), request.Email); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dtos.ReplyDTO{
		Message: "auth.resend_verification_email_success",
		Data:    nil,
	})
}

func (h *Handler) SendResetPasswordEmail(c *gin.Context) {
	var request SendResetPasswordEmailRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(err)
		return
	}

	if !h.handleRateLimit(c, AuthRateLimitSendResetPassword, request.Email) {
		return
	}

	if err := h.service.SendResetPasswordEmail(c.Request.Context(), request.Email); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dtos.ReplyDTO{
		Message: "auth.reset_password_email_sent_success",
		Data:    nil,
	})
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var request VerifyEmailRequestDTO

	if err := c.ShouldBindQuery(&request); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), request.Token); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Message: "auth.verify_email_success",
		Data:    nil,
	})
}

func (h *Handler) handleRateLimit(c *gin.Context, action AuthRateLimitAction, identity string) bool {
	allowed, retryAfter := h.rateLimiter.Allow(action, identity, c.ClientIP(), time.Now())
	if allowed {
		return true
	}

	c.JSON(http.StatusTooManyRequests, dtos.ReplyDTO{
		Message: "auth.rate_limited",
		Data: gin.H{
			"retry_after_seconds": int(retryAfter.Seconds()),
		},
	})
	return false
}

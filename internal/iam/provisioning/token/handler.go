package token

import (
	"net/http"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/dtos"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type resendVerificationEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) ResendVerificationEmail(c *gin.Context) {
	var input resendVerificationEmailInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.SendVerificationEmail(c.Request.Context(), input.Email); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dtos.ReplyDTO{
		Message: "token.resend_verification_email_success",
		Data:    nil,
	})
}

type resetPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var input resetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), input.Email); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dtos.ReplyDTO{
		Message: "token.reset_password_success",
		Data:    nil,
	})
}

type verifyEmailInput struct {
	Token string `form:"token" binding:"required"`
}

func (h *Handler) VerifyEmail(c *gin.Context) {
	var input verifyEmailInput

	if err := c.ShouldBindQuery(&input); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), input.Token); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Message: "auth.verify_email_success",
		Data:    nil,
	})
}

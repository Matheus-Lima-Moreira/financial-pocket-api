package auth

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

type registerInput struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type refreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type verifyEmailInput struct {
	Token string `form:"token" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var input registerInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Register(c.Request.Context(), RegisterInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
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
	var input registerInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), input.Email, input.Password)
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
	var input refreshInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	tokens, err := h.service.RefreshToken(c.Request.Context(), input.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Data:    tokens,
		Message: "auth.refresh_success",
	})
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

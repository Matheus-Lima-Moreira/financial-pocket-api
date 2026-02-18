package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type registerInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type refreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var input registerInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Register(c.Request.Context(), input.Email, input.Password); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
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

	c.JSON(http.StatusOK, tokens)
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

	c.JSON(http.StatusOK, tokens)
}

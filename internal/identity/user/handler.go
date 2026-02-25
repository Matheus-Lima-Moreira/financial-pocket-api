package user

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
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type listInput struct {
	Page int `form:"page" binding:"required,min=1"`
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

	c.JSON(http.StatusCreated, dtos.ReplyDTO{
		Message: "user created successfully",
		Data:    nil,
	})
}

func (h *Handler) List(c *gin.Context) {
	var input listInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.Error(err)
		return
	}

	users, pagination, err := h.service.List(c.Request.Context(), input.Page)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Data:       users,
		Pagination: pagination,
		Message:    "users listed successfully",
	})
}

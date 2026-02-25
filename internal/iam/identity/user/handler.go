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

type listInput struct {
	Page int `form:"page" binding:"required,min=1"`
}

type detailsInput struct {
	ID uint `uri:"id" binding:"required,min=1"`
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
		Message:    "user.listed",
	})
}

func (h *Handler) Details(c *gin.Context) {
	var input detailsInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.Error(err)
		return
	}

	user, err := h.service.Details(c.Request.Context(), input.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Data:    user,
		Message: "user.details",
	})
}

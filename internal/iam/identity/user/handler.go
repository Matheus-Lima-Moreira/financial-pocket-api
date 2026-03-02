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

func (h *Handler) List(c *gin.Context) {
	var request ListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.Error(err)
		return
	}

	users, pagination, err := h.service.List(c.Request.Context(), request.Page)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Code:       dtos.SUCCESS,
		Data:       users,
		Pagination: pagination,
		Message:    "user.listed",
	})
}

func (h *Handler) Details(c *gin.Context) {
	var request DetailsRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.Error(err)
		return
	}

	user, err := h.service.Details(c.Request.Context(), request.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Code:    dtos.SUCCESS,
		Data:    user,
		Message: "user.details",
	})
}

func (h *Handler) Profile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dtos.ReplyDTO{
		Code:    dtos.SUCCESS,
		Data:    user,
		Message: "user.profile-details",
	})
}

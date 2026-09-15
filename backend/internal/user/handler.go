package user

import (
	apperrors "backend/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	users, err := h.service.GetAll()
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		apperrors.Handle(c, err)
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetMe(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		apperrors.Handle(c, apperrors.NewUnauthorized("unauthenticated"))
		return
	}

	user, err := h.service.GetByID(userID.(uint))
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		apperrors.Handle(c, err)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Handle(c, apperrors.NewInvalidInput(err.Error()))
		return
	}

	user, err := h.service.Update(id, req)
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateMe(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		apperrors.Handle(c, apperrors.NewUnauthorized("unauthenticated"))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Handle(c, apperrors.NewInvalidInput(err.Error()))
		return
	}

	user, err := h.service.Update(userID.(uint), req)
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		apperrors.Handle(c, err)
		return
	}

	if err := h.service.Delete(id); err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, apperrors.NewInvalidInput("invalid user id")
	}
	return uint(id), nil
}

package device

import (
	apperrors "backend/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	devices, err := h.service.GetAll()
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		apperrors.Handle(c, err)
		return
	}

	device, err := h.service.GetByID(id)
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, device)
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		apperrors.Handle(c, apperrors.NewUnauthorized("unauthenticated"))
		return
	}

	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Handle(c, apperrors.NewInvalidInput(err.Error()))
		return
	}

	device, err := h.service.Create(userID.(uuid.UUID), req)
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusCreated, device)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		apperrors.Handle(c, err)
		return
	}

	var req UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Handle(c, apperrors.NewInvalidInput(err.Error()))
		return
	}

	device, err := h.service.Update(id, req)
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, device)
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

func parseID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, apperrors.NewInvalidInput("invalid device id")
	}
	return id, nil
}

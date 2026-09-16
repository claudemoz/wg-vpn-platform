package server

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
	servers, err := h.service.GetAll()
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, servers)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		apperrors.Handle(c, err)
		return
	}

	server, err := h.service.GetByID(id)
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, server)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Handle(c, apperrors.NewInvalidInput(err.Error()))
		return
	}

	server, err := h.service.Create(req)
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusCreated, server)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		apperrors.Handle(c, err)
		return
	}

	var req UpdateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Handle(c, apperrors.NewInvalidInput(err.Error()))
		return
	}

	server, err := h.service.Update(id, req)
	if err != nil {
		apperrors.Handle(c, err)
		return
	}
	c.JSON(http.StatusOK, server)
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
		return uuid.Nil, apperrors.NewInvalidInput("invalid server id")
	}
	return id, nil
}

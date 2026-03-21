package handler

import (
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FineHandler struct {
	service service.FineService
}

func NewFineHandler(service service.FineService) *FineHandler {
	return &FineHandler{service: service}
}

func (h *FineHandler) ListFines(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	fines, total, err := h.service.ListFines(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  fines,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *FineHandler) GetFine(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fine ID"})
		return
	}

	fine, err := h.service.GetFine(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fine not found"})
		return
	}

	c.JSON(http.StatusOK, fine)
}

func (h *FineHandler) CreateFine(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.CreateFineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fine, err := h.service.CreateFine(userID, role, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, fine)
}


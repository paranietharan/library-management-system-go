package handler

import (
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LendingHandler struct {
	service service.LendingService
}

func NewLendingHandler(service service.LendingService) *LendingHandler {
	return &LendingHandler{service: service}
}

func (h *LendingHandler) ListLendings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	lendings, total, err := h.service.ListLendings(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  lendings,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *LendingHandler) GetLending(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lending ID"})
		return
	}

	lending, err := h.service.GetLending(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lending not found"})
		return
	}

	c.JSON(http.StatusOK, lending)
}

func (h *LendingHandler) CreateLending(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.CreateLendingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lending, err := h.service.CreateLending(userID, role, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lending)
}

func (h *LendingHandler) UpdateLending(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lending ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.UpdateLendingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lending, err := h.service.UpdateLending(userID, role, uint(id), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lending)
}

func (h *LendingHandler) DeleteLending(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lending ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	if err := h.service.DeleteLending(userID, role, uint(id)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lending deleted successfully"})
}


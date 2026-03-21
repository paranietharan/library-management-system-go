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

// @Summary List lendings
// @Tags lendings
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param search query string false "Search"
// @Success 200 {object} gin.H
// @Router /lendings [get]
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

// @Summary Get lending
// @Tags lendings
// @Security BearerAuth
// @Produce json
// @Param id path int true "Lending ID"
// @Success 200 {object} gin.H
// @Router /lendings/{id} [get]
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

// @Summary Create lending
// @Tags lendings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateLendingRequest true "Create lending request"
// @Success 201 {object} gin.H
// @Router /lendings [post]
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

// @Summary Update lending
// @Tags lendings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Lending ID"
// @Param request body dto.UpdateLendingRequest true "Update lending request"
// @Success 200 {object} gin.H
// @Router /lendings/{id} [put]
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

// @Summary Delete lending
// @Tags lendings
// @Security BearerAuth
// @Produce json
// @Param id path int true "Lending ID"
// @Success 200 {object} gin.H
// @Router /lendings/{id} [delete]
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

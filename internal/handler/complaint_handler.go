package handler

import (
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ComplaintHandler struct {
	service service.ComplaintService
}

func NewComplaintHandler(service service.ComplaintService) *ComplaintHandler {
	return &ComplaintHandler{service: service}
}

// @Summary List complaints
// @Tags complaints
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param search query string false "Search"
// @Success 200 {object} gin.H
// @Router /complaints [get]
func (h *ComplaintHandler) ListComplaints(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	complaints, total, err := h.service.ListComplaints(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  complaints,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// @Summary Get complaint
// @Tags complaints
// @Security BearerAuth
// @Produce json
// @Param id path int true "Complaint ID"
// @Success 200 {object} gin.H
// @Router /complaints/{id} [get]
func (h *ComplaintHandler) GetComplaint(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid complaint ID"})
		return
	}

	complaint, err := h.service.GetComplaint(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Complaint not found"})
		return
	}

	c.JSON(http.StatusOK, complaint)
}

// @Summary Create complaint
// @Tags complaints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateComplaintRequest true "Create complaint request"
// @Success 201 {object} gin.H
// @Router /complaints [post]
func (h *ComplaintHandler) CreateComplaint(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateComplaintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	complaint, err := h.service.CreateComplaint(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, complaint)
}

// @Summary Update complaint
// @Tags complaints
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Complaint ID"
// @Param request body dto.UpdateComplaintRequest true "Update complaint request"
// @Success 200 {object} gin.H
// @Router /complaints/{id} [put]
func (h *ComplaintHandler) UpdateComplaint(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid complaint ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.UpdateComplaintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	complaint, err := h.service.UpdateComplaint(userID, role, uint(id), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, complaint)
}

// @Summary Delete complaint
// @Tags complaints
// @Security BearerAuth
// @Produce json
// @Param id path int true "Complaint ID"
// @Success 200 {object} gin.H
// @Router /complaints/{id} [delete]
func (h *ComplaintHandler) DeleteComplaint(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid complaint ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	if err := h.service.DeleteComplaint(userID, role, uint(id)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Complaint deleted successfully"})
}

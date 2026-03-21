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


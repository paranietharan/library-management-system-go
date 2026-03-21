package handler

import (
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	service service.ReservationService
}

func NewReservationHandler(service service.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

// @Summary List reservations
// @Tags reservations
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param search query string false "Search"
// @Success 200 {object} gin.H
// @Router /reservations [get]
func (h *ReservationHandler) ListReservations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	reservations, total, err := h.service.ListReservations(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  reservations,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// @Summary Get reservation
// @Tags reservations
// @Security BearerAuth
// @Produce json
// @Param id path int true "Reservation ID"
// @Success 200 {object} gin.H
// @Router /reservations/{id} [get]
func (h *ReservationHandler) GetReservation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	reservation, err := h.service.GetReservation(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// @Summary Create reservation
// @Tags reservations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateReservationRequest true "Create reservation request"
// @Success 201 {object} gin.H
// @Router /reservations [post]
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := h.service.CreateReservation(userID, role, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// @Summary Update reservation
// @Tags reservations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID"
// @Param request body dto.UpdateReservationRequest true "Update reservation request"
// @Success 200 {object} gin.H
// @Router /reservations/{id} [put]
func (h *ReservationHandler) UpdateReservation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.UpdateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reservation, err := h.service.UpdateReservation(userID, role, uint(id), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// @Summary Delete reservation
// @Tags reservations
// @Security BearerAuth
// @Produce json
// @Param id path int true "Reservation ID"
// @Success 200 {object} gin.H
// @Router /reservations/{id} [delete]
func (h *ReservationHandler) DeleteReservation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	if err := h.service.DeleteReservation(userID, role, uint(id)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted successfully"})
}

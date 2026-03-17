package http

import (
	"net/http"

	"go-challenge-agenda/services/api/internal/domain"
	"go-challenge-agenda/services/api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AvailabilityHandler struct {
	uc *usecase.AvailabilityUsecase
}

func NewAvailabilityHandler(uc *usecase.AvailabilityUsecase) *AvailabilityHandler {
	return &AvailabilityHandler{uc: uc}
}

// Get godoc
// @Summary     Get available slots for a doctor on a given date
// @Tags        availability
// @Produce     json
// @Param       id    path      string  true   "Doctor ID"
// @Param       date  query     string  true   "Date (YYYY-MM-DD)"
// @Param       type  query     string  false  "Reservation type: first_visit or follow_up"
// @Success     200   {object}  domain.AvailabilityResponse
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /doctors/{id}/availability [get]
func (h *AvailabilityHandler) Get(c *gin.Context) {
	doctorID := c.Param("id")

	var req domain.GetAvailabilityRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.uc.GetAvailability(c.Request.Context(), doctorID, req.Date, req.Type)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

package http

import (
	"net/http"
	"time"

	agendav1 "go-challenge-agenda/gen/agenda/v1"
	"go-challenge-agenda/services/api/internal/domain"
	"go-challenge-agenda/services/api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	uc     *usecase.ReservationUsecase
	agenda usecase.AgendaPort
}

func NewReservationHandler(uc *usecase.ReservationUsecase, agenda usecase.AgendaPort) *ReservationHandler {
	return &ReservationHandler{uc: uc, agenda: agenda}
}

// Create godoc
// @Summary     Create a reservation
// @Tags        reservations
// @Accept      json
// @Produce     json
// @Param       body  body      domain.CreateReservationRequest  true  "Reservation request"
// @Success     201   {object}  domain.ReservationResponse
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /reservations [post]
func (h *ReservationHandler) Create(c *gin.Context) {
	var req domain.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.Create(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Get godoc
// @Summary     Get a reservation by ID
// @Tags        reservations
// @Produce     json
// @Param       id   path      string  true  "Reservation ID"
// @Success     200  {object}  domain.ReservationResponse
// @Failure     404  {object}  map[string]string
// @Router      /reservations/{id} [get]
func (h *ReservationHandler) Get(c *gin.Context) {
	resp, err := h.agenda.GetReservation(c.Request.Context(), &agendav1.GetReservationRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, protoReservationToDTO(resp.Reservation))
}

// List godoc
// @Summary     List reservations for a doctor in a time range
// @Tags        reservations
// @Produce     json
// @Param       doctor_id  query     string  false  "Doctor ID"
// @Param       from       query     string  true   "From datetime (RFC3339)"
// @Param       to         query     string  true   "To datetime (RFC3339)"
// @Success     200        {array}   domain.ReservationResponse
// @Failure     400        {object}  map[string]string
// @Router      /reservations [get]
func (h *ReservationHandler) List(c *gin.Context) {
	doctorID := c.Query("doctor_id")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params are required"})
		return
	}
	if _, err := time.Parse(time.RFC3339, fromStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from"})
		return
	}

	resp, err := h.agenda.ListReservations(c.Request.Context(), &agendav1.ListReservationsRequest{
		DoctorId: doctorID,
		From:     fromStr,
		To:       toStr,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	var list []domain.ReservationResponse
	for _, r := range resp.Reservations {
		list = append(list, *protoReservationToDTO(r))
	}
	c.JSON(http.StatusOK, list)
}

// Cancel godoc
// @Summary     Cancel a reservation
// @Tags        reservations
// @Param       id   path  string  true  "Reservation ID"
// @Success     204
// @Failure     500  {object}  map[string]string
// @Router      /reservations/{id} [delete]
func (h *ReservationHandler) Cancel(c *gin.Context) {
	if err := h.uc.Cancel(c.Request.Context(), c.Param("id")); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func protoReservationToDTO(r *agendav1.Reservation) *domain.ReservationResponse {
	if r == nil {
		return nil
	}
	typeStr := "follow_up"
	if r.Type == agendav1.ReservationType_RESERVATION_TYPE_FIRST_VISIT {
		typeStr = "first_visit"
	}
	statusStr := "confirmed"
	if r.Status == agendav1.ReservationStatus_RESERVATION_STATUS_CANCELLED {
		statusStr = "cancelled"
	}
	return &domain.ReservationResponse{
		ID: r.Id, DoctorID: r.DoctorId, PatientID: r.PatientId,
		StartsAt: r.StartsAt, EndsAt: r.EndsAt, Type: typeStr, Status: statusStr,
	}
}

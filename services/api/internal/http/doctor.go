package http

import (
	"net/http"

	"go-challenge-agenda/services/api/internal/domain"
	"go-challenge-agenda/services/api/internal/usecase"

	agendav1 "go-challenge-agenda/gen/agenda/v1"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	agendaPort usecase.AgendaPort
}

func NewDoctorHandler(agendaPort usecase.AgendaPort) *DoctorHandler {
	return &DoctorHandler{agendaPort: agendaPort}
}

// List godoc
// @Summary     List all doctors
// @Tags        doctors
// @Produce     json
// @Success     200  {array}   domain.DoctorResponse
// @Failure     500  {object}  map[string]string
// @Router      /doctors [get]
func (h *DoctorHandler) List(c *gin.Context) {
	resp, err := h.agendaPort.ListDoctors(c.Request.Context(), &agendav1.ListDoctorsRequest{})
	if err != nil {
		_ = c.Error(err)
		return
	}
	var doctors []domain.DoctorResponse
	for _, d := range resp.Doctors {
		doctors = append(doctors, protoDoctorToDTO(d))
	}
	c.JSON(http.StatusOK, doctors)
}

// Get godoc
// @Summary     Get a doctor by ID
// @Tags        doctors
// @Produce     json
// @Param       id   path      string  true  "Doctor ID"
// @Success     200  {object}  domain.DoctorResponse
// @Failure     404  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /doctors/{id} [get]
func (h *DoctorHandler) Get(c *gin.Context) {
	resp, err := h.agendaPort.GetDoctor(c.Request.Context(), &agendav1.GetDoctorRequest{Id: c.Param("id")})
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, protoDoctorToDTO(resp.Doctor))
}

func protoDoctorToDTO(d *agendav1.Doctor) domain.DoctorResponse {
	resp := domain.DoctorResponse{ID: d.Id, Name: d.Name, Specialty: d.Specialty}
	for _, wh := range d.WorkingHours {
		resp.WorkingHours = append(resp.WorkingHours, domain.WorkingHoursResp{
			Weekday: int(wh.Weekday),
			From:    wh.From,
			To:      wh.To,
		})
	}
	return resp
}

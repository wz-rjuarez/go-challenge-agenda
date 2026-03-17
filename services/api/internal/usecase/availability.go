package usecase

import (
	"context"

	agendav1 "go-challenge-agenda/gen/agenda/v1"
	"go-challenge-agenda/services/api/internal/domain"
)

// AgendaPort defines the interface for communicating with the agenda service.
// This abstracts away the concrete gRPC client implementation.
type AgendaPort interface {
	GetAvailability(ctx context.Context, req *agendav1.GetAvailabilityRequest) (*agendav1.GetAvailabilityResponse, error)
	CreateReservation(ctx context.Context, req *agendav1.CreateReservationRequest) (*agendav1.CreateReservationResponse, error)
	CancelReservation(ctx context.Context, req *agendav1.CancelReservationRequest) (*agendav1.CancelReservationResponse, error)
	GetReservation(ctx context.Context, req *agendav1.GetReservationRequest) (*agendav1.GetReservationResponse, error)
	ListReservations(ctx context.Context, req *agendav1.ListReservationsRequest) (*agendav1.ListReservationsResponse, error)
	ListDoctors(ctx context.Context, req *agendav1.ListDoctorsRequest) (*agendav1.ListDoctorsResponse, error)
	GetDoctor(ctx context.Context, req *agendav1.GetDoctorRequest) (*agendav1.GetDoctorResponse, error)
	ListPatients(ctx context.Context, req *agendav1.ListPatientsRequest) (*agendav1.ListPatientsResponse, error)
	GetPatient(ctx context.Context, req *agendav1.GetPatientRequest) (*agendav1.GetPatientResponse, error)
	CreatePatient(ctx context.Context, req *agendav1.CreatePatientRequest) (*agendav1.CreatePatientResponse, error)
	UpdatePatient(ctx context.Context, req *agendav1.UpdatePatientRequest) (*agendav1.UpdatePatientResponse, error)
	DeletePatient(ctx context.Context, req *agendav1.DeletePatientRequest) (*agendav1.DeletePatientResponse, error)
}

// AvailabilityUsecase handles availability-related business logic.
type AvailabilityUsecase struct {
	agenda AgendaPort
}

func NewAvailabilityUsecase(agenda AgendaPort) *AvailabilityUsecase {
	return &AvailabilityUsecase{agenda: agenda}
}

func (u *AvailabilityUsecase) GetAvailability(ctx context.Context, doctorID, date, resType string) (*domain.AvailabilityResponse, error) {
	pbType := agendav1.ReservationType_RESERVATION_TYPE_FOLLOW_UP
	if resType == "first_visit" {
		pbType = agendav1.ReservationType_RESERVATION_TYPE_FIRST_VISIT
	}

	resp, err := u.agenda.GetAvailability(ctx, &agendav1.GetAvailabilityRequest{
		DoctorId:        doctorID,
		Date:            date,
		ReservationType: pbType,
	})
	if err != nil {
		return nil, err
	}

	result := &domain.AvailabilityResponse{}
	for _, s := range resp.Slots {
		result.Slots = append(result.Slots, domain.AvailableSlot{StartsAt: s.StartsAt, EndsAt: s.EndsAt})
	}
	for _, r := range resp.FreeRanges {
		result.FreeRanges = append(result.FreeRanges, domain.TimeRange{From: r.From, To: r.To})
	}
	return result, nil
}

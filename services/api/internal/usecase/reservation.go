package usecase

import (
	"context"
	"fmt"
	"time"

	agendav1 "go-challenge-agenda/gen/agenda/v1"
	"go-challenge-agenda/services/api/internal/domain"
)

// ReservationUsecase handles reservation-related business logic.
type ReservationUsecase struct {
	agenda AgendaPort
}

func NewReservationUsecase(agenda AgendaPort) *ReservationUsecase {
	return &ReservationUsecase{agenda: agenda}
}

func (u *ReservationUsecase) Create(ctx context.Context, req *domain.CreateReservationRequest) (*domain.ReservationResponse, error) {
	// Validate starts_at
	if _, err := time.Parse(time.RFC3339, req.StartsAt); err != nil {
		return nil, fmt.Errorf("invalid starts_at: %w", err)
	}

	pbType := agendav1.ReservationType_RESERVATION_TYPE_FOLLOW_UP
	if req.Type == "first_visit" {
		pbType = agendav1.ReservationType_RESERVATION_TYPE_FIRST_VISIT
	}

	resp, err := u.agenda.CreateReservation(ctx, &agendav1.CreateReservationRequest{
		DoctorId:     req.DoctorID,
		StartsAt:     req.StartsAt,
		Type:         pbType,
		PatientId:    req.PatientID,
		PatientName:  req.PatientName,
		PatientPhone: req.PatientPhone,
		PatientEmail: req.PatientEmail,
	})
	if err != nil {
		return nil, fmt.Errorf("agenda.CreateReservation: %w", err)
	}

	return protoReservationToDTO(resp.Reservation), nil
}

func (u *ReservationUsecase) Cancel(ctx context.Context, id string) error {
	_, err := u.agenda.CancelReservation(ctx, &agendav1.CancelReservationRequest{Id: id})
	return err
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
		ID:        r.Id,
		DoctorID:  r.DoctorId,
		PatientID: r.PatientId,
		StartsAt:  r.StartsAt,
		EndsAt:    r.EndsAt,
		Type:      typeStr,
		Status:    statusStr,
	}
}

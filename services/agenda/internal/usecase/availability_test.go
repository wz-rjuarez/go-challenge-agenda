package usecase_test

import (
	"context"
	"testing"
	"time"

	"go-challenge-agenda/services/agenda/internal/domain"
	"go-challenge-agenda/services/agenda/internal/domain/mocks"
	"go-challenge-agenda/services/agenda/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mondayDoctor() *domain.Doctor {
	return &domain.Doctor{
		ID:        "doc-001",
		Name:      "Dr. Test",
		Specialty: "General",
		WorkingHours: []domain.WorkingHours{
			{Weekday: domain.Monday, From: "09:00", To: "17:00"},
		},
	}
}

func nextMonday() time.Time {
	t := time.Now().UTC()
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, 1)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func TestGetAvailability_HappyPath(t *testing.T) {
	date := nextMonday()
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, time.UTC)
	dayEnd := time.Date(date.Year(), date.Month(), date.Day(), 17, 0, 0, 0, time.UTC)

	doctorRepo := mocks.NewDoctorRepository(t)
	reservationRepo := mocks.NewReservationRepository(t)
	blockedSlotRepo := mocks.NewBlockedSlotRepository(t)

	doctorRepo.EXPECT().GetDoctor(context.Background(), "doc-001").Return(mondayDoctor(), nil)
	reservationRepo.EXPECT().ListReservations(context.Background(), "doc-001", dayStart, dayEnd).Return(nil, nil)
	blockedSlotRepo.EXPECT().ListBlockedSlots(context.Background(), "doc-001", dayStart, dayEnd).Return(nil, nil)

	uc := usecase.NewAvailabilityUsecase(doctorRepo, reservationRepo, blockedSlotRepo)

	result, err := uc.GetAvailability(context.Background(), "doc-001", date, domain.ReservationTypeFollowUp)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Slots)
	assert.NotEmpty(t, result.FreeRanges)
}

// TestGetAvailability_WithBlockedSlots verifies that blocked slots remove time from availability.
// This test FAILS because blocked slots are not yet factored into the availability calculation.
func TestGetAvailability_WithBlockedSlots(t *testing.T) {
	date := nextMonday()
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, time.UTC)
	dayEnd := time.Date(date.Year(), date.Month(), date.Day(), 17, 0, 0, 0, time.UTC)
	blockedStart := time.Date(date.Year(), date.Month(), date.Day(), 10, 0, 0, 0, time.UTC)
	blockedEnd := time.Date(date.Year(), date.Month(), date.Day(), 11, 0, 0, 0, time.UTC)

	doctorRepo := mocks.NewDoctorRepository(t)
	reservationRepo := mocks.NewReservationRepository(t)
	blockedSlotRepo := mocks.NewBlockedSlotRepository(t)

	doctorRepo.EXPECT().GetDoctor(context.Background(), "doc-001").Return(mondayDoctor(), nil)
	reservationRepo.EXPECT().ListReservations(context.Background(), "doc-001", dayStart, dayEnd).Return(nil, nil)
	blockedSlotRepo.EXPECT().ListBlockedSlots(context.Background(), "doc-001", dayStart, dayEnd).Return([]*domain.BlockedSlot{
		{StartsAt: blockedStart, EndsAt: blockedEnd},
	}, nil)

	uc := usecase.NewAvailabilityUsecase(doctorRepo, reservationRepo, blockedSlotRepo)

	result, err := uc.GetAvailability(context.Background(), "doc-001", date, domain.ReservationTypeFollowUp)
	require.NoError(t, err)

	blocked := &domain.BlockedSlot{StartsAt: blockedStart, EndsAt: blockedEnd}
	_ = blocked

	for _, slot := range result.Slots {
		if !slot.StartsAt.Before(blockedStart) && slot.StartsAt.Before(blockedEnd) {
			t.Errorf("slot %v falls within blocked period [%v, %v]", slot.StartsAt, blockedStart, blockedEnd)
		}
	}
}

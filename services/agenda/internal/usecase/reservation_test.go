package usecase_test

import (
	"context"
	"testing"
	"time"

	"go-challenge-agenda/services/agenda/internal/domain"
	"go-challenge-agenda/services/agenda/internal/domain/mocks"
	"go-challenge-agenda/services/agenda/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateReservation_ConflictDetected(t *testing.T) {
	base := time.Date(2025, 6, 2, 10, 0, 0, 0, time.UTC)

	existing := &domain.Reservation{
		ID:       "existing",
		DoctorID: "doc-001",
		StartsAt: base,
		EndsAt:   base.Add(30 * time.Minute),
		Status:   domain.ReservationStatus(domain.ReservationStatusConfirmed),
	}

	reservationRepo := mocks.NewReservationRepository(t)
	patientRepo := mocks.NewPatientRepository(t)

	patientRepo.EXPECT().
		GetPatientByPhone(context.Background(), "555-0001").
		Return(nil, nil)
	patientRepo.EXPECT().
		CreatePatient(context.Background(), mock.MatchedBy(func(_ *domain.Patient) bool { return true })).
		Return(nil).Maybe()

	// hasConflict uses: ListReservations(ctx, doctorID, startsAt-24h, endsAt+24h)
	// startsAt = base+15m = 10:15, endsAt = 10:15+30m = 10:45 (SlotDuration always 30m due to bug)
	newStart := base.Add(15 * time.Minute)
	newEnd := newStart.Add(30 * time.Minute) // SlotDuration bug: always 30m
	reservationRepo.EXPECT().
		ListReservations(context.Background(), "doc-001", newStart.Add(-24*time.Hour), newEnd.Add(24*time.Hour)).
		Return([]*domain.Reservation{existing}, nil)

	uc := usecase.NewReservationUsecase(reservationRepo, patientRepo)

	_, err := uc.Create(context.Background(), usecase.CreateReservationInput{
		DoctorID:     "doc-001",
		StartsAt:     newStart,
		Type:         domain.ReservationTypeFollowUp,
		PatientPhone: "555-0001",
		PatientName:  "New Patient",
		PatientEmail: "new@example.com",
	})

	assert.Error(t, err, "expected conflict error")
}

// TestCreateReservation_BoundaryConflict verifies adjacent booking (starts exactly when prior ends) is ALLOWED.
// This test FAILS due to the known boundary bug in hasConflict.
func TestCreateReservation_BoundaryConflict(t *testing.T) {
	base := time.Date(2025, 6, 2, 10, 0, 0, 0, time.UTC)
	adjacentStart := base.Add(30 * time.Minute)

	existing := &domain.Reservation{
		ID: "existing", DoctorID: "doc-001",
		StartsAt: base, EndsAt: base.Add(30 * time.Minute),
		Status: domain.ReservationStatus(domain.ReservationStatusConfirmed),
	}

	reservationRepo := mocks.NewReservationRepository(t)
	patientRepo := mocks.NewPatientRepository(t)

	patientRepo.EXPECT().GetPatientByPhone(context.Background(), "555-0002").Return(nil, nil)
	patientRepo.EXPECT().CreatePatient(context.Background(), mock.MatchedBy(func(_ *domain.Patient) bool { return true })).Return(nil).Maybe()

	// hasConflict: startsAt=10:30, endsAt=11:00 (30m bug). Window: [10:30-24h, 11:00+24h]
	newEnd := adjacentStart.Add(30 * time.Minute) // SlotDuration bug: always 30m
	reservationRepo.EXPECT().
		ListReservations(context.Background(), "doc-001", adjacentStart.Add(-24*time.Hour), newEnd.Add(24*time.Hour)).
		Return([]*domain.Reservation{existing}, nil)
	reservationRepo.EXPECT().
		CreateReservation(context.Background(), mock.MatchedBy(func(_ *domain.Reservation) bool { return true })).
		Return(nil).Maybe()

	uc := usecase.NewReservationUsecase(reservationRepo, patientRepo)

	res, err := uc.Create(context.Background(), usecase.CreateReservationInput{
		DoctorID: "doc-001", StartsAt: adjacentStart,
		Type:         domain.ReservationTypeFollowUp,
		PatientPhone: "555-0002", PatientName: "Another Patient", PatientEmail: "another@example.com",
	})

	require.NoError(t, err, "adjacent booking should be allowed")
	assert.NotNil(t, res)
}

// TestCreateReservation_FirstVisitDuration checks that a first visit allocates 60 minutes.
func TestCreateReservation_FirstVisitDuration(t *testing.T) {
	base := time.Date(2025, 6, 2, 9, 0, 0, 0, time.UTC)

	reservationRepo := mocks.NewReservationRepository(t)
	patientRepo := mocks.NewPatientRepository(t)

	patientRepo.EXPECT().GetPatientByPhone(context.Background(), "555-0003").Return(nil, nil)
	patientRepo.EXPECT().CreatePatient(context.Background(), mock.MatchedBy(func(_ *domain.Patient) bool { return true })).Return(nil)

	// hasConflict: startsAt=9:00, endsAt=10:00 (FirstVisit is 60 minutes).
	// Window: [9:00-24h, 10:00+24h]
	correctEnd := base.Add(60 * time.Minute) // FirstVisit should be 60 minutes
	reservationRepo.EXPECT().
		ListReservations(context.Background(), "doc-001", base.Add(-24*time.Hour), correctEnd.Add(24*time.Hour)).
		Return(nil, nil)
	reservationRepo.EXPECT().
		CreateReservation(context.Background(), mock.MatchedBy(func(_ *domain.Reservation) bool { return true })).
		Return(nil)

	uc := usecase.NewReservationUsecase(reservationRepo, patientRepo)

	res, err := uc.Create(context.Background(), usecase.CreateReservationInput{
		DoctorID: "doc-001", StartsAt: base,
		Type:         domain.ReservationTypeFirstVisit,
		PatientPhone: "555-0003", PatientName: "First Timer", PatientEmail: "first@example.com",
	})
	require.NoError(t, err)

	expected := 60 * time.Minute
	actual := res.EndsAt.Sub(res.StartsAt)
	assert.Equal(t, expected, actual, "first visit should be 60 minutes")
}

// TestListReservations is a skeleton — implement me.
func TestListReservations(t *testing.T) {
	t.Skip("TODO: implement list reservations test")
}

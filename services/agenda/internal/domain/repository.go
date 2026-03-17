package domain

import (
	"context"
	"time"
)

type DoctorRepository interface {
	GetDoctor(ctx context.Context, id string) (*Doctor, error)
	ListDoctors(ctx context.Context) ([]*Doctor, error)
}

type PatientRepository interface {
	GetPatient(ctx context.Context, id string) (*Patient, error)
	GetPatientByPhone(ctx context.Context, phone string) (*Patient, error)
	CreatePatient(ctx context.Context, p *Patient) error
	ListPatients(ctx context.Context) ([]*Patient, error)
	UpdatePatient(ctx context.Context, p *Patient) error
	DeletePatient(ctx context.Context, id string) error
}

type ReservationRepository interface {
	CreateReservation(ctx context.Context, r *Reservation) error
	GetReservation(ctx context.Context, id string) (*Reservation, error)
	// ListReservations returns reservations overlapping [from, to].
	// If doctorID is not empty, filters by doctor_id.
	// If patientID is not empty, filters by patient_id.
	ListReservations(ctx context.Context, doctorID, patientID string, from, to time.Time) ([]*Reservation, error)
	UpdateReservation(ctx context.Context, r *Reservation) error
	CancelReservation(ctx context.Context, id string) error
}

type BlockedSlotRepository interface {
	CreateBlockedSlot(ctx context.Context, b *BlockedSlot) error
	GetBlockedSlot(ctx context.Context, id string) (*BlockedSlot, error)
	ListBlockedSlots(ctx context.Context, doctorID string, from, to time.Time) ([]*BlockedSlot, error)
	DeleteBlockedSlot(ctx context.Context, id string) error
}

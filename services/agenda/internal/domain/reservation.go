package domain

import "time"

type ReservationType int

const (
	ReservationTypeUnspecified ReservationType = iota
	ReservationTypeFirstVisit
	ReservationTypeFollowUp
)

// SlotDuration returns the duration for a reservation type.
// FirstVisit requires 60 minutes, FollowUp requires 30 minutes.
func (t ReservationType) SlotDuration() time.Duration {
	switch t {
	case ReservationTypeFirstVisit:
		return 60 * time.Minute
	case ReservationTypeFollowUp:
		return 30 * time.Minute
	default:
		return 30 * time.Minute
	}
}

type ReservationStatus int

const (
	ReservationStatusUnspecified ReservationStatus = iota
	ReservationStatusConfirmed
	ReservationStatusCancelled
)

type Reservation struct {
	ID        string
	DoctorID  string
	PatientID string
	StartsAt  time.Time
	EndsAt    time.Time
	Type      ReservationType
	Status    ReservationStatus
}

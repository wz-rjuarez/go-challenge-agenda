package sqlite

import (
	"context"
	"fmt"
	"time"

	"go-challenge-agenda/services/agenda/internal/domain"
	"go-challenge-agenda/services/agenda/internal/repository/models"

	"gorm.io/gorm"
)

type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) CreateReservation(ctx context.Context, res *domain.Reservation) error {
	m := models.ReservationToModel(res)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ReservationRepository) GetReservation(ctx context.Context, id string) (*domain.Reservation, error) {
	var m models.Reservation
	res := r.db.WithContext(ctx).First(&m, "id = ?", id)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("reservation not found: %s", id)
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return models.ReservationFromModel(&m), nil
}

// ListReservations returns reservations overlapping [from, to].
// If doctorID is not empty, filters by doctor_id.
// If patientID is not empty, filters by patient_id.
func (r *ReservationRepository) ListReservations(ctx context.Context, doctorID, patientID string, from, to time.Time) ([]*domain.Reservation, error) {
	var ms []models.Reservation
	query := r.db.WithContext(ctx).Where("starts_at < ? AND ends_at > ?", to.UTC(), from.UTC())

	if doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}
	if patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	err := query.Find(&ms).Error
	if err != nil {
		return nil, err
	}
	reservations := make([]*domain.Reservation, len(ms))
	for i, m := range ms {
		m := m
		reservations[i] = models.ReservationFromModel(&m)
	}
	return reservations, nil
}

func (r *ReservationRepository) UpdateReservation(ctx context.Context, res *domain.Reservation) error {
	m := models.ReservationToModel(res)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *ReservationRepository) CancelReservation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", int(domain.ReservationStatusCancelled)).Error
}

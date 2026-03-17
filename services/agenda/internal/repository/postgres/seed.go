package postgres

import (
	"fmt"
	"math/rand"
	"time"

	"go-challenge-agenda/services/agenda/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&models.Doctor{}).Count(&count)
	if count > 0 {
		return nil
	}

	doctors := []models.Doctor{
		{ID: "doc-001", Name: "Dr. Ana García", Specialty: "General Practice"},
		{ID: "doc-002", Name: "Dr. Luis Mendoza", Specialty: "Cardiology"},
		{ID: "doc-003", Name: "Dr. Sara Patel", Specialty: "Pediatrics"},
	}

	whs := []models.WorkingHours{
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-001", Weekday: 1, FromTime: "09:00", ToTime: "17:00"},
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-001", Weekday: 2, FromTime: "09:00", ToTime: "17:00"},
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-001", Weekday: 3, FromTime: "09:00", ToTime: "17:00"},
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-001", Weekday: 4, FromTime: "09:00", ToTime: "17:00"},
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-001", Weekday: 5, FromTime: "09:00", ToTime: "13:00"},
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-002", Weekday: 1, FromTime: "08:00", ToTime: "16:00"},
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-002", Weekday: 3, FromTime: "08:00", ToTime: "16:00"},
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-002", Weekday: 5, FromTime: "08:00", ToTime: "12:00"},
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-003", Weekday: 2, FromTime: "10:00", ToTime: "18:00"},
		{ID: fmt.Sprintf("wh-%s", uuid.NewString()[:8]), DoctorID: "doc-003", Weekday: 4, FromTime: "10:00", ToTime: "18:00"},
	}

	blockedSlots := seedBlockedSlots()

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&doctors).Error; err != nil {
			return err
		}
		if err := tx.Create(&whs).Error; err != nil {
			return err
		}
		return tx.Create(&blockedSlots).Error
	})
}

// seedBlockedSlots generates 1–2 random blocked slots per doctor per week
// for the next 3 months, always within each doctor's working hours.
// Uses a fixed seed so the output is deterministic across restarts.
func seedBlockedSlots() []models.BlockedSlot {
	rng := rand.New(rand.NewSource(42))

	// working days per doctor: map[doctorID][]weekday
	workdays := map[string][]int{
		"doc-001": {1, 2, 3, 4, 5},
		"doc-002": {1, 3, 5},
		"doc-003": {2, 4},
	}

	// earliest start hour per doctor (their shift start)
	shiftStart := map[string]int{
		"doc-001": 9,
		"doc-002": 8,
		"doc-003": 10,
	}

	// latest end hour (last possible slot end)
	shiftEnd := map[string]int{
		"doc-001": 17,
		"doc-002": 16,
		"doc-003": 18,
	}

	reasons := []string{
		"Administrative meeting",
		"Lunch break",
		"Conference call",
		"Training session",
		"Personal appointment",
		"Research time",
	}

	now := time.Now().UTC().Truncate(24 * time.Hour)
	end := now.AddDate(0, 3, 0)

	var slots []models.BlockedSlot

	for doctorID, days := range workdays {
		for week := now; week.Before(end); week = week.AddDate(0, 0, 7) {
			// 1 or 2 blocked slots this week
			count := 1 + rng.Intn(2)
			// pick `count` distinct weekdays from this doctor's workdays
			chosen := rng.Perm(len(days))[:count]

			for _, idx := range chosen {
				wd := days[idx]

				// find the date of this weekday in the current week
				day := week
				for int(day.Weekday()) != wd {
					day = day.AddDate(0, 0, 1)
				}
				if !day.Before(end) {
					continue
				}

				// random start hour within [shiftStart, shiftEnd-2]
				startHour := shiftStart[doctorID] + rng.Intn(shiftEnd[doctorID]-shiftStart[doctorID]-1)
				// random duration: 30, 60, or 90 minutes
				durationMins := []int{30, 60, 90}[rng.Intn(3)]

				startsAt := time.Date(day.Year(), day.Month(), day.Day(), startHour, 0, 0, 0, time.UTC)
				endsAt := startsAt.Add(time.Duration(durationMins) * time.Minute)

				slots = append(slots, models.BlockedSlot{
					ID:             fmt.Sprintf("bs-%s", uuid.NewString()[:8]),
					DoctorID:       doctorID,
					StartsAt:       startsAt,
					EndsAt:         endsAt,
					Reason:         reasons[rng.Intn(len(reasons))],
					RecurrenceType: 0,
				})
			}
		}
	}

	return slots
}

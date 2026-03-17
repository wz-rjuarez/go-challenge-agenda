package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOccurrences_NoRecurrence(t *testing.T) {
	// Non-recurring slot should return itself if it overlaps the range
	slot := &BlockedSlot{
		ID:             "slot-1",
		DoctorID:       "doc-001",
		StartsAt:       time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:         time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:         "Meeting",
		RecurrenceType: RecurrenceNone,
	}

	from := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	occurrences := slot.Occurrences(from, to)
	assert.Len(t, occurrences, 1)
	assert.Equal(t, slot.StartsAt, occurrences[0].StartsAt)
	assert.Equal(t, slot.EndsAt, occurrences[0].EndsAt)
}

func TestOccurrences_NoRecurrence_OutOfRange(t *testing.T) {
	// Non-recurring slot should return nil if it doesn't overlap the range
	slot := &BlockedSlot{
		ID:             "slot-1",
		DoctorID:       "doc-001",
		StartsAt:       time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:         time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:         "Meeting",
		RecurrenceType: RecurrenceNone,
	}

	from := time.Date(2024, 1, 14, 9, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 14, 12, 0, 0, 0, time.UTC)

	occurrences := slot.Occurrences(from, to)
	assert.Nil(t, occurrences)
}

func TestOccurrences_Daily(t *testing.T) {
	// Daily recurrence should expand
	until := time.Date(2024, 1, 17, 23, 59, 59, 0, time.UTC)
	slot := &BlockedSlot{
		ID:              "slot-1",
		DoctorID:        "doc-001",
		StartsAt:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:          "Daily meeting",
		RecurrenceType:  RecurrenceDaily,
		RecurrenceUntil: &until,
	}

	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC)

	occurrences := slot.Occurrences(from, to)

	// Should have 3 occurrences: Jan 15, 16, 17
	assert.Len(t, occurrences, 3)
	assert.Equal(t, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), occurrences[0].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC), occurrences[1].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 17, 10, 0, 0, 0, time.UTC), occurrences[2].StartsAt)
}

func TestOccurrences_Weekly(t *testing.T) {
	// Weekly recurrence on Monday
	until := time.Date(2024, 2, 5, 23, 59, 59, 0, time.UTC)
	slot := &BlockedSlot{
		ID:              "slot-1",
		DoctorID:        "doc-001",
		StartsAt:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), // Monday
		EndsAt:          time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:          "Weekly meeting",
		RecurrenceType:  RecurrenceWeekly,
		RecurrenceUntil: &until,
	}

	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)

	occurrences := slot.Occurrences(from, to)

	// Should have 4 occurrences: Jan 15, 22, 29, Feb 5
	assert.Len(t, occurrences, 4)
	assert.Equal(t, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), occurrences[0].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 22, 10, 0, 0, 0, time.UTC), occurrences[1].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 29, 10, 0, 0, 0, time.UTC), occurrences[2].StartsAt)
	assert.Equal(t, time.Date(2024, 2, 5, 10, 0, 0, 0, time.UTC), occurrences[3].StartsAt)
}

func TestOccurrences_Monthly(t *testing.T) {
	// Monthly recurrence on the 15th
	until := time.Date(2024, 4, 15, 23, 59, 59, 0, time.UTC)
	slot := &BlockedSlot{
		ID:              "slot-1",
		DoctorID:        "doc-001",
		StartsAt:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:          "Monthly review",
		RecurrenceType:  RecurrenceMonthly,
		RecurrenceUntil: &until,
	}

	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)

	occurrences := slot.Occurrences(from, to)

	// Should have 4 occurrences: Jan 15, Feb 15, Mar 15, Apr 15
	assert.Len(t, occurrences, 4)
	assert.Equal(t, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), occurrences[0].StartsAt)
	assert.Equal(t, time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC), occurrences[1].StartsAt)
	assert.Equal(t, time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC), occurrences[2].StartsAt)
	assert.Equal(t, time.Date(2024, 4, 15, 10, 0, 0, 0, time.UTC), occurrences[3].StartsAt)
}

func TestOccurrences_Monthly_EdgeCase_EndOfMonth(t *testing.T) {
	// Monthly recurrence on the 31st, which doesn't exist in all months
	until := time.Date(2024, 4, 30, 23, 59, 59, 0, time.UTC)
	slot := &BlockedSlot{
		ID:              "slot-1",
		DoctorID:        "doc-001",
		StartsAt:        time.Date(2024, 1, 31, 10, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 31, 11, 0, 0, 0, time.UTC),
		Reason:          "Monthly review",
		RecurrenceType:  RecurrenceMonthly,
		RecurrenceUntil: &until,
	}

	from := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)

	occurrences := slot.Occurrences(from, to)

	// Should have 4 occurrences with proper handling of Feb 29 (leap year)
	// Jan 31, Feb 29 (leap year adjusts to last day), Mar 31, Apr 30 (adjusts to last day)
	assert.Len(t, occurrences, 4)
	assert.Equal(t, time.Date(2024, 1, 31, 10, 0, 0, 0, time.UTC), occurrences[0].StartsAt)
	assert.Equal(t, time.Date(2024, 2, 29, 10, 0, 0, 0, time.UTC), occurrences[1].StartsAt) // Feb 2024 is leap year
	assert.Equal(t, time.Date(2024, 3, 31, 10, 0, 0, 0, time.UTC), occurrences[2].StartsAt)
	assert.Equal(t, time.Date(2024, 4, 30, 10, 0, 0, 0, time.UTC), occurrences[3].StartsAt) // April only has 30 days
}

func TestOccurrences_RespectRecurrenceUntil(t *testing.T) {
	// Recurrence should stop at RecurrenceUntil even if the range extends further
	until := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	slot := &BlockedSlot{
		ID:              "slot-1",
		DoctorID:        "doc-001",
		StartsAt:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:          "Limited daily",
		RecurrenceType:  RecurrenceDaily,
		RecurrenceUntil: &until,
	}

	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC) // Much longer range

	occurrences := slot.Occurrences(from, to)

	// Should only have occurrences until Jan 31
	// Jan 15-31 = 17 days
	assert.Len(t, occurrences, 17)
	assert.Equal(t, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), occurrences[0].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 31, 10, 0, 0, 0, time.UTC), occurrences[16].StartsAt)
}

func TestOccurrences_FilterByRange(t *testing.T) {
	// Should filter occurrences by the [from, to] range
	until := time.Date(2024, 2, 28, 23, 59, 59, 0, time.UTC)
	slot := &BlockedSlot{
		ID:              "slot-1",
		DoctorID:        "doc-001",
		StartsAt:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:          "Daily meeting",
		RecurrenceType:  RecurrenceDaily,
		RecurrenceUntil: &until,
	}

	from := time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 25, 23, 59, 59, 0, time.UTC)

	occurrences := slot.Occurrences(from, to)

	// Should have 6 occurrences: Jan 20, 21, 22, 23, 24, 25
	assert.Len(t, occurrences, 6)
	assert.Equal(t, time.Date(2024, 1, 20, 10, 0, 0, 0, time.UTC), occurrences[0].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 25, 10, 0, 0, 0, time.UTC), occurrences[5].StartsAt)
}

func TestOccurrences_PreservesMetadata(t *testing.T) {
	// All occurrences should preserve the original slot's metadata
	until := time.Date(2024, 1, 17, 23, 59, 59, 0, time.UTC)
	slot := &BlockedSlot{
		ID:              "slot-original",
		DoctorID:        "doc-special",
		StartsAt:        time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 15, 15, 30, 0, 0, time.UTC),
		Reason:          "Important meeting",
		RecurrenceType:  RecurrenceDaily,
		RecurrenceUntil: &until,
	}

	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC)

	occurrences := slot.Occurrences(from, to)

	for _, occ := range occurrences {
		assert.Equal(t, "slot-original", occ.ID)
		assert.Equal(t, "doc-special", occ.DoctorID)
		assert.Equal(t, "Important meeting", occ.Reason)
		assert.Equal(t, RecurrenceDaily, occ.RecurrenceType)
		assert.Equal(t, &until, occ.RecurrenceUntil)
		// Duration should be preserved (1 hour)
		assert.Equal(t, time.Hour, occ.EndsAt.Sub(occ.StartsAt))
	}
}

func TestOccurrences_NoRecurrenceUntil(t *testing.T) {
	// When RecurrenceUntil is nil, should still respect the query range
	slot := &BlockedSlot{
		ID:              "slot-1",
		DoctorID:        "doc-001",
		StartsAt:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:          "Ongoing daily",
		RecurrenceType:  RecurrenceDaily,
		RecurrenceUntil: nil, // No end date
	}

	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 20, 23, 59, 59, 0, time.UTC)

	occurrences := slot.Occurrences(from, to)

	// Should have 6 occurrences: Jan 15-20
	assert.Len(t, occurrences, 6)
	assert.Equal(t, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), occurrences[0].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 20, 10, 0, 0, 0, time.UTC), occurrences[5].StartsAt)
}

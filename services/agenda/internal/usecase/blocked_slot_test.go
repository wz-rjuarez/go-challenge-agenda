package usecase_test

import (
	"context"
	"testing"
	"time"

	"go-challenge-agenda/services/agenda/internal/domain"
	"go-challenge-agenda/services/agenda/internal/repository/sqlite"
	"go-challenge-agenda/services/agenda/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockedSlotUsecase_List_ExpandsRecurrences(t *testing.T) {
	// Setup in-memory database
	db, err := sqlite.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, sqlite.Migrate(db))

	repo := sqlite.NewBlockedSlotRepository(db)
	uc := usecase.NewBlockedSlotUsecase(repo)
	ctx := context.Background()

	// Create a daily recurring blocked slot
	until := time.Date(2024, 1, 20, 23, 59, 59, 0, time.UTC)
	slot := &domain.BlockedSlot{
		ID:              "test-slot-1",
		DoctorID:        "doc-001",
		StartsAt:        time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Reason:          "Daily meeting",
		RecurrenceType:  domain.RecurrenceDaily,
		RecurrenceUntil: &until,
	}

	err = repo.CreateBlockedSlot(ctx, slot)
	require.NoError(t, err)

	// List blocked slots for a date range
	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 20, 23, 59, 59, 0, time.UTC)

	result, err := uc.List(ctx, "doc-001", from, to)
	require.NoError(t, err)

	// Should have 6 occurrences (Jan 15-20)
	assert.Len(t, result, 6)

	// Verify dates
	assert.Equal(t, time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC), result[0].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 16, 9, 0, 0, 0, time.UTC), result[1].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 17, 9, 0, 0, 0, time.UTC), result[2].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 18, 9, 0, 0, 0, time.UTC), result[3].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 19, 9, 0, 0, 0, time.UTC), result[4].StartsAt)
	assert.Equal(t, time.Date(2024, 1, 20, 9, 0, 0, 0, time.UTC), result[5].StartsAt)

	// All should preserve metadata
	for _, occ := range result {
		assert.Equal(t, "doc-001", occ.DoctorID)
		assert.Equal(t, "Daily meeting", occ.Reason)
		assert.Equal(t, domain.RecurrenceDaily, occ.RecurrenceType)
	}
}

func TestBlockedSlotUsecase_List_WeeklyRecurrence(t *testing.T) {
	// Setup in-memory database
	db, err := sqlite.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, sqlite.Migrate(db))

	repo := sqlite.NewBlockedSlotRepository(db)
	uc := usecase.NewBlockedSlotUsecase(repo)
	ctx := context.Background()

	// Create a weekly recurring blocked slot (every Monday)
	until := time.Date(2024, 2, 12, 23, 59, 59, 0, time.UTC)
	slot := &domain.BlockedSlot{
		ID:              "test-slot-2",
		DoctorID:        "doc-002",
		StartsAt:        time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC), // Monday
		EndsAt:          time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC),
		Reason:          "Weekly team meeting",
		RecurrenceType:  domain.RecurrenceWeekly,
		RecurrenceUntil: &until,
	}

	err = repo.CreateBlockedSlot(ctx, slot)
	require.NoError(t, err)

	// List blocked slots for a date range
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 2, 29, 23, 59, 59, 0, time.UTC)

	result, err := uc.List(ctx, "doc-002", from, to)
	require.NoError(t, err)

	// Should have 5 occurrences (Jan 15, 22, 29, Feb 5, 12)
	assert.Len(t, result, 5)

	// Verify it's every Monday
	for _, occ := range result {
		assert.Equal(t, time.Monday, occ.StartsAt.Weekday(), "Should be Monday: %v", occ.StartsAt)
	}
}

func TestBlockedSlotUsecase_List_MonthlyRecurrence(t *testing.T) {
	// Setup in-memory database
	db, err := sqlite.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, sqlite.Migrate(db))

	repo := sqlite.NewBlockedSlotRepository(db)
	uc := usecase.NewBlockedSlotUsecase(repo)
	ctx := context.Background()

	// Create a monthly recurring blocked slot on the 15th
	until := time.Date(2024, 5, 15, 23, 59, 59, 0, time.UTC)
	slot := &domain.BlockedSlot{
		ID:              "test-slot-3",
		DoctorID:        "doc-003",
		StartsAt:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:          "Monthly review",
		RecurrenceType:  domain.RecurrenceMonthly,
		RecurrenceUntil: &until,
	}

	err = repo.CreateBlockedSlot(ctx, slot)
	require.NoError(t, err)

	// List blocked slots for a date range
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 1, 23, 59, 59, 0, time.UTC)

	result, err := uc.List(ctx, "doc-003", from, to)
	require.NoError(t, err)

	// Should have 5 occurrences (Jan 15, Feb 15, Mar 15, Apr 15, May 15)
	assert.Len(t, result, 5)

	// Verify it's the 15th of each month
	for _, occ := range result {
		assert.Equal(t, 15, occ.StartsAt.Day(), "Should be 15th of month: %v", occ.StartsAt)
	}
}

func TestBlockedSlotUsecase_List_NoRecurrence(t *testing.T) {
	// Setup in-memory database
	db, err := sqlite.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, sqlite.Migrate(db))

	repo := sqlite.NewBlockedSlotRepository(db)
	uc := usecase.NewBlockedSlotUsecase(repo)
	ctx := context.Background()

	// Create a non-recurring blocked slot
	slot := &domain.BlockedSlot{
		ID:             "test-slot-4",
		DoctorID:       "doc-001",
		StartsAt:       time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndsAt:         time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Reason:         "One-time block",
		RecurrenceType: domain.RecurrenceNone,
	}

	err = repo.CreateBlockedSlot(ctx, slot)
	require.NoError(t, err)

	// List blocked slots for a date range
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	result, err := uc.List(ctx, "doc-001", from, to)
	require.NoError(t, err)

	// Should have only 1 occurrence
	assert.Len(t, result, 1)
	assert.Equal(t, time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), result[0].StartsAt)
}

func TestBlockedSlotUsecase_List_MultipleSlots(t *testing.T) {
	// Setup in-memory database
	db, err := sqlite.Open(":memory:")
	require.NoError(t, err)
	require.NoError(t, sqlite.Migrate(db))

	repo := sqlite.NewBlockedSlotRepository(db)
	uc := usecase.NewBlockedSlotUsecase(repo)
	ctx := context.Background()

	// Create multiple blocked slots with different recurrence patterns
	until1 := time.Date(2024, 1, 20, 23, 59, 59, 0, time.UTC)
	slot1 := &domain.BlockedSlot{
		ID:              "test-slot-5",
		DoctorID:        "doc-001",
		StartsAt:        time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
		EndsAt:          time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Reason:          "Daily meeting",
		RecurrenceType:  domain.RecurrenceDaily,
		RecurrenceUntil: &until1,
	}

	slot2 := &domain.BlockedSlot{
		ID:             "test-slot-6",
		DoctorID:       "doc-001",
		StartsAt:       time.Date(2024, 1, 18, 14, 0, 0, 0, time.UTC),
		EndsAt:         time.Date(2024, 1, 18, 15, 0, 0, 0, time.UTC),
		Reason:         "One-time event",
		RecurrenceType: domain.RecurrenceNone,
	}

	err = repo.CreateBlockedSlot(ctx, slot1)
	require.NoError(t, err)
	err = repo.CreateBlockedSlot(ctx, slot2)
	require.NoError(t, err)

	// List blocked slots for a date range
	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 20, 23, 59, 59, 0, time.UTC)

	result, err := uc.List(ctx, "doc-001", from, to)
	require.NoError(t, err)

	// Should have 7 occurrences (6 from daily + 1 one-time)
	assert.Len(t, result, 7)
}

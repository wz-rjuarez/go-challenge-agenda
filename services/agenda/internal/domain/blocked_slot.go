package domain

import "time"

type RecurrenceType int

const (
	RecurrenceNone RecurrenceType = iota
	RecurrenceDaily
	RecurrenceWeekly
	RecurrenceMonthly
)

type BlockedSlot struct {
	ID              string
	DoctorID        string
	StartsAt        time.Time
	EndsAt          time.Time
	Reason          string
	RecurrenceType  RecurrenceType
	RecurrenceUntil *time.Time
}

// Occurrences returns all occurrences of this blocked slot within [from, to].
// For recurring slots, it expands the recurrence pattern (daily/weekly/monthly)
// until RecurrenceUntil or the end of the range, whichever comes first.
func (b *BlockedSlot) Occurrences(from, to time.Time) []BlockedSlot {
	// Non-recurring slot
	if b.RecurrenceType == RecurrenceNone {
		if b.StartsAt.After(to) || b.EndsAt.Before(from) {
			return nil
		}
		return []BlockedSlot{*b}
	}

	// Recurring slot - expand occurrences
	var occurrences []BlockedSlot

	// Determine the end boundary for expansion
	// If RecurrenceUntil is set, we need to include all occurrences on that day
	endBoundary := to
	if b.RecurrenceUntil != nil {
		// Add 23:59:59 to include the entire day
		recurrenceEndOfDay := time.Date(
			b.RecurrenceUntil.Year(), b.RecurrenceUntil.Month(), b.RecurrenceUntil.Day(),
			23, 59, 59, 999999999, b.RecurrenceUntil.Location(),
		)
		if recurrenceEndOfDay.Before(to) {
			endBoundary = recurrenceEndOfDay
		}
	}

	// Start from the base slot time
	current := b.StartsAt
	duration := b.EndsAt.Sub(b.StartsAt)

	// For monthly recurrence, remember the original day of month
	originalDay := b.StartsAt.Day()

	for {
		// Calculate end time for this occurrence
		occEnd := current.Add(duration)

		// Check if current occurrence starts beyond the end boundary
		if current.After(endBoundary) {
			break
		}

		// If this occurrence overlaps with [from, to], include it
		if !current.After(to) && !occEnd.Before(from) {
			occurrences = append(occurrences, BlockedSlot{
				ID:              b.ID,
				DoctorID:        b.DoctorID,
				StartsAt:        current,
				EndsAt:          occEnd,
				Reason:          b.Reason,
				RecurrenceType:  b.RecurrenceType,
				RecurrenceUntil: b.RecurrenceUntil,
			})
		}

		// Move to next occurrence based on recurrence type
		switch b.RecurrenceType {
		case RecurrenceDaily:
			current = current.AddDate(0, 0, 1)
		case RecurrenceWeekly:
			current = current.AddDate(0, 0, 7)
		case RecurrenceMonthly:
			// Add one month, then adjust to the correct day
			year, month, _ := current.Date()
			nextMonth := month + 1
			nextYear := year
			if nextMonth > 12 {
				nextMonth = 1
				nextYear++
			}

			// Try to use the original day, but adjust if the month doesn't have enough days
			nextDay := originalDay
			// Get the last day of the next month
			lastDayOfNextMonth := time.Date(nextYear, nextMonth+1, 0, 0, 0, 0, 0, current.Location()).Day()
			if nextDay > lastDayOfNextMonth {
				nextDay = lastDayOfNextMonth
			}

			current = time.Date(nextYear, nextMonth, nextDay,
				current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), current.Location())
		default:
			// Should not reach here, but break to avoid infinite loop
			break
		}

		// Safety check: if we've moved beyond our end boundary, stop
		if current.After(endBoundary) {
			break
		}
	}

	return occurrences
}

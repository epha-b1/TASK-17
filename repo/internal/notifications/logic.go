package notifications

import "time"

func InDNDWindow(now time.Time, start, end time.Time) bool {
	n := now.UTC()
	startAt := time.Date(n.Year(), n.Month(), n.Day(), start.Hour(), start.Minute(), start.Second(), 0, time.UTC)
	endAt := time.Date(n.Year(), n.Month(), n.Day(), end.Hour(), end.Minute(), end.Second(), 0, time.UTC)
	if !endAt.After(startAt) {
		if n.Before(endAt) {
			startAt = startAt.Add(-24 * time.Hour)
		} else {
			endAt = endAt.Add(24 * time.Hour)
		}
	}
	return (n.Equal(startAt) || n.After(startAt)) && n.Before(endAt)
}

func DNDEnd(now time.Time, start, end time.Time) time.Time {
	n := now.UTC()
	endAt := time.Date(n.Year(), n.Month(), n.Day(), end.Hour(), end.Minute(), end.Second(), 0, time.UTC)
	if endAt.Before(n) || endAt.Equal(n) {
		endAt = endAt.Add(24 * time.Hour)
	}
	if !time.Date(n.Year(), n.Month(), n.Day(), end.Hour(), end.Minute(), end.Second(), 0, time.UTC).After(
		time.Date(n.Year(), n.Month(), n.Day(), start.Hour(), start.Minute(), start.Second(), 0, time.UTC),
	) {
		if n.Before(endAt.Add(-24 * time.Hour)) {
			endAt = endAt.Add(-24 * time.Hour)
		}
	}
	return endAt
}

func AllowByFrequencyCap(existingToday int) bool {
	return existingToday < 3
}

func RetryBackoff(attemptCount int) (time.Duration, bool) {
	if attemptCount >= 5 {
		return 0, false
	}
	mins := 1 << attemptCount
	return time.Duration(mins) * time.Minute, true
}

package unit_tests

import (
	"testing"
	"time"

	"parkops/internal/notifications"
)

func TestDNDDefer(t *testing.T) {
	now := time.Date(2026, 4, 1, 23, 30, 0, 0, time.UTC)
	start := time.Date(2026, 4, 1, 22, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 1, 6, 0, 0, 0, time.UTC)

	if !notifications.InDNDWindow(now, start, end) {
		t.Fatalf("expected now to be inside DND window")
	}
	deferTo := notifications.DNDEnd(now, start, end)
	if deferTo.Hour() != 6 || deferTo.Minute() != 0 {
		t.Fatalf("expected DND defer to 06:00, got %s", deferTo.Format(time.RFC3339))
	}
}

func TestFrequencyCap(t *testing.T) {
	if !notifications.AllowByFrequencyCap(0) || !notifications.AllowByFrequencyCap(2) {
		t.Fatalf("expected first three reminders allowed")
	}
	if notifications.AllowByFrequencyCap(3) {
		t.Fatalf("expected fourth reminder to be rejected")
	}
}

func TestRetryBackoff(t *testing.T) {
	d, ok := notifications.RetryBackoff(0)
	if !ok || d != time.Minute {
		t.Fatalf("attempt0 expected 1m got %v ok=%v", d, ok)
	}
	d, ok = notifications.RetryBackoff(3)
	if !ok || d != 8*time.Minute {
		t.Fatalf("attempt3 expected 8m got %v ok=%v", d, ok)
	}
	_, ok = notifications.RetryBackoff(5)
	if ok {
		t.Fatalf("attempt5 expected no retry")
	}
}

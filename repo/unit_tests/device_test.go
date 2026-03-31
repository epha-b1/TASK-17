package unit_tests

import (
	"testing"
	"time"

	"parkops/internal/devices"
)

func TestDeviceDeduplicationRule(t *testing.T) {
	replayed, skipped := devices.ReplayDecision(1)
	if replayed || !skipped {
		t.Fatalf("expected already replayed event to be skipped")
	}
}

func TestOutOfOrderWithinWindow(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	late, reordered := devices.ClassifySequence(100, now.Add(-5*time.Minute), 99, now, 10*time.Minute)
	if late || !reordered {
		t.Fatalf("expected reordered within window, got late=%v reordered=%v", late, reordered)
	}
}

func TestLateFlagOutsideWindow(t *testing.T) {
	now := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	late, reordered := devices.ClassifySequence(100, now.Add(-11*time.Minute), 90, now, 10*time.Minute)
	if !late || reordered {
		t.Fatalf("expected late flag outside window, got late=%v reordered=%v", late, reordered)
	}
}

func TestReplayDecision(t *testing.T) {
	replayed, skipped := devices.ReplayDecision(0)
	if !replayed || skipped {
		t.Fatalf("expected first replay request to replay event")
	}

	replayed, skipped = devices.ReplayDecision(2)
	if replayed || !skipped {
		t.Fatalf("expected repeat replay request to skip")
	}
}

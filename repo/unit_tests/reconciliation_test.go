package unit_tests

import (
	"testing"

	"parkops/internal/reconciliation"
)

func TestReconciliationPositiveDeltaGeneratesCompensatingHold(t *testing.T) {
	eventType, delta, needed := reconciliation.DecideCompensatingEvent(5, 7)
	if !needed || eventType != reconciliation.CompensatingHold || delta != 2 {
		t.Fatalf("expected hold delta=2, got needed=%v type=%q delta=%d", needed, eventType, delta)
	}
}

func TestReconciliationNegativeDeltaGeneratesCompensatingRelease(t *testing.T) {
	eventType, delta, needed := reconciliation.DecideCompensatingEvent(7, 5)
	if !needed || eventType != reconciliation.CompensatingRelease || delta != 2 {
		t.Fatalf("expected release delta=2, got needed=%v type=%q delta=%d", needed, eventType, delta)
	}
}

func TestReconciliationZeroDeltaGeneratesNothing(t *testing.T) {
	eventType, delta, needed := reconciliation.DecideCompensatingEvent(6, 6)
	if needed || eventType != "" || delta != 0 {
		t.Fatalf("expected no compensating event, got needed=%v type=%q delta=%d", needed, eventType, delta)
	}
}

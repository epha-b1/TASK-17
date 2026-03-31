package unit_tests

import (
	"testing"

	"parkops/internal/exceptions"
)

func TestExceptionCreatedFromDeviceEvent(t *testing.T) {
	exType, ok := exceptions.ExceptionTypeForEvent("camera_error")
	if !ok || exType != "camera_error" {
		t.Fatalf("expected camera_error to create exception, got ok=%v type=%q", ok, exType)
	}

	if _, ok := exceptions.ExceptionTypeForEvent("camera_ping"); ok {
		t.Fatalf("expected non-exception event type to be ignored")
	}
}

func TestAcknowledgeTransitionsStatus(t *testing.T) {
	next, transitioned := exceptions.AcknowledgeTransition("open")
	if !transitioned || next != "acknowledged" {
		t.Fatalf("expected open to transition to acknowledged, got transitioned=%v next=%q", transitioned, next)
	}

	next, transitioned = exceptions.AcknowledgeTransition("acknowledged")
	if transitioned || next != "acknowledged" {
		t.Fatalf("expected acknowledged status to remain unchanged, got transitioned=%v next=%q", transitioned, next)
	}
}

package unit_tests

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"parkops/internal/tracking"
)

func TestDriftThresholdSuspectOver500MetersIn30Seconds(t *testing.T) {
	last := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	incoming := last.Add(30 * time.Second)
	if !tracking.IsSuspectJump(last, incoming, 37.7749, -122.4194, 37.7810, -122.4194) {
		t.Fatalf("expected jump over 500m in 30s to be suspect")
	}
}

func TestStopDetectionAtExactly3Minutes(t *testing.T) {
	if !tracking.ShouldCreateStop(3*time.Minute, 50) {
		t.Fatalf("expected stop creation at exactly 3 minutes and 50 meters")
	}
}

func TestTrustedTimestampHMACValidation(t *testing.T) {
	deviceTime := "2026-01-01T10:00:00Z"
	secret := "RES-123"
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(deviceTime))
	sig := hex.EncodeToString(mac.Sum(nil))

	if !tracking.ValidateDeviceTimeHMAC(deviceTime, sig, secret) {
		t.Fatalf("expected signature validation to pass")
	}
	if tracking.ValidateDeviceTimeHMAC(deviceTime, sig, "wrong-secret") {
		t.Fatalf("expected signature validation to fail for wrong secret")
	}
}

func TestSuspectDiscardedWithoutConfirmation(t *testing.T) {
	if tracking.ConfirmsSuspect(37.7810, -122.4194, 37.7900, -122.4194) {
		t.Fatalf("expected distant next report to not confirm suspect")
	}
}

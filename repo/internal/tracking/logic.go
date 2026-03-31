package tracking

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strings"
	"time"
)

const (
	EarthRadiusMeters   = 6371000.0
	DriftMeters         = 500.0
	DriftWindow         = 30 * time.Second
	StopDistanceMeters  = 50.0
	StopDurationMinimum = 3 * time.Minute
)

func DistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := radians(lat2 - lat1)
	dLon := radians(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(radians(lat1))*math.Cos(radians(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadiusMeters * c
}

func IsSuspectJump(lastAt, incomingAt time.Time, lastLat, lastLon, incomingLat, incomingLon float64) bool {
	if incomingAt.Before(lastAt) {
		return false
	}
	if incomingAt.Sub(lastAt) > DriftWindow {
		return false
	}
	return DistanceMeters(lastLat, lastLon, incomingLat, incomingLon) > DriftMeters
}

func ConfirmsSuspect(suspectLat, suspectLon, incomingLat, incomingLon float64) bool {
	return DistanceMeters(suspectLat, suspectLon, incomingLat, incomingLon) <= DriftMeters
}

func ShouldCreateStop(stationaryDuration time.Duration, movedMeters float64) bool {
	return stationaryDuration >= StopDurationMinimum && movedMeters <= StopDistanceMeters
}

func ValidateDeviceTimeHMAC(deviceTime, signatureHex, secret string) bool {
	if strings.TrimSpace(deviceTime) == "" || strings.TrimSpace(signatureHex) == "" || strings.TrimSpace(secret) == "" {
		return false
	}
	provided, err := hex.DecodeString(strings.TrimSpace(signatureHex))
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(deviceTime))
	expected := mac.Sum(nil)
	return hmac.Equal(provided, expected)
}

func radians(v float64) float64 {
	return v * math.Pi / 180
}

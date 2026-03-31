package exceptions

import "strings"

func ExceptionTypeForEvent(eventType string) (string, bool) {
	normalized := strings.TrimSpace(strings.ToLower(eventType))
	switch normalized {
	case "gate_stuck", "sensor_offline", "camera_error":
		return normalized, true
	default:
		return "", false
	}
}

func AcknowledgeTransition(status string) (next string, transitioned bool) {
	if strings.EqualFold(strings.TrimSpace(status), "open") {
		return "acknowledged", true
	}
	return strings.TrimSpace(strings.ToLower(status)), false
}

package devices

import "time"

func ClassifySequence(lastSequence int64, lastSeenAt time.Time, incomingSequence int64, now time.Time, window time.Duration) (late bool, reordered bool) {
	if incomingSequence >= lastSequence || lastSequence == 0 {
		return false, false
	}
	if lastSeenAt.IsZero() {
		return true, false
	}
	if now.Sub(lastSeenAt) > window {
		return true, false
	}
	return false, true
}

func ReplayDecision(replayCount int) (replayed bool, skipped bool) {
	if replayCount > 0 {
		return false, true
	}
	return true, false
}

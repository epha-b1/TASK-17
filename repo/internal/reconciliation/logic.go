package reconciliation

const (
	CompensatingHold    = "hold"
	CompensatingRelease = "release"
)

func DecideCompensatingEvent(snapshotStalls, eventDerivedStalls int) (eventType string, stallDelta int, needed bool) {
	if eventDerivedStalls > snapshotStalls {
		return CompensatingHold, eventDerivedStalls - snapshotStalls, true
	}
	if eventDerivedStalls < snapshotStalls {
		return CompensatingRelease, snapshotStalls - eventDerivedStalls, true
	}
	return "", 0, false
}

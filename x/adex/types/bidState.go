package types

const (
	BidStateUnknown = 0
	BidStateActive = 1
	// fail states
	BidStateCancelled = 2
	BidStateTimedOut = 3
	BidStateFailed = 4

	// success states
	BidStateSucceeded = 5
)

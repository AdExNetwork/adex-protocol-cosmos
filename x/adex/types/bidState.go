package types

const (
	BidStateUnknown = 0
	BidStateActive = 1
	// fail states
	BidStateCanceled = 2
	BidStateExpired = 3
	BidStateFailed = 4

	// success states
	BidStateSucceeded = 5
)

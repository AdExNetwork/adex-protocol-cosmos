package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The commitment refund is stored in our validUntil store, in order to refund expired commitments
type CommitmentRefund struct {
	BidId BidId `json:"bidId"`
	Beneficiary sdk.AccAddress `json:"beneficiary"`
	TotalReward sdk.Coins `json:"reward"`
}

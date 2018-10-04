package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Commitment struct {
	BidId []byte `json:"bidId"`
	Reward sdk.Coins `json:"reward"`
	ValidUntil int64 `json:"validUntil"`
	Advertiser sdk.AccAddress `json:"advertiser"`
	Publisher sdk.AccAddress `json:"publisher"`
	Validators []Validator `json:"validators"`
}
// @TODO: validation, including all Validator rewards is not > Reward
// alternatively, this may not be a requirement, but we have to ensure all those amounts are escrowed on commitmentStart

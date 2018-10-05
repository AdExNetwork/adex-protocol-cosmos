package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MIN_VALIDATOR_COUNT = 2
)

type Commitment struct {
	BidId []byte `json:"bidId"`
	TotalReward sdk.Coins `json:"totalReward"`
	ValidUntil int64 `json:"validUntil"`
	Advertiser sdk.AccAddress `json:"advertiser"`
	Publisher sdk.AccAddress `json:"publisher"`
	Validators []Validator `json:"validators"`
}

func (commitment Commitment) IsValid() bool {
	if len(commitment.Validators) < MIN_VALIDATOR_COUNT {
		return false
	}

	return true
}

// @TODO: GetHash()
// @TODO: IsValid() : should check if the validator reward IsLT the sum of all validator rewards (same as on eth)
// @TODO: FromBid() : last arg would be extraValidatorAddr; test if != nil, also test if .IsValid or smth

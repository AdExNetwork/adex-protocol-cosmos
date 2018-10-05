package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	minValidatorCount = 2
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
	if commitment.Validators == nil {
		return false
	}
	if len(commitment.Validators) < minValidatorCount {
		return false
	}

	if commitment.ValidUntil == 0 {
		return false
	}

	if commitment.TotalReward == nil {
		return false
	}
	// @TODO: should we have those checks for the addresses?

	validatorRewards := sdk.Coins{}
	for _, validator := range commitment.Validators {
		validatorRewards = validatorRewards.Plus(validator.Reward)
	}
	if commitment.TotalReward.IsLT(validatorRewards) {
		return false
	}

	return true
}

// @TODO: GetHash()
// @TODO: IsValid() : should check if the validator reward IsLT the sum of all validator rewards (same as on eth)
// @TODO: FromBid() : last arg would be extraValidatorAddr; test if != nil, also test if .IsValid or smth

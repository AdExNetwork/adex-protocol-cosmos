package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/crypto/sha3"
	"encoding/json"
)

const (
	minValidatorCount = 2
)

type Commitment struct {
	BidId [32]byte `json:"bidId"`
	TotalReward sdk.Coins `json:"totalReward"`
	ValidUntil int64 `json:"validUntil"`
	Advertiser sdk.AccAddress `json:"advertiser"`
	Publisher sdk.AccAddress `json:"publisher"`
	Validators []Validator `json:"validators"`
}

// @TODO; tests
func (commitment Commitment) IsValid() bool {
	if commitment.Validators == nil {
		return false
	}
	if len(commitment.Validators) < minValidatorCount {
		return false
	}

	if commitment.ValidUntil <= 0 {
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

func (commitment Commitment) Hash() [32]byte {
	b, err := json.Marshal(commitment)
	if err != nil {
		panic(err)
	}
	return sha3.Sum256(b)
}

func NewCommitmentFromBid(bid Bid, publisher sdk.AccAddress, validUntil int64, extraValidator sdk.AccAddress) Commitment {
	validators := bid.Validators
	if extraValidator != nil && !extraValidator.Empty() {
		validators = append(validators, Validator{
			Address: extraValidator,
			// The extra validator should not be allowed to set their own reward
			Reward: sdk.Coins{},
		})
	}
	return Commitment{
		BidId: bid.Hash(),
		TotalReward: bid.TotalReward,
		ValidUntil: validUntil,
		Publisher: publisher,
		Advertiser: bid.Advertiser,
		Validators: validators,
	}
}

// @TODO: FromBid() : last arg would be extraValidatorAddr; test if != nil, also test if .IsValid or smth

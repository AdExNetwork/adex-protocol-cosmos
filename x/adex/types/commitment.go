package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/crypto/sha3"
	"github.com/tendermint/tendermint/crypto"
	"encoding/json"
)

const (
	minValidatorCount = 2
)

type CommitmentId [32]byte

type Commitment struct {
	BidId BidId `json:"bidId"`
	TotalReward sdk.Coins `json:"totalReward"`
	ValidUntil int64 `json:"validUntil"`
	Advertiser sdk.AccAddress `json:"advertiser"`
	Publisher sdk.AccAddress `json:"publisher"`
	Validators []Validator `json:"validators"`
}

// @TODO; tests
func (commitment Commitment) IsValid() bool {
	// @TODO: figure out if we need the nil ref/slice checks; it all depends on what happens when deserializing
	if commitment.Validators == nil || commitment.TotalReward == nil {
		return false
	}
	if len(commitment.Validators) < minValidatorCount {
		return false
	}
	for _, validator := range commitment.Validators {
		if !validator.IsValid() {
			return false
		}
	}

	if commitment.ValidUntil <= 0 {
		return false
	}

	if !commitment.TotalReward.IsNotNegative() {
		return false
	}

	validatorRewards := sdk.Coins{}
	for _, validator := range commitment.Validators {
		validatorRewards = validatorRewards.Plus(validator.Reward)
	}
	if commitment.TotalReward.IsLT(validatorRewards) {
		return false
	}

	return true
}

func (commitment Commitment) Hash() CommitmentId {
	b, err := json.Marshal(commitment)
	if err != nil {
		panic(err)
	}
	return sha3.Sum256(b)
}

func NewCommitmentFromBid(bid Bid, publisher sdk.AccAddress, validUntil int64, extraValidator crypto.PubKey) Commitment {
	validators := bid.Validators
	if extraValidator != nil {
		validators = append(validators, Validator{
			PubKey: extraValidator,
			// The extra validator should not be allowed to set their own reward
			Reward: sdk.Coins{},
		})
	}
	return Commitment{
		BidId: bid.Hash(),
		TotalReward: bid.TotalReward,
		ValidUntil: validUntil,
		Publisher: publisher,
		Advertiser: bid.GetAdvertiserAddress(),
		Validators: validators,
	}
}

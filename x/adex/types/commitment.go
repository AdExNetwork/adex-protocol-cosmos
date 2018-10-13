package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/crypto/sha3"
	"github.com/tendermint/tendermint/crypto"
	"encoding/json"
	"errors"
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
func (commitment Commitment) Validate() error {
	// @TODO: figure out if we need the nil ref/slice checks; it all depends on what happens when deserializing
	if commitment.Validators == nil || commitment.TotalReward == nil || commitment.Advertiser == nil || commitment.Publisher == nil {
		return errors.New("unexpected nil value")
	}
	if len(commitment.Validators) < minValidatorCount {
		return errors.New("insufficient number of validators")
	}
	for _, validator := range commitment.Validators {
		if !validator.IsValid() {
			return errors.New("invalid validator")
		}
	}

	if commitment.ValidUntil <= 0 {
		return errors.New("commitment validUntil must be positive")
	}

	if !commitment.TotalReward.IsNotNegative() {
		return errors.New("commitment TotalReward must be positive")
	}

	validatorRewards := sdk.Coins{}
	for _, validator := range commitment.Validators {
		validatorRewards = validatorRewards.Plus(validator.Reward)
	}
	if commitment.TotalReward.IsLT(validatorRewards) {
		return errors.New("TotalReward must be more than validatorRewards")
	}

	return nil
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

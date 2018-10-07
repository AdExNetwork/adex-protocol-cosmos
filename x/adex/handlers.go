package adex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	costCommitmentStart = 3000
	costCommitmentFinalize = 5000
)

func NewHandler(k bank.Keeper, ak Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
			case types.CommitmentStartMsg:
				return handleCommitmentStart(k, ak, ctx, msg)
			case types.CommitmentFinalizeMsg:
				return handleCommitmentFinalize(k, ak, ctx, msg)
			default:
				errMsg := "Unrecognized adex Msg type"
				return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCommitmentStart(k bank.Keeper, ak Keeper, ctx sdk.Context, msg types.CommitmentStartMsg) sdk.Result {
	ctx.GasMeter().ConsumeGas(costCommitmentStart, "commitmentStart")

	bidId := msg.Bid.Hash()
	if ak.GetBidState(ctx, bidId) != types.BidStateUnknown {
		return sdk.ErrUnknownRequest("a commitment for this bid already exists").Result()
	}
	validUntil := ctx.BlockHeader().Time.Unix() + msg.Bid.Timeout
	commitment := types.NewCommitmentFromBid(msg.Bid, msg.Publisher, validUntil, msg.ExtraValidatorAddr)
	if !commitment.IsValid() {
		return sdk.ErrUnknownRequest("commitment is not valid").Result()
	}

	ak.SetBidActive(ctx, bidId, commitment.Hash(), commitment.ValidUntil)

	_, _, err := k.SubtractCoins(ctx, commitment.Advertiser, commitment.TotalReward)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleCommitmentFinalize(k bank.Keeper, ak Keeper, ctx sdk.Context, msg types.CommitmentFinalizeMsg) sdk.Result {
	ctx.GasMeter().ConsumeGas(costCommitmentFinalize, "commitmentFinalize")

	commitmentId := msg.Commitment.Hash()
	if !ak.IsBidActive(ctx, msg.Commitment.BidId, commitmentId) {
		return sdk.ErrUnknownRequest("there is no active bid with that commitment").Result()
	}

	// @TODO: check signatures for each validator here
	rewardedValidators := msg.Commitment.Validators

	var newState uint8
	var rewardRecepient sdk.AccAddress
	if len(msg.Vote) == 1 && msg.Vote[0] == 0 {
		newState = types.BidStateFailed
		rewardRecepient = msg.Commitment.Advertiser
	} else {
		newState = types.BidStateSucceeded
		rewardRecepient = msg.Commitment.Publisher
	}

	// Mark the bid as failed/suceeded
	ak.SetBidState(ctx, msg.Commitment.BidId, newState)

	// Distribute rewards
	remainingReward := msg.Commitment.TotalReward
	for _, validator := range rewardedValidators {
		remainingReward = remainingReward.Minus(validator.Reward)
		_, _, err := k.AddCoins(ctx, validator.Address, validator.Reward)
		if err != nil {
			return err.Result()
		}
	}
	_, _, err := k.AddCoins(ctx, rewardRecepient, remainingReward)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// @TODO handle timeout on endblocker
// @TODO: handle bid cancel

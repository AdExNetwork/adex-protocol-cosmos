package adex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	signedmsg "github.com/cosmos/cosmos-sdk/adex/x/adex/signedmsg"
)

const (
	costBidCancel = 1000
	costCommitmentStart = 3000
	costCommitmentFinalize = 5000
)

func NewHandler(k bank.Keeper, ak Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
			case types.BidCancelMsg:
				return handleBidCancel(ak, ctx, msg)
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

func handleBidCancel(ak Keeper, ctx sdk.Context, msg types.BidCancelMsg) sdk.Result {
	ctx.GasMeter().ConsumeGas(costBidCancel, "bidCancel")

	bidId := msg.Bid.Hash()
	if ak.GetBidState(ctx, bidId) != types.BidStateUnknown {
		return sdk.ErrUnknownRequest("a commitment for this bid already exists").Result()
	}

	ak.SetBidState(ctx, bidId, types.BidStateCanceled)

	return sdk.Result{}
}

func handleCommitmentStart(k bank.Keeper, ak Keeper, ctx sdk.Context, msg types.CommitmentStartMsg) sdk.Result {
	ctx.GasMeter().ConsumeGas(costCommitmentStart, "commitmentStart")

	bidId := msg.Bid.Hash()
	if ak.GetBidState(ctx, bidId) != types.BidStateUnknown {
		return sdk.ErrUnknownRequest("a commitment for this bid already exists").Result()
	}
	validUntil := ctx.BlockHeader().Time.Unix() + msg.Bid.Timeout
	commitment := types.NewCommitmentFromBid(msg.Bid, msg.Publisher, validUntil, msg.ExtraValidatorPubKey)
	if !commitment.IsValid() {
		// @TODO: detailed info on why the commitment is not valid
		return sdk.ErrUnknownRequest("commitment is not valid").Result()
	}

	ak.SetBidActiveCommitment(ctx, bidId, commitment)

	_, _, err := k.SubtractCoins(ctx, commitment.Advertiser, commitment.TotalReward)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleCommitmentFinalize(k bank.Keeper, ak Keeper, ctx sdk.Context, msg types.CommitmentFinalizeMsg) sdk.Result {
	ctx.GasMeter().ConsumeGas(costCommitmentFinalize, "commitmentFinalize")

	// check if the state is correct
	commitmentId := msg.Commitment.Hash()
	if !ak.IsBidActiveCommitment(ctx, msg.Commitment.BidId, commitmentId) {
		return sdk.ErrUnknownRequest("there is no active bid with that commitment").Result()
	}

	// Check signatures for each validator
	expectSigned := append(commitmentId[:], msg.Vote...)
	validatorsWhoVoted := make([]types.Validator, 0)
	for i, validator := range msg.Commitment.Validators {
		if signedmsg.IsSigned(validator.PubKey, expectSigned, msg.Signatures[i]) {
			validatorsWhoVoted = append(validatorsWhoVoted, validator)
		}
	}
	if len(validatorsWhoVoted)*3 < len(msg.Commitment.Validators)*2 {
		return sdk.ErrUnknownRequest("not enough valid signatures: 2/3 of validators or more required").Result()
	}

	// a vote of 1 zero byte means failure to deliver
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
	for _, validator := range validatorsWhoVoted {
		remainingReward = remainingReward.Minus(validator.Reward)
		_, _, err := k.AddCoins(ctx, validator.GetAccAddress(), validator.Reward)
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

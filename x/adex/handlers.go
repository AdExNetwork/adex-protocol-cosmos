package adex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"log"
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
	// @TODO: remove this
	log.Println(msg)

	// then on finalize/timeout, assuming the commitment exists, we distribute the balances back; we should credit the validator rewards for validators
	// that did not sign back to the advertiser
	ctx.GasMeter().ConsumeGas(costCommitmentFinalize, "commitmentFinalize")
	// @TODO always do a transfer of coins, i.e. check if someone has the balance
	// unlike solidity, functions here won't revert() under you, so everything must be checked at a top level
	return sdk.Result{}
}

// @TODO handle timeout

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

func NewHandler(k bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
			case types.CommitmentStartMsg:
				return handleCommitmentStart(k, ctx, msg)
			case types.CommitmentFinalizeMsg:
				return handleCommitmentFinalize(k, ctx, msg)
			default:
				errMsg := "Unrecognized adex Msg type"
				return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCommitmentStart(k bank.Keeper, ctx sdk.Context, msg types.CommitmentStartMsg) sdk.Result {
	// @TODO: remove this
	log.Println(msg, msg.Bid.Hash())

	// @TODO: more granular
	ctx.GasMeter().ConsumeGas(costCommitmentStart, "commitmentStart")
	// @TODO can we do Bid.GetTotalReward()
	// k.HasCoins(ctx, msg.Bid.Advertiser, msg.Bid.Reward)
	// for validator := msg.Bid.Validators

	// @TODO: since we presume the bid is valid (cause Validatebasic on the msg). we construct a commitment and check if that is valid
	// then, we check if the advertiser has all the balances for bid.Reward
	// after we construct the commitment, check if the commitment.GetTotalReward() is less than Bid.Reward .IsGTE, .IsLT
	// if they do, we proceed to deduct them and mark the commitment as existant
	// then on finalize/timeout, assuming the commitment exists, we distribute the balances back; we should credit the validator rewards for validators
	// that did not sign back to the advertiser
	_, _, err := k.AddCoins(ctx, msg.Publisher, sdk.Coins{{"adex", sdk.NewInt(20)}})
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleCommitmentFinalize(k bank.Keeper, ctx sdk.Context, msg types.CommitmentFinalizeMsg) sdk.Result {
	// @TODO: remove this
	log.Println(msg)
	ctx.GasMeter().ConsumeGas(costCommitmentFinalize, "commitmentFinalize")
	// @TODO always do a transfer of coins, i.e. check if someone has the balance
	// unlike solidity, functions here won't revert() under you, so everything must be checked at a top level
	return sdk.Result{}
}

// @TODO handle timeout

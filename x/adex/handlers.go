package adex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"log"
)

func NewHandler(k bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
			case types.CommitmentStartMsg:
				return handleCommitmentStart(k, ctx, msg)
			// @TODO: fianlize
			default:
				errMsg := "Unrecognized adex Msg type"
				return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCommitmentStart(k bank.Keeper, ctx sdk.Context, msg types.CommitmentStartMsg) sdk.Result {
	// @TODO: remove this
	log.Println(msg)

	_, _, err := k.AddCoins(ctx, msg.Publisher, sdk.Coins{{"adex", sdk.NewInt(20)}})
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

//func handleCommitmentFinalize()
// @TODO always do a transfer of coins, i.e. check if someone has the balance
// unlike solidity, functions here won't revert() under you, so everything must be checked at a top level

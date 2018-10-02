package adex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func NewHandler(k bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
			case types.ClaimTokenMsg:
				return handleClaimToken(k, ctx, msg)
			default:
				errMsg := "Unrecognized adex Msg type"
				return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleClaimToken(k bank.Keeper, ctx sdk.Context, msg types.ClaimTokenMsg) sdk.Result {
	// @TODO: validate ethereum sigs, etc.
	_, _, err := k.AddCoins(ctx, msg.Submitter, sdk.Coins{{"adex", sdk.NewInt(msg.Amount)}})
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

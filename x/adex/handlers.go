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

func handleClaimToken(k bank.Keeper, ctx sdk.Context, msg sdk.Msg) sdk.Result {
	//k.AddCoins(ctx, submitter, coins)
	return sdk.Result{}
}

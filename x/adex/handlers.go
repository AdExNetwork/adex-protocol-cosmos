package adex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

func NewHandler(k bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
			errMsg := "Unrecognized adex Msg type"
			return sdk.ErrUnknownRequest(errMsg).Result()
	}
}

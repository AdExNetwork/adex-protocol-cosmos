package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Validator struct {
	Address sdk.AccAddress `json:"address"`
	Reward sdk.Coins `json:"reward"`
}

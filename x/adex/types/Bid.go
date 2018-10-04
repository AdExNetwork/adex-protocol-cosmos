package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Bid struct {
	Advertiser sdk.AccAddress `json:"advertiser"`
	// @TODO: adUnit, goal
	Timeout int64 `json:"timeout"`
	Reward sdk.Coins `json:"reward"`
	Nonce uint `json:"nonce"`
	Validators []Validator `json:"validators"`
}

// @TODO; validation

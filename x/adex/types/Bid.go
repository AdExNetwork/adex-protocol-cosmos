package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Bid struct {
	Advertiser sdk.AccAddress `json:"advertiser"`
	// @TODO: adUnit, goal
	Timeout int64 `json:"timeout"`
	TotalReward sdk.Coins `json:"totalReward"`
	Nonce uint `json:"nonce"`
	Validators []Validator `json:"validators"`
}

// @TODO: GetHash()
// @TODO: IsValid() - not sure if this is needed

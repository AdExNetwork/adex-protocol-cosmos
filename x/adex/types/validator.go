package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type Validator struct {
	PubKey crypto.PubKey `json:"pubkey"`
	Reward sdk.Coins `json:"reward"`
}

func (v Validator) IsValid() bool {
	return v.PubKey != nil && v.Reward != nil && v.Reward.IsNotNegative()
}

func (v Validator) GetAccAddress() sdk.AccAddress {
	return sdk.AccAddress(v.PubKey.Address())
}

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth"
)

type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

type GenesisAccount struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

func NewGenesisAccount(aa *auth.BaseAccount) *GenesisAccount {
	return &GenesisAccount{
		Address: aa.Address,
		Coins:   aa.Coins.Sort(),
	}
}

func (ga *GenesisAccount) ToAppAccount() (acc *auth.BaseAccount, err error) {
	return &auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}, nil
}

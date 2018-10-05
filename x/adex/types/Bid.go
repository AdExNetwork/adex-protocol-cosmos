package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/crypto/sha3"
)

type Bid struct {
	Advertiser sdk.AccAddress `json:"advertiser"`
	// @TODO: adUnit, goal
	TotalReward sdk.Coins `json:"totalReward"`
	Timeout int64 `json:"timeout"`
	Nonce uint `json:"nonce"`
	Validators []Validator `json:"validators"`
}

func (bid Bid) IsValid() bool {
	// @TODO: nonce > 0, timeout is valid, etc.
	return true
}

func (bid Bid) Hash() [32]byte {
	b, err := json.Marshal(bid)
	if err != nil {
		panic(err)
	}
	return sha3.Sum256(b)
}

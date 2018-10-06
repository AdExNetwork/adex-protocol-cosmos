package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/crypto/sha3"
)

const (
	// 1 year in seconds
	maxTimeout = 31536000
)

type BidId [32]byte

type Bid struct {
	Advertiser sdk.AccAddress `json:"advertiser"`
	// @TODO: adUnit, goal
	TotalReward sdk.Coins `json:"totalReward"`
	Timeout int64 `json:"timeout"`
	Nonce int64 `json:"nonce"`
	Validators []Validator `json:"validators"`
}

func (bid Bid) IsValid() bool {
	return bid.Timeout > 0 && bid.Timeout < maxTimeout && !bid.Advertiser.Empty()
}

func (bid Bid) Hash() BidId {
	b, err := json.Marshal(bid)
	if err != nil {
		panic(err)
	}
	return sha3.Sum256(b)
}

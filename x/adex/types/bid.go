package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/crypto/sha3"
	"github.com/tendermint/tendermint/crypto"
	signedmsg "github.com/cosmos/cosmos-sdk/adex/x/adex/signedmsg"
	"errors"
)

const (
	// 1 year in seconds
	maxTimeout = 31536000
)

type BidId [32]byte

type Bid struct {
	AdvertiserPubKey crypto.PubKey `json:"advertiser"`
	AdUnit []byte `json:"adUnit"`
	Goal []byte `json:"byte"`
	TotalReward sdk.Coins `json:"totalReward"`
	Timeout int64 `json:"timeout"`
	Nonce int64 `json:"nonce"`
	Validators []Validator `json:"validators"`
}

func (bid Bid) Validate() error {
	if !(bid.Timeout > 0 && bid.Timeout < maxTimeout) {
		return errors.New("invalid bid timeout")
	}
	if bid.AdvertiserPubKey == nil {
		return errors.New("AdvertiserPubKey required")
	}
	return nil
}

func (bid Bid) Hash() BidId {
	b, err := json.Marshal(bid)
	if err != nil {
		panic(err)
	}
	return sha3.Sum256(b)
}

func (bid Bid) GetAdvertiserAddress() sdk.AccAddress {
	return sdk.AccAddress(bid.AdvertiserPubKey.Address())
}

func (bid Bid) IsValidSignature(sig []byte) bool {
	bidId := bid.Hash()
	return signedmsg.IsSigned(bid.AdvertiserPubKey, bidId[:], sig)
}

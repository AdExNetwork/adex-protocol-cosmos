package types

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	errors "github.com/cosmos/cosmos-sdk/adex/x/adex/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CommitmentStartMsg struct {
	// the Bid is created and signed by the Advertiser
	Bid Bid `json:"bid"`
	BidSig []byte `json:"bidSig"`
	ExtraValidatorPubKey crypto.PubKey `json:"extraValidatorPubKey"`
	// and accepted (turned into a Commitment) by the publisher
	Publisher sdk.AccAddress `json:"publisher"`
}

func (msg CommitmentStartMsg) Name() string {
	return "CommitmentStartMsg"
}

func (msg CommitmentStartMsg) Type() string {
	return "adex"
}

func (msg CommitmentStartMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg CommitmentStartMsg) GetSigners() []sdk.AccAddress {
	// This message should be signed by the publisher, but the advertiser previously signed the Bid
	return []sdk.AccAddress{msg.Publisher}
}

func (msg CommitmentStartMsg) ValidateBasic() sdk.Error {
	errBid := msg.Bid.Validate()
	if errBid != nil {
		return errors.ErrInvalidBid(errors.DefaultCodespace, errBid)
	}
	if !msg.Bid.IsValidSignature(msg.BidSig) {
		return errors.ErrInvalidBidSignature(errors.DefaultCodespace)
	}
	return nil
}

func (msg CommitmentStartMsg) String() string {
	return fmt.Sprintf("CommitmentStartMsg{%v, %v}", msg.Bid, msg.Publisher)
}

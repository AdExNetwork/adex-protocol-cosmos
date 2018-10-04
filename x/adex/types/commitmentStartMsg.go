package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"log"
)

const (
	MAX_EXTRA_VALIDATORS = 1
)

type CommitmentStartMsg struct {
	// the Bid is created and signed by the Advertiser
	Bid Bid `json:"bid"`
	BidSig []byte `json:"bidSig"`
	// and accepted (turned into a Commitment) by the publisher
	Publisher sdk.AccAddress `json:"publisher"`
	ExtraValidators []Validator `json:"extraValidators"`
}

func (msg CommitmentStartMsg) Name() string {
	return "CommitmentStartMsg"
}

func (msg CommitmentStartMsg) Type() string {
	return "adex"
}

func (msg CommitmentStartMsg) GetSignBytes() []byte {
	// @TODO: try amino for this shit
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	// @TODO: remove this; currently for debugging
	log.Println(string(b))
	return b
}

func (msg CommitmentStartMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Publisher}
}

func (msg CommitmentStartMsg) ValidateBasic() sdk.Error {
	if len(msg.ExtraValidators) > MAX_EXTRA_VALIDATORS {
		// @TODO: our own error
		return sdk.ErrUnknownRequest("invalid amount")
	}
	// @TODO: call .Bid.IsValid()
	// @TODO: validate the sig

	// NOTE: the .Publisher must be the signer of the mssage (see GetSigners)
	return nil
}

func (msg CommitmentStartMsg) String() string {
	return fmt.Sprintf("CommitmentStartMsg{%v, %v}", msg.Bid, msg.Publisher)
}

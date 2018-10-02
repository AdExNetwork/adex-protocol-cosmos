package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ClaimTokenMsg struct {
	Amount	int64	`json:"amount"`
	Submitter sdk.AccAddress `json:"submitter"`
	// @TODO ethSig, ethAddr; the ethSig must be signing the cosmos addr of the requester
}

// Implement msg
func (msg ClaimTokenMsg) Type() string {
	return "adex"
}

// Get Implements Msg
func (msg ClaimTokenMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes Implements Msg
func (msg ClaimTokenMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners Implements Msg
func (msg ClaimTokenMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Submitter}
}

// ValidateBasic Implements Msg
func (msg ClaimTokenMsg) ValidateBasic() sdk.Error {
	if msg.Amount < 0 {
		return sdk.ErrInvalidCoins("invalid amount")
	}

	return nil
}

func (msg ClaimTokenMsg) String() string {
	return fmt.Sprintf("ClaimTokenMsg{%v}", msg.Amount)
}

// @TODO: implement the msg interface

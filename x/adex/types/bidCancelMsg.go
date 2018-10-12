package types

import (
	"encoding/json"
	"fmt"
	errors "github.com/cosmos/cosmos-sdk/adex/x/adex/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BidCancelMsg struct {
	Bid Bid `json:"bid"`
}

func (msg BidCancelMsg) Name() string {
	return "BidCancelMsg"
}

func (msg BidCancelMsg) Type() string {
	return "adex"
}

func (msg BidCancelMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg BidCancelMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Bid.GetAdvertiserAddress()}
}

func (msg BidCancelMsg) ValidateBasic() sdk.Error {
	// @TODO: think of skipping validation here, in case we upgraded the app and we want to cancel bids that are no longer valid
	if !msg.Bid.IsValid() {
		return errors.ErrInvalidBid(errors.DefaultCodespace, "IsValid failed")
	}
	return nil
}

func (msg BidCancelMsg) String() string {
	return fmt.Sprintf("BidCancelMsg{%v, %v}", msg.Bid)
}

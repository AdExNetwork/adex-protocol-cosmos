package signedmsg

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func IsSigned(addr sdk.AccAddress, data []byte, sig []byte) bool {
	// @TODO
	return false
}

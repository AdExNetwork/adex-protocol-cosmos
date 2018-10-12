package errors

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 42

	// adex module reserves error 4200-4300
	CodeInvalidBid sdk.CodeType = 400
	CodeInvalidBidSignature sdk.CodeType = 401
	CodeInvalidCommitment sdk.CodeType = 402
	CodeInvalidSigCount sdk.CodeType = 403
	CodeInvalidVote sdk.CodeType = 404
	CodeUnexpectedBidState sdk.CodeType = 405
	CodeNotEnoughSignatures sdk.CodeType = 406
)

func ErrInvalidBid(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBid, fmt.Sprintf("invalid bid: %v", msg))
}

func ErrInvalidBidSignature(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBidSignature, "invalid bid signature")
}

func ErrInvalidCommitment(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidCommitment, fmt.Sprintf("invalid commitment: %v", msg))
}

func ErrInvalidSigCount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSigCount, "wrong signature count: must be same as number of validators")
}

func ErrInvalidVote(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidVote, fmt.Sprintf("invalid vote: %v", msg))
}

func ErrUnexpectedBidState(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnexpectedBidState, fmt.Sprintf("unexpected bid state: %v", msg))
}

func ErrNotEnoughSignatures(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNotEnoughSignatures, msg)
}


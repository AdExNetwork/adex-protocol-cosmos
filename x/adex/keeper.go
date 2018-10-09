package adex

import (
        sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"encoding/binary"
	"bytes"
        codec "github.com/cosmos/cosmos-sdk/codec"
)

const (
	bidStateKey = "bidState"
	validUntilKey = "bidValidUntil"
)

type Keeper struct {
        storeKey  sdk.StoreKey
	cdc       *codec.Codec
}

func NewKeeper(key sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{ storeKey: key, cdc: cdc }
}

func (k Keeper) SetBidState(ctx sdk.Context, bidId types.BidId, state uint8) {
        store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
	// @TODO: validate if between 2 and 5
        store.Set(bidId[:], []byte{ byte(state) })
}

func (k Keeper) SetBidActiveCommitment(ctx sdk.Context, bidId types.BidId, commitment types.Commitment) {
        store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
	commitmentId := commitment.Hash()
	store.Set(bidId[:], commitmentId[:])
	k.setCommitmentValidUntil(ctx, bidId, commitment)
}

func (k Keeper) IsBidActiveCommitment(ctx sdk.Context, bidId types.BidId, commitmentId types.CommitmentId) bool {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
	expectedCommitmentId := store.Get(bidId[:])
	return bytes.Equal(expectedCommitmentId, commitmentId[:])
}

func (k Keeper) GetBidState(ctx sdk.Context, bidId types.BidId) uint8 {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
	val := store.Get(bidId[:])
	if len(val) == 0 {
		return types.BidStateUnknown
	}
	if len(val) == 32 {
		return types.BidStateActive
	}
	if len(val) == 1 {
		// @TODO: validate if between 2 (BidStateCanceled) and 5 (BidStateSuccess)
		return uint8(val[0])
	}
	panic("unknown bid state")
}

//
// Commitment expiry & refunds
//
type CommitmentRefunds struct {
	All []types.CommitmentRefund
}

func (k Keeper) CleanupCommitmentsExpiringBetween(ctx sdk.Context, start int64, end int64, f func(id types.CommitmentRefund) error) {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(validUntilKey))

	startB := make([]byte, 8)
	endB := make([]byte, 8)
	binary.LittleEndian.PutUint64(startB, uint64(start))
	binary.LittleEndian.PutUint64(endB, uint64(end))

	toCleanup := make([][]byte, 0)

	iter := store.Iterator(startB[:], endB[:])
	for {
		if !iter.Valid() {
			iter.Close()
			break
		}

		var refunds CommitmentRefunds
		k.cdc.MustUnmarshalBinary(iter.Value(), &refunds)

		for _, refund := range refunds.All {
			// Errors on cleaning up commitments are absolutely not allowed to happen
			err := f(refund)
			if err != nil {
				panic(err)
			}
		}

		toCleanup = append(toCleanup, iter.Key())

		iter.Next()
	}

	for _, key := range toCleanup {
		store.Delete(key)
	}
}

func (k Keeper) setCommitmentValidUntil(ctx sdk.Context, bidId types.BidId, commitment types.Commitment) {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(validUntilKey))

	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(commitment.ValidUntil))

	var refunds CommitmentRefunds
	bz := store.Get(key)
	if bz != nil {
		k.cdc.MustUnmarshalBinary(bz, &refunds)
	}
	refunds.All = append(refunds.All, types.CommitmentRefund{
		BidId: commitment.BidId,
		TotalReward: commitment.TotalReward,
		Beneficiary: commitment.Advertiser,
	})
	bz = k.cdc.MustMarshalBinary(refunds)
	store.Set(key, bz)
}

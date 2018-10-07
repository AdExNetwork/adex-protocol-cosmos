package adex

import (
        sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"encoding/binary"
	"bytes"
)

const (
	bidStateKey = "bidState"
	validUntilKey = "bidValidUntil"
)

type Keeper struct {
        storeKey  sdk.StoreKey
}

func NewKeeper(key sdk.StoreKey) Keeper {
	return Keeper{ storeKey: key }
}

func (k Keeper) SetBidState(ctx sdk.Context, bidId types.BidId, state uint8) {
        store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
	// @TODO: validate if between 2 and 5
        store.Set(bidId[:], []byte{ byte(state) })
}

func (k Keeper) SetBidActive(ctx sdk.Context, bidId types.BidId, commitmentId types.CommitmentId, validUntil int64) {
        store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
	store.Set(bidId[:], commitmentId[:])
	k.setCommitmentValidUntil(ctx, bidId, validUntil)
}

func (k Keeper) IsBidActive(ctx sdk.Context, bidId types.BidId, commitmentId types.CommitmentId) bool {
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
		// @TODO: validate if between 2 (BidStateCancelled) and 5 (BidStateSuccess)
		return uint8(val[0])
	}
	panic("unknown bid state")
}

func (k Keeper) IterateCommitmentsExpiringBetween(ctx sdk.Context, start int64, end int64, f func(id types.BidId)) {
	// @TODO: perhaps we should make our own iterator, where we unfold each value into bidId's
	store := ctx.KVStore(k.storeKey).Prefix([]byte(validUntilKey))
	startB := make([]byte, 8)
	endB := make([]byte, 8)
	binary.LittleEndian.PutUint64(startB, uint64(start))
	binary.LittleEndian.PutUint64(endB, uint64(end))
	iter := store.Iterator(startB[:], endB[:])
	for {
		if !iter.Valid() {
			iter.Close()
			break
		}

		allIds := iter.Value()
		if len(allIds) % 32 != 0 {
			panic("invalid data in the validUntil store")
		}

		for i := 0; i < len(allIds)/32; i++ {
			start := int64(i*32)
			end = start+32
			id := types.BidId{}
			copy(id[:], allIds[start:end])
			f(id)
		}

		iter.Next()
	}
}

func (k Keeper) setCommitmentValidUntil(ctx sdk.Context, bidId types.BidId, validUntil int64) {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(validUntilKey))

	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(validUntil))

	bids := store.Get(key)
	bids = append(bids, bidId[:]...)
	store.Set(key, bids)
}

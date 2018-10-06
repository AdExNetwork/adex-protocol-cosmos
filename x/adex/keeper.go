package adex

import (
        sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"encoding/binary"
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
	k.setBidValidUntil(ctx, bidId, validUntil)
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

// @TODO: how can we safely cast here?
func (k Keeper) GetValidUntilIter(ctx sdk.Context, start int64, end int64) sdk.Iterator {
	// @TODO: perhaps we should make our own iterator, where we unfold each value into bidId's
	store := ctx.KVStore(k.storeKey).Prefix([]byte(validUntilKey))
	startB := make([]byte, 8)
	endB := make([]byte, 8)
	binary.LittleEndian.PutUint64(startB, uint64(start))
	binary.LittleEndian.PutUint64(endB, uint64(end))
	return store.Iterator(startB[:], endB[:])
}

func (k Keeper) setBidValidUntil(ctx sdk.Context, bidId types.BidId, validUntil int64) {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(validUntilKey))

	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(validUntil))

	bids := store.Get(key)
	bids = append(bids, bidId[:]...)
	store.Set(key, bids)
}

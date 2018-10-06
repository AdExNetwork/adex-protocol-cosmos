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
	k.SetBidValidUntil(ctx, bidId, validUntil)
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

// @TODO: make this private
// @TODO: how can we safely cast here?
func (k Keeper) SetBidValidUntil(ctx sdk.Context, bidId types.BidId, validUntil int64) {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(validUntilKey))
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(validUntil))
	// @TODO: add the bidId on to the current value
	store.Set(b, bidId[:])
}

func (k Keeper) GetValidUntilIter(ctx sdk.Context, start int64, end int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey).Prefix([]byte("bidValidUntil"))
	startB := make([]byte, 8)
	endB := make([]byte, 8)
	binary.LittleEndian.PutUint64(startB, uint64(start))
	binary.LittleEndian.PutUint64(endB, uint64(end))
	return store.Iterator(startB[:], endB[:])
}

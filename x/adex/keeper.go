package adex

import (
        sdk "github.com/cosmos/cosmos-sdk/types"
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

// @TODO: think of making the BidId a type
// @TODO: think of passing structs here?
func (k Keeper) SetBidState(ctx sdk.Context, bidId [32]byte, state byte) {
        store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
        store.Set(bidId[:], []byte{ state })
}

func (k Keeper) GetBidState(ctx sdk.Context, bidId [32]byte) byte {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
	return store.Get(bidId[:])[0]
}

// @TODO: think of safely casting int to uint32
func (k Keeper) SetBidValidUntil(ctx sdk.Context, bidId [32]byte, validUntil uint32) {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(validUntilKey))
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, validUntil)
	// @TODO: add the bidId on to the value
	store.Set(b, bidId[:])
}

func (k Keeper) GetValidUntilIter(ctx sdk.Context, start uint32, end uint32) sdk.Iterator {
	// @TODO: safe cast, binary.Little...
	store := ctx.KVStore(k.storeKey).Prefix([]byte("bidValidUntil"))
	startB := make([]byte, 4)
	endB := make([]byte, 4)
	binary.LittleEndian.PutUint32(startB, start)
	binary.LittleEndian.PutUint32(endB, end)
	return store.Iterator(startB[:], endB[:])
}

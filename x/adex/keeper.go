package adex

import (
        sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	bidStateKey = "bidState"
	validUntilKey = "commitmentValidUntil"
)

type Keeper struct {
        storeKey  sdk.StoreKey
}

func NewKeeper(key sdk.StoreKey) Keeper {
	return Keeper{ storeKey: key }
}

func (k Keeper) SetBidState(ctx sdk.Context, bidId [32]byte, state byte) {
        store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
        store.Set(bidId[:], []byte{ state })
}

func (k Keeper) GetBidState(ctx sdk.Context, bidId [32]byte) byte {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
	return store.Get(bidId[:])[0]
}

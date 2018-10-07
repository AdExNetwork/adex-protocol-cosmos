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

func (k Keeper) SetBidActiveCommitment(ctx sdk.Context, bidId types.BidId, commitment types.Commitment) {
        store := ctx.KVStore(k.storeKey).Prefix([]byte(bidStateKey))
	commitmentId := commitment.Hash()
	store.Set(bidId[:], commitmentId[:])
	k.setCommitmentValidUntil(ctx, bidId, commitment)
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

func (k Keeper) CleanupCommitmentsExpiringBetween(ctx sdk.Context, start int64, end int64, f func(id types.CommitmentRefund) error) {
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
			// @TODO
			f(types.CommitmentRefund{ BidId: id, TotalReward: sdk.Coins{}, Beneficiary: sdk.AccAddress{} })
		}

		iter.Next()
	}
	// @TODO: clean them up from the state tree
}

func (k Keeper) setCommitmentValidUntil(ctx sdk.Context, bidId types.BidId, commitment types.Commitment) {
	store := ctx.KVStore(k.storeKey).Prefix([]byte(validUntilKey))

	// @TODO: put types.CommitmentRefund
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(commitment.ValidUntil))

	bids := store.Get(key)
	bids = append(bids, bidId[:]...)
	store.Set(key, bids)
}

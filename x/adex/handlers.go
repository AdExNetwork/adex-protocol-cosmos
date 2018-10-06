package adex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"log"
)

const (
	costCommitmentStart = 3000
	costCommitmentFinalize = 5000
)

func NewHandler(k bank.Keeper, ak Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
			case types.CommitmentStartMsg:
				return handleCommitmentStart(k, ak, ctx, msg)
			case types.CommitmentFinalizeMsg:
				return handleCommitmentFinalize(k, ak, ctx, msg)
			default:
				errMsg := "Unrecognized adex Msg type"
				return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleCommitmentStart(k bank.Keeper, ak Keeper, ctx sdk.Context, msg types.CommitmentStartMsg) sdk.Result {
	// @NOTE start real impl
	// @TODO: more granular gas cost
	ctx.GasMeter().ConsumeGas(costCommitmentStart, "commitmentStart")
	bidId := msg.Bid.Hash()
	if ak.GetBidState(ctx, bidId) != types.BidStateUnknown {
		return sdk.ErrUnknownRequest("a commitment for this bid already exists").Result()
	}
	// @TODO: set validUntil
	commitment := types.NewCommitmentFromBid(msg.Bid, msg.Publisher, 0, msg.ExtraValidatorAddr)
	// @TODO; check commitment.IsValid()
	log.Println(commitment)
	// finally, k.SubtractCoins, ak.SetBidState, ak.SetBidValidUntil
	// @NOTE end real impl

	// @TODO: remove this test code
	ak.SetBidValidUntil(ctx, bidId, uint32(30))
	ak.SetBidValidUntil(ctx, bidId, uint32(8))
	ak.SetBidValidUntil(ctx, bidId, uint32(2))
	ak.SetBidValidUntil(ctx, bidId, uint32(40))
	ak.SetBidValidUntil(ctx, bidId, uint32(44))
	ak.SetBidValidUntil(ctx, bidId, uint32(300))
	ak.SetBidValidUntil(ctx, bidId, uint32(304))
	// end is exclusive
	iterator := ak.GetValidUntilIter(ctx, 3, 305)
	for ; ; {
		if !iterator.Valid() {
			iterator.Close()
			break
		}
		log.Println("found value at key", iterator.Key())
		iterator.Next()
	}
	log.Println("iteration finished!!")


	// for validator := msg.Bid.Validators

	// @TODO: since we presume the bid is valid (cause Validatebasic on the msg). we construct a commitment and check if that is valid
	// then, we check if the advertiser has all the balances for bid.Reward
	// after we construct the commitment, check if the commitment.GetTotalReward() is less than Bid.Reward .IsGTE, .IsLT
	// if they do, we proceed to deduct them and mark the commitment as existant
	// then on finalize/timeout, assuming the commitment exists, we distribute the balances back; we should credit the validator rewards for validators
	// that did not sign back to the advertiser
	_, _, err := k.AddCoins(ctx, msg.Publisher, sdk.Coins{{"adex", sdk.NewInt(20)}})
	if err != nil {
		return err.Result()
	}
	// @TODO: we can safely SubtractCoins here for the Commitment.TotalReward

	return sdk.Result{}
}

func handleCommitmentFinalize(k bank.Keeper, ak Keeper, ctx sdk.Context, msg types.CommitmentFinalizeMsg) sdk.Result {
	// @TODO: remove this
	log.Println(msg)
	ctx.GasMeter().ConsumeGas(costCommitmentFinalize, "commitmentFinalize")
	// @TODO always do a transfer of coins, i.e. check if someone has the balance
	// unlike solidity, functions here won't revert() under you, so everything must be checked at a top level
	return sdk.Result{}
}

// @TODO handle timeout

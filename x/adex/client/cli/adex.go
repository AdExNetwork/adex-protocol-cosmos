package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

// –––––––––––– Flags ––––––––––––––––

const (
	FlagAmount     = "amount"
)

func PostCmdCommitmentStart(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commitment",
		Short: "Start a commitment for a bid",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
                                WithCodec(cdc).
                                WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			publisher, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			//amount := viper.GetInt64(FlagAmount)
			// @TODO
			// @TODO: instead of empty slices, nil should be used
			// othrewise after encoding and decoding through amino, it still ends up as a nil slice
			validators := []types.Validator{
				types.Validator{ Address: publisher, Reward: sdk.Coins{{ "adex", sdk.NewInt(1) }} },
				types.Validator{ Address: publisher, Reward: sdk.Coins{{ "adex", sdk.NewInt(2) }} },
			}
			bid := types.Bid{
				Advertiser: publisher,
				Timeout: 20,
				Validators: validators,
				TotalReward: sdk.Coins{{ "adex", sdk.NewInt(99) }},
			}
			msg := types.CommitmentStartMsg{
				Bid: bid,
				BidSig: nil,
				Publisher: publisher,
				// @TODO: nil extra validator
				ExtraValidatorAddr: publisher,
			}

			cliCtx.PrintResponse = true
			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	// @TODO
	cmd.Flags().Int(FlagAmount, 10000, "Amount to claim")
	cmd.MarkFlagRequired(FlagAmount)

	return cmd
}

// @TODO: PostCmdCommitmentFinalize


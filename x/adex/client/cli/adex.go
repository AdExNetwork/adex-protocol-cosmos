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
	"github.com/spf13/viper"
)

// –––––––––––– Flags ––––––––––––––––

const (
	FlagAmount     = "amount"
)

func PostCmdClaimToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim",
		Short: "Claim your tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
                                WithCodec(cdc).
                                WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			submitter, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			amount := viper.GetInt64(FlagAmount)
			msg := types.ClaimTokenMsg{ Amount: amount, Submitter: submitter }

			cliCtx.PrintResponse = true
			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Int(FlagAmount, 10000, "Amount to claim")
	cmd.MarkFlagRequired(FlagAmount)

	return cmd
}


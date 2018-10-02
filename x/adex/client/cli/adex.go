package cli

import (
	"os"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	types "github.com/cosmos/cosmos-sdk/adex/x/adex/types"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// –––––––––––– Flags ––––––––––––––––

const (
	FlagAmount     = "amount"
)

func PostCmdClaimToken(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim",
		Short: "Claim your tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().
                                WithCodec(cdc).
                                WithLogger(os.Stdout).
                                WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			submitter, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}

			amount := viper.GetInt64(FlagAmount)
			msg := types.ClaimTokenMsg{ Amount: amount, Submitter: submitter }

			ctx.PrintResponse = true
			return utils.SendTx(txCtx, ctx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Int(FlagAmount, 10000, "Amount to claim")
	cmd.MarkFlagRequired(FlagAmount)

	return cmd
}


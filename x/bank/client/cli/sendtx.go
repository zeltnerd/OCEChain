package cli

import (
	"github.com/OCEChain/OCEChain/client/context"
	"github.com/OCEChain/OCEChain/client/utils"
	"github.com/OCEChain/OCEChain/codec"
	sdk "github.com/OCEChain/OCEChain/types"
	authcmd "github.com/OCEChain/OCEChain/x/auth/client/cli"
	authtxb "github.com/OCEChain/OCEChain/x/auth/client/txbuilder"
	"github.com/OCEChain/OCEChain/x/bank/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo     = "to"
	flagAmount = "amount"
)

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create and sign a send tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			/* 			account, err := cliCtx.GetAccount(from)
			   			if err != nil {
			   				return err
			   			}

			   			// ensure account has enough coins
			   			if !account.GetCoins().IsAllGTE(coins) {
			   				return errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			   			} */

			// build and sign the transaction, then broadcast to Tendermint
			msg := client.CreateMsg(from, to, coins)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTo, "", "Address to send coins")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagAmount)

	return cmd
}

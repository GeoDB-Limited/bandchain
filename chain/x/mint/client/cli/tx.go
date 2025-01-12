package cli

import (
	"bufio"
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

const (
	flagReceiver = "receiver"
	flagAmount   = "amount"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	mintCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "mint transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	mintCmd.AddCommand(flags.PostCommands(GetCmdMsgWithdrawCoinsToAccFromTreasury(cdc))...)

	return mintCmd
}

// GetCmdMsgWithdrawCoinsToAccFromTreasury implements minting transaction command.
func GetCmdMsgWithdrawCoinsToAccFromTreasury(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw-coins (--receiver [receiver]) (--amount [amount])",
		Short: "Withdraw some coins for account",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			receiverStr, err := cmd.Flags().GetString(flagReceiver)
			if err != nil {
				return sdkerrors.Wrapf(err, "flag: %s", flagReceiver)
			}
			receiver, err := sdk.AccAddressFromBech32(receiverStr)
			if err != nil {
				return sdkerrors.Wrapf(err, "receiver: %s", receiverStr)
			}

			amountStr, err := cmd.Flags().GetString(flagAmount)
			if err != nil {
				return sdkerrors.Wrapf(err, "flag: %s", flagAmount)
			}
			amount, err := sdk.ParseCoins(amountStr)
			if err != nil {
				return sdkerrors.Wrapf(err, "amount: %s", amountStr)
			}

			msg := types.NewMsgWithdrawCoinsToAccFromTreasury(amount, receiver, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return sdkerrors.Wrapf(err, "amount: %s receiver: %s", amount, receiverStr)
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

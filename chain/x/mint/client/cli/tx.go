package cli

import (
	"bufio"
	"fmt"
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
	flagDepositor = "depositor"
	flagAmount    = "amount"
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
	mintCmd.AddCommand(flags.PostCommands(GetCmdMintTokens(cdc))...)

	return mintCmd
}

// GetCmdMintTokens implements minting transaction command.
func GetCmdMintTokens(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-tokens (--depositor [depositor]) (--amount [amount])",
		Short: "Mint some tokens for account",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			depositorStr, err := cmd.Flags().GetString(flagDepositor)
			if err != nil {
				return sdkerrors.Wrapf(err, fmt.Sprintf("flag: %s", flagDepositor))
			}
			depositor, err := sdk.AccAddressFromBech32(depositorStr)
			if err != nil {
				return sdkerrors.Wrapf(err, fmt.Sprintf("depositor: %s", depositorStr))
			}

			amountStr, err := cmd.Flags().GetString(flagAmount)
			if err != nil {
				return sdkerrors.Wrapf(err, fmt.Sprintf("flag: %s", flagAmount))
			}
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return sdkerrors.Wrapf(types.ErrInvalidAmountToMint, fmt.Sprintf("amount: %s", amountStr))
			}

			msg := types.NewMsgMintTokens(amount, depositor, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return sdkerrors.Wrapf(err, fmt.Sprintf("amount: %s depositor: %s", amount, depositorStr))
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagDepositor, "", "Depositor")
	cmd.Flags().String(flagAmount, "", "Amount to mint")

	return cmd
}

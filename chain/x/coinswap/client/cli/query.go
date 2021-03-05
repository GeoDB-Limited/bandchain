package cli

import (
	"encoding/json"
	"fmt"
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"net/http"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	coinswapCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the coinswap module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	coinswapCmd.AddCommand(flags.GetCommands(
		GetQueryCmdParams(storeKey, cdc),
		GetQueryCmdRate(storeKey, cdc),
	)...)
	return coinswapCmd
}

func printOutput(cliCtx context.CLIContext, cdc *codec.Codec, bz []byte, out interface{}) error {
	var result types.QueryResult
	if err := json.Unmarshal(bz, &result); err != nil {
		return err
	}
	if result.Status != http.StatusOK {
		return cliCtx.PrintOutput(result.Result)
	}
	cdc.MustUnmarshalJSON(result.Result, out)
	return cliCtx.PrintOutput(out)
}

// GetQueryCmdParams implements the query parameters command.
func GetQueryCmdParams(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "params",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", route, types.QueryParams))
			if err != nil {
				return err
			}
			return printOutput(cliCtx, cdc, bz, &types.Params{})
		},
	}
}

// todo maybe query with data???
func GetQueryCmdRate(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "rate [from-denom] [to-denom]",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", route, types.QueryRate))
			if err != nil {
				return err
			}
			return printOutput(cliCtx, cdc, bz, &types.QueryRateResult{})
		},
	}
}

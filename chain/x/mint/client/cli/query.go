package cli

import (
	"fmt"
	"github.com/GeoDB-Limited/odincore/chain/x/common/client/cli"
	"github.com/GeoDB-Limited/odincore/chain/x/mint/internal/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetQueryCmd returns the cli query commands for the minting module.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	mintingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the minting module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	mintingQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryParams(cdc),
			GetCmdQueryInflation(cdc),
			GetCmdQueryAnnualProvisions(cdc),
			GetCmdQueryEthIntegrationAddress(cdc),
			GetCmdQueryTreasuryPool(cdc),
		)...,
	)

	return mintingQueryCmd
}

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryInflation(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "inflation",
		Short: "Query the current minting inflation value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryInflation)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var inflation sdk.Dec
			if err := cdc.UnmarshalJSON(res, &inflation); err != nil {
				return err
			}

			return cliCtx.PrintOutput(inflation)
		},
	}
}

// GetCmdQueryAnnualProvisions implements a command to return the current minting
// annual provisions value.
func GetCmdQueryAnnualProvisions(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "annual-provisions",
		Short: "Query the current minting annual provisions value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAnnualProvisions)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var inflation sdk.Dec
			if err := cdc.UnmarshalJSON(res, &inflation); err != nil {
				return err
			}

			return cliCtx.PrintOutput(inflation)
		},
	}
}

func GetCmdQueryEthIntegrationAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "eth-integration-address",
		Short: "Query current ETH integration address",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryEthIntegrationAddress)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var ethIntegrationAddress string
			if err := cdc.UnmarshalJSON(res, &ethIntegrationAddress); err != nil {
				return err
			}
			return cliCtx.PrintOutput(ethIntegrationAddress)
		},
	}
}

// GetCmdQueryTreasuryPool returns the command for fetching treasury pool info
func GetCmdQueryTreasuryPool(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "treasury-pool",
		Args:  cobra.NoArgs,
		Short: "Query the amount of coins in the treasury pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bz, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTreasuryPool))
			if err != nil {
				return err
			}
			return cli.PrintOutput(cliCtx, cdc, bz, &[]sdk.Coin{})
		},
	}
}

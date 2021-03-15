package cli

import (
	"encoding/json"
	"github.com/GeoDB-Limited/odincore/chain/x/common/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"net/http"
)

func PrintOutput(cliCtx context.CLIContext, cdc *codec.Codec, bz []byte, out interface{}) error {
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

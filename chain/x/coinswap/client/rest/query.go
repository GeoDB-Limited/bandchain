package rest

import (
	"fmt"
	"github.com/GeoDB-Limited/odincore/chain/x/coinswap/types"
	clientcmn "github.com/GeoDB-Limited/odincore/chain/x/oracle/client/common"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"net/http"
)

func getParamsHandler(cliCtx context.CLIContext, route string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", route, types.QueryParams))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		clientcmn.PostProcessQueryResponse(w, cliCtx.WithHeight(height), bz)
	}
}

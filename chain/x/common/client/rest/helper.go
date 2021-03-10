package rest

import (
	"encoding/json"
	commontypes "github.com/GeoDB-Limited/odincore/chain/x/common/types"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func PostProcessQueryResponse(w http.ResponseWriter, cliCtx context.CLIContext, bz []byte) {
	var result commontypes.QueryResult
	if err := json.Unmarshal(bz, &result); err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.Status)
	rest.PostProcessResponse(w, cliCtx, result.Result)
}

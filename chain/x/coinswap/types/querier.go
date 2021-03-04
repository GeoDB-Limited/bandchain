package types

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"net/http"
)

// Query endpoints supported by the coinswap Querier.
const (
	QueryParams = "params"
	QueryRate   = "rate"
)

// QueryResult wraps querier result with HTTP status to return to application.
type QueryResult struct {
	Status int             `json:"status"`
	Result json.RawMessage `json:"result"`
}

// QueryOK creates and marshals a QueryResult instance with HTTP status OK.
func QueryOK(result interface{}) ([]byte, error) {
	return json.MarshalIndent(QueryResult{
		Status: http.StatusOK,
		Result: codec.MustMarshalJSONIndent(ModuleCdc, result),
	}, "", "  ")
}

// QueryBadRequest creates and marshals a QueryResult instance with HTTP status BadRequest.
func QueryBadRequest(result interface{}) ([]byte, error) {
	return json.MarshalIndent(QueryResult{
		Status: http.StatusBadRequest,
		Result: codec.MustMarshalJSONIndent(ModuleCdc, result),
	}, "", "  ")
}

// QueryNotFound creates and marshals a QueryResult instance with HTTP status NotFound.
func QueryNotFound(result interface{}) ([]byte, error) {
	return json.MarshalIndent(QueryResult{
		Status: http.StatusNotFound,
		Result: codec.MustMarshalJSONIndent(ModuleCdc, result),
	}, "", "  ")
}

type QueryRateResult struct {
	Rate sdk.Int `json:"rate"`
	From Denom   `json:"from"`
	To   Denom   `json:"to"`
}

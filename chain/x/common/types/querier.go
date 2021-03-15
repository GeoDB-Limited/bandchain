package types

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	"net/http"
)

// QueryResult wraps querier result with HTTP status to return to application.
type QueryResult struct {
	Status int             `json:"status"`
	Result json.RawMessage `json:"result"`
}

// QueryOK creates and marshals a QueryResult instance with HTTP status OK.
func QueryOK(cdc *codec.Codec, result interface{}) ([]byte, error) {
	return json.MarshalIndent(QueryResult{
		Status: http.StatusOK,
		Result: codec.MustMarshalJSONIndent(cdc, result),
	}, "", "  ")
}

// QueryBadRequest creates and marshals a QueryResult instance with HTTP status BadRequest.
func QueryBadRequest(cdc *codec.Codec, result interface{}) ([]byte, error) {
	return json.MarshalIndent(QueryResult{
		Status: http.StatusBadRequest,
		Result: codec.MustMarshalJSONIndent(cdc, result),
	}, "", "  ")
}

// QueryNotFound creates and marshals a QueryResult instance with HTTP status NotFound.
func QueryNotFound(cdc *codec.Codec, result interface{}) ([]byte, error) {
	return json.MarshalIndent(QueryResult{
		Status: http.StatusNotFound,
		Result: codec.MustMarshalJSONIndent(cdc, result),
	}, "", "  ")
}

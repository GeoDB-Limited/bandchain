package keeper_test

import (
	"encoding/json"
	commontypes "github.com/GeoDB-Limited/odincore/chain/x/common/types"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/GeoDB-Limited/odincore/chain/x/common/testapp"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/keeper"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

func TestQueryPendingRequests(t *testing.T) {
	_, ctx, k := testapp.CreateTestInput(true)

	// Add 3 requests
	k.SetRequestLastExpired(ctx, 40)
	k.SetRequest(ctx, 41, defaultRequest())
	k.SetRequest(ctx, 42, defaultRequest())
	k.SetRequest(ctx, 43, defaultRequest())
	k.SetRequestCount(ctx, 43)

	// Fulfill some requests
	k.SetReport(ctx, 41, types.NewReport(testapp.Validator1.ValAddress, true, nil))
	k.SetReport(ctx, 42, types.NewReport(testapp.Validator2.ValAddress, true, nil))

	q := keeper.NewQuerier(k)

	tests := []struct {
		name     string
		args     []string
		expected []types.RequestID
	}{
		{
			name:     "Get all pending requests",
			args:     []string{},
			expected: []types.RequestID{41, 42, 43},
		},
		{
			name:     "Get pending requests for validator1",
			args:     []string{testapp.Validator1.ValAddress.String()},
			expected: []types.RequestID{42, 43},
		},
		{
			name:     "Get pending requests for validator2",
			args:     []string{testapp.Validator2.ValAddress.String()},
			expected: []types.RequestID{41, 43},
		},
		{
			name:     "Get pending requests for validator3",
			args:     []string{testapp.Validator3.ValAddress.String()},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := q(ctx, append([]string{types.QueryPendingRequests}, tt.args...), abci.RequestQuery{})
			require.NoError(t, err)

			var queryRequest commontypes.QueryResult
			require.NoError(t, json.Unmarshal(raw, &queryRequest))

			var requestIDs []types.RequestID
			types.ModuleCdc.MustUnmarshalJSON(queryRequest.Result, &requestIDs)

			require.Equal(t, tt.expected, requestIDs)
		})
	}
}

package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func queryRequest(route string, cliCtx context.CLIContext, rid string) ([]byte, int64, error) {
	res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryRequests, rid), nil)
	if err != nil {
		return nil, 0, err
	}
	return res, height, nil
}

func QuerySearchLatestRequest(
	route string, cliCtx context.CLIContext, oid, calldata, askCount, minCount string,
) ([]byte, int64, error) {
	query := fmt.Sprintf("%s.%s='%s' AND %s.%s='%s' AND %s.%s='%s' AND %s.%s='%s'",
		types.EventTypeRequest, types.AttributeKeyOracleScriptID, oid,
		types.EventTypeRequest, types.AttributeKeyCalldata, calldata,
		types.EventTypeRequest, types.AttributeKeyAskCount, askCount,
		types.EventTypeRequest, types.AttributeKeyMinCount, minCount,
	)
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, 0, err
	}
	resTxs, err := node.TxSearch(query, !cliCtx.TrustNode, 1, 30, "desc")
	if err != nil {
		return nil, 0, err
	}
	for _, tx := range resTxs.Txs {
		if !cliCtx.TrustNode {
			err := utils.ValidateTxResult(cliCtx, tx)
			if err != nil {
				return nil, 0, err
			}
		}
		logs, _ := sdk.ParseABCILogs(tx.TxResult.Log)
		for _, log := range logs {
			for _, ev := range log.Events {
				if ev.Type != types.EventTypeRequest {
					continue
				}
				rid := ""
				ok := true
				for _, attr := range ev.Attributes {
					if attr.Key == types.AttributeKeyID {
						rid = attr.Value
					}
					if attr.Key == types.AttributeKeyOracleScriptID && attr.Value != oid ||
						attr.Key == types.AttributeKeyCalldata && attr.Value != calldata ||
						attr.Key == types.AttributeKeyAskCount && attr.Value != askCount ||
						attr.Key == types.AttributeKeyMinCount && attr.Value != minCount {
						ok = false
						break
					}
				}
				if ok && rid != "" {
					res, h, err := queryRequest(route, cliCtx, rid)
					if err != nil {
						return nil, 0, err
					}
					var out types.QueryRequestResult
					cliCtx.Codec.MustUnmarshalJSON(res, &out)
					if out.Result != nil {
						return res, h, nil
					}
				}
			}
		}
	}
	return nil, 0, fmt.Errorf("request not found")
}

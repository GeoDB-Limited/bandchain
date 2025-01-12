package keeper

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"

	"github.com/GeoDB-Limited/odincore/chain/pkg/bandrng"
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
	owasm "github.com/GeoDB-Limited/odincore/go-owasm/api"
)

// GetRandomValidators returns a pseudorandom subset of active validators. Each validator has
// chance of getting selected directly proportional to the amount of voting power it has.
func (k Keeper) GetRandomValidators(ctx sdk.Context, size int, id int64) ([]sdk.ValAddress, error) {
	valOperators := []sdk.ValAddress{}
	valPowers := []uint64{}
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx,
		func(idx int64, val exported.ValidatorI) (stop bool) {
			if k.GetValidatorStatus(ctx, val.GetOperator()).IsActive {
				valOperators = append(valOperators, val.GetOperator())
				valPowers = append(valPowers, val.GetTokens().Uint64())
			}
			return false
		})
	if len(valOperators) < size {
		return nil, sdkerrors.Wrapf(
			types.ErrInsufficientValidators, "%d < %d", len(valOperators), size)
	}
	rng, err := bandrng.NewRng(k.GetRollingSeed(ctx), sdk.Uint64ToBigEndian(uint64(id)), []byte(ctx.ChainID()))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrBadDrbgInitialization, err.Error())
	}
	tryCount := int(k.GetParamUint64(ctx, types.KeySamplingTryCount))
	chosenValIndexes := bandrng.ChooseSomeMaxWeight(rng, valPowers, size, tryCount)
	validators := make([]sdk.ValAddress, size)
	for i, idx := range chosenValIndexes {
		validators[i] = valOperators[idx]
	}
	return validators, nil
}

// PrepareRequest takes an request specification object, performs the prepare call, and saves
// the request object to store. Also emits events related to the request.
func (k Keeper) PrepareRequest(ctx sdk.Context, r types.RequestWithSenderSpec) error {
	askCount := r.GetAskCount()
	if askCount > k.GetParamUint64(ctx, types.KeyMaxAskCount) {
		return sdkerrors.Wrapf(types.ErrInvalidAskCount, "got: %d, max: %d", askCount, k.GetParamUint64(ctx, types.KeyMaxAskCount))
	}
	// Consume gas for data requests. We trust that we have reasonable params that don't cause overflow.
	ctx.GasMeter().ConsumeGas(k.GetParamUint64(ctx, types.KeyBaseRequestGas), "BASE_REQUEST_FEE")
	ctx.GasMeter().ConsumeGas(askCount*k.GetParamUint64(ctx, types.KeyPerValidatorRequestGas), "PER_VALIDATOR_REQUEST_FEE")
	// Get a random validator set to perform this request.
	validators, err := k.GetRandomValidators(ctx, int(askCount), k.GetRequestCount(ctx)+1)
	if err != nil {
		return err
	}
	// Create a request object. Note that RawRequestIDs will be populated after preparation is done.
	req := types.NewRequest(
		r.GetOracleScriptID(), r.GetCalldata(), validators, r.GetMinCount(),
		ctx.BlockHeight(), ctx.BlockTime(), r.GetClientID(), nil,
	)
	// Create an execution environment and call Owasm prepare function.
	env := types.NewPrepareEnv(req, int64(k.GetParamUint64(ctx, types.KeyMaxRawRequestCount)), int64(k.GetParamUint64(ctx, types.KeyMaxDataSize)))
	script, err := k.GetOracleScript(ctx, req.OracleScriptID)
	if err != nil {
		return err
	}
	code := k.GetFile(script.Filename)

	maxDataSize := k.GetParamUint64(ctx, types.KeyMaxDataSize)
	output, err := owasm.Prepare(code, types.WasmPrepareGas, int64(maxDataSize), env)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrBadWasmExecution, err.Error())
	}

	// calculate fee for each data source
	req.RawRequests = env.GetRawRequests()
	fee := k.GetDataRequesterBasicFeeParam(ctx)

	err = k.supplyKeeper.SendCoinsFromAccountToModule(ctx, r.GetSender(), types.ModuleName, sdk.NewCoins(fee.Value()))
	if err != nil {
		return sdkerrors.Wrap(err, "sending coins from account to module")
	}

	oraclePool := k.GetOraclePool(ctx)
	oraclePool.DataProvidersPool = oraclePool.DataProvidersPool.Add(sdk.NewDecCoinFromCoin(fee.Value()))
	k.SetOraclePool(ctx, oraclePool)

	// Preparation complete! Nothing can go wrong now (naive). It's time to collect raw request ids.
	if len(req.RawRequests) == 0 {
		return types.ErrEmptyRawRequests
	}
	// We now have everything we need to the request, so let's add it to the store.
	rid := k.AddRequest(ctx, req)
	// Emit an event describing a data request and asked validators.
	event := sdk.NewEvent(types.EventTypeRequest)
	event = event.AppendAttributes(
		sdk.NewAttribute(types.AttributeKeyID, fmt.Sprintf("%d", rid)),
		sdk.NewAttribute(types.AttributeKeyClientID, req.ClientID),
		sdk.NewAttribute(types.AttributeKeyOracleScriptID, fmt.Sprintf("%d", req.OracleScriptID)),
		sdk.NewAttribute(types.AttributeKeyCalldata, hex.EncodeToString(req.Calldata)),
		sdk.NewAttribute(types.AttributeKeyAskCount, fmt.Sprintf("%d", askCount)),
		sdk.NewAttribute(types.AttributeKeyMinCount, fmt.Sprintf("%d", req.MinCount)),
		sdk.NewAttribute(types.AttributeKeyGasUsed, fmt.Sprintf("%d", output.GasUsed)),
	)
	for _, val := range req.RequestedValidators {
		event = event.AppendAttributes(sdk.NewAttribute(types.AttributeKeyValidator, val.String()))
	}
	ctx.EventManager().EmitEvent(event)
	// Emit an event for each of the raw data requests
	for _, rawReq := range env.GetRawRequests() {
		ds, err := k.GetDataSource(ctx, rawReq.DataSourceID)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeRawRequest,
			sdk.NewAttribute(types.AttributeKeyDataSourceID, fmt.Sprintf("%d", rawReq.DataSourceID)),
			sdk.NewAttribute(types.AttributeKeyDataSourceHash, ds.Filename),
			sdk.NewAttribute(types.AttributeKeyExternalID, fmt.Sprintf("%d", rawReq.ExternalID)),
			sdk.NewAttribute(types.AttributeKeyCalldata, string(rawReq.Calldata)),
		))
	}
	return nil
}

// ResolveRequest resolves the given request and saves the result to the store. The function
// assumes that the given request is in a resolvable state with sufficient reporters.
func (k Keeper) ResolveRequest(ctx sdk.Context, reqID types.RequestID) {
	req := k.MustGetRequest(ctx, reqID)
	env := types.NewExecuteEnv(req, k.GetReports(ctx, reqID))
	script := k.MustGetOracleScript(ctx, req.OracleScriptID)
	code := k.GetFile(script.Filename)
	maxDataSize := k.GetParamUint64(ctx, types.KeyMaxDataSize)
	output, err := owasm.Execute(code, types.WasmExecuteGas, int64(maxDataSize), env)
	if err != nil {
		k.ResolveFailure(ctx, reqID, err.Error())
	} else if env.Retdata == nil {
		k.ResolveFailure(ctx, reqID, "no return data")
	} else {
		k.ResolveSuccess(ctx, reqID, env.Retdata, output.GasUsed)
	}
}

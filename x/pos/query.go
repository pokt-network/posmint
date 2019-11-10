package pos

import (
	"fmt"
	"github.com/pokt-network/posmint/client/context"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/types"
)

func QueryValidator(cdc *codec.Codec, addr sdk.ValAddress) (types.Validator, error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	res, _, err := cliCtx.QueryStore(types.KeyForValByAllVals(addr), types.StoreKey)
	if err != nil {
		return types.Validator{}, err
	}
	if len(res) == 0 {
		return types.Validator{}, fmt.Errorf("no validator found with address %s", addr)
	}
	return types.MustUnmarshalValidator(cdc, res), nil
}

func QueryValidators(cdc *codec.Codec) (types.Validators, error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	resKVs, _, err := cliCtx.QuerySubspace(types.AllValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
	}
	return validators, nil
}

func QueryStakedValidators(cdc *codec.Codec) (types.Validators, error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	resKVs, _, err := cliCtx.QuerySubspace(types.StakedValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
	}
	return validators, nil
}

func QueryUnstakedValidators(cdc *codec.Codec) (types.Validators, error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	resKVs, _, err := cliCtx.QuerySubspace(types.UnstakedValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
	}
	return validators, nil
}

func QueryUnstakingValidators(cdc *codec.Codec) (types.Validators, error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	resKVs, _, err := cliCtx.QuerySubspace(types.UnstakingValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
	}
	return validators, nil
}

func QuerySigningInfo(cdc *codec.Codec, ctx sdk.Context, consAddr sdk.ConsAddress) (types.ValidatorSigningInfo, error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	key := types.GetValidatorSigningInfoKey(consAddr)
	res, _, err := cliCtx.QueryStore(key, types.StoreKey)
	if err != nil {
		return types.ValidatorSigningInfo{}, err
	}
	if len(res) == 0 {
		return types.ValidatorSigningInfo{}, fmt.Errorf("validator %s not found in slashing store", consAddr)
	}
	var signingInfo types.ValidatorSigningInfo
	cdc.MustUnmarshalBinaryLengthPrefixed(res, &signingInfo)
	return types.ValidatorSigningInfo{}, cliCtx.PrintOutput(signingInfo)
}

func QuerySupply(cdc *codec.Codec) (stakedCoins sdk.Int, unstakedCoins sdk.Int, err error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	stakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/stakedPool", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	unstakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/unstakedPool", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	var stakedPool types.StakingPool
	if err := cdc.UnmarshalJSON(stakedPoolBytes, &stakedPool); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	var unstakedPool types.StakingPool
	if err := cdc.UnmarshalJSON(unstakedPoolBytes, &unstakedPool); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	return stakedPool.Tokens, unstakedPool.Tokens, nil
}

func QueryDAO(cdc *codec.Codec) (daoCoins sdk.Int, err error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	daoPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/dao", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, err
	}
	var daoPool types.DAOPool
	if err := cdc.UnmarshalJSON(daoPoolBytes, &daoPool); err != nil {
		return sdk.Int{}, err
	}
	return daoPool.Tokens, err
}

func QueryPOSParams(cdc *codec.Codec) (types.Params, error) {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}

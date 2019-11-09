package keeper

import (
	"fmt"
	"github.com/pokt-network/posmint/client/context"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/types"
	"github.com/tendermint/tendermint/crypto"
)

func (k Keeper) QueryValidator(addr sdk.ValAddress) (types.Validator, error) {
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	res, _, err := cliCtx.QueryStore(types.KeyForValByAllVals(addr), types.StoreKey)
	if err != nil {
		return types.Validator{}, err
	}
	if len(res) == 0 {
		return types.Validator{}, fmt.Errorf("no validator found with address %s", addr)
	}
	return types.MustUnmarshalValidator(k.cdc, res), nil
}

func (k Keeper) QueryValidators() (types.Validators, error) {
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	resKVs, _, err := cliCtx.QuerySubspace(types.AllValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(k.cdc, kv.Value))
	}
	return validators, nil
}

func (k Keeper) QueryStakedValidators() (types.Validators, error) {
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	resKVs, _, err := cliCtx.QuerySubspace(types.StakedValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(k.cdc, kv.Value))
	}
	return validators, nil
}

func (k Keeper) QueryUnstakedValidators() (types.Validators, error) {
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	resKVs, _, err := cliCtx.QuerySubspace(types.UnstakedValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(k.cdc, kv.Value))
	}
	return validators, nil
}

func (k Keeper) QueryUnstakingValidators() (types.Validators, error) {
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	resKVs, _, err := cliCtx.QuerySubspace(types.UnstakingValidatorsKey, types.StoreKey)
	if err != nil {
		return types.Validators{}, err
	}
	var validators types.Validators
	for _, kv := range resKVs {
		validators = append(validators, types.MustUnmarshalValidator(k.cdc, kv.Value))
	}
	return validators, nil
}

func (k Keeper) QuerySigningInfo(ctx sdk.Context, address sdk.ValAddress) (types.ValidatorSigningInfo, error) {
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	pk, err := k.getPubKeyRelation(ctx, crypto.Address(address))
	if err != nil {
		return types.ValidatorSigningInfo{}, err
	}
	consAddr := sdk.ConsAddress(pk.Address())
	key := types.GetValidatorSigningInfoKey(consAddr)
	res, _, err := cliCtx.QueryStore(key, types.StoreKey)
	if err != nil {
		return types.ValidatorSigningInfo{}, err
	}
	if len(res) == 0 {
		return types.ValidatorSigningInfo{}, fmt.Errorf("validator %s not found in slashing store", consAddr)
	}
	var signingInfo types.ValidatorSigningInfo
	k.cdc.MustUnmarshalBinaryLengthPrefixed(res, &signingInfo)
	return types.ValidatorSigningInfo{}, cliCtx.PrintOutput(signingInfo)
}

func (k Keeper) QuerySupply() (stakedCoins sdk.Int, unstakedCoins sdk.Int, err error) {
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	stakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/stakedPool", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	unstakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/unstakedPool", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	var stakedPool types.StakingPool
	if err := k.cdc.UnmarshalJSON(stakedPoolBytes, &stakedPool); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	var unstakedPool types.StakingPool
	if err := k.cdc.UnmarshalJSON(unstakedPoolBytes, &unstakedPool); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	return stakedPool.Tokens, unstakedPool.Tokens, nil
}

func (k Keeper) QueryDAO() (daoCoins sdk.Int, err error) {
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	daoPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/dao", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, err
	}
	var daoPool types.DAOPool
	if err := k.cdc.UnmarshalJSON(daoPoolBytes, &daoPool); err != nil {
		return sdk.Int{}, err
	}
	return daoPool.Tokens, err
}

func (k Keeper) QueryPOSParams() (types.Params, error) {
	cliCtx := context.NewCLIContext().WithCodec(k.cdc)
	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	k.cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}

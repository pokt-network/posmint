package staking

import (
	"fmt"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/exported"
	"github.com/pokt-network/posmint/x/pos/types"
)

// InitGenesis sets the pool and parameters for the provided keeper.  For each
// validator in data, it sets that validator in the keeper along with manually
// setting the indexes. In addition, it also sets any delegations found in
// data. Finally, it updates the staked validators.
// Returns final validator set after applying all declaration and delegations
func InitGenesis(ctx sdk.Context, keeper Keeper, accountKeeper types.AccountKeeper,
	supplyKeeper types.SupplyKeeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	stakedTokens := sdk.ZeroInt()

	// We need to pretend to be "n blocks before genesis", where "n" is the
	// validator update delay, so that e.g. slashing periods are correctly
	// initialized for the validator set e.g. with a one-block offset - the
	// first TM block is at height 1, so state updates applied from
	// genesis.json are in block 0.
	ctx = ctx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)
	// set the parameters from the data
	keeper.SetParams(ctx, data.Params)
	// set the 'previous state total power' from the data
	keeper.SetPrevStateValidatorsPower(ctx, data.PrevStateTotalPower)

	for _, validator := range data.Validators {
		// Call the creation hook if not exported
		if !data.Exported {
			keeper.BeforeValidatorCreated(ctx, validator.Address)
		}
		// set the validators from the data
		keeper.SetValidator(ctx, validator)

		// Manually set indices for the first time
		keeper.SetValidatorByConsAddr(ctx, validator)
		keeper.SetStakedValidator(ctx, validator)

		// Call the creation hook if not exported
		if !data.Exported {
			keeper.AddPubKeyRelation(ctx, validator.GetConsPubKey())
			keeper.AfterValidatorCreated(ctx, validator.Address)
		}

		// update unstaking validators if necessary
		if validator.IsUnstaking() {
			keeper.SetUnstakingValidator(ctx, validator)
		}

		if validator.IsStaked() {
			stakedTokens = stakedTokens.Add(validator.GetTokens())
		}
	}

	stakedCoins := sdk.NewCoins(sdk.NewCoin(data.Params.StakeDenom, stakedTokens))

	// check if the staked pool accounts exists
	stakedPool := keeper.GetStakedPool(ctx)
	if stakedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.StakedPoolName))
	}
	// check if the dao pool account exists
	daoPool := keeper.GetDAOPool(ctx)
	if daoPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.DAOPoolName))
	}

	// add coins if not provided on genesis
	if stakedPool.GetCoins().IsZero() {
		if err := stakedPool.SetCoins(stakedCoins); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, stakedPool)
	} else {
		if !stakedPool.GetCoins().IsEqual(stakedCoins) {
			panic(fmt.Sprintf("%s module account total does not equal the amount in each validator account", types.StakedPoolName))
		}
	}

	// don't need to run Tendermint updates if we exported
	if data.Exported {
		for _, lv := range data.PrevStateValidatorPowers {
			// set the staked validator powers from the previous state
			keeper.SetPrevStateValPower(ctx, lv.Address, lv.Power)
			validator, found := keeper.GetValidator(ctx, lv.Address)
			if !found {
				panic(fmt.Sprintf("validator %s not found", lv.Address))
			}
			update := validator.ABCIValidatorUpdate()
			update.Power = lv.Power // keep the next-val-set offset, use the prevState power for the first block
			res = append(res, update)
		}
	} else {
		// run tendermint updates
		res = keeper.UpdateTendermintValidators(ctx)
	}

	// slashing init genesis below todo

	keeper.IterateAndExecuteOverVals(ctx,
		func(index int64, validator exported.ValidatorI) bool {
			keeper.AddPubKeyRelation(ctx, validator.GetConsPubKey())
			return false
		},
	)

	for addr, info := range data.SigningInfos {
		address, err := sdk.ConsAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		keeper.SetValidatorSigningInfo(ctx, address, info)
	}

	for addr, array := range data.MissedBlocks {
		address, err := sdk.ConsAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		for _, missed := range array {
			keeper.SetValidatorMissedBlockBitArray(ctx, address, missed.Index, missed.Missed)
		}
	}

	keeper.Paramstore.SetParamSet(ctx, &data.Params)

	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, validators, and stakes found in
// the keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	prevStateTotalPower := keeper.PrevStateValidatorsPower(ctx)
	validators := keeper.GetAllValidators(ctx)
	var prevStateValidatorPowers []types.PrevStateBlockValidatorPower
	keeper.IterateAndExecuteOverPrevStateValsByPower(ctx, func(addr sdk.ValAddress, power int64) (stop bool) {
		prevStateValidatorPowers = append(prevStateValidatorPowers, types.PrevStateBlockValidatorPower{Address: addr, Power: power})
		return false
	})

	return types.GenesisState{
		Params:                   params,
		PrevStateTotalPower:      prevStateTotalPower,
		PrevStateValidatorPowers: prevStateValidatorPowers,
		Validators:               validators,
		Exported:                 true,
	}
}

// WriteValidators returns a slice of staked genesis validators.
func WriteValidators(ctx sdk.Context, keeper Keeper) (vals []tmtypes.GenesisValidator) {
	keeper.IterateAndExecuteOverPrevStateVals(ctx, func(_ int64, validator exported.ValidatorI) (stop bool) {
		vals = append(vals, tmtypes.GenesisValidator{
			PubKey: validator.GetConsPubKey(),
			Power:  validator.GetConsensusPower(),
			Name:   validator.GetAddress().String(),
		})

		return false
	})

	return
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data types.GenesisState) error {
	err := validateGenesisStateValidators(data.Validators)
	if err != nil {
		return err
	}
	err = data.Params.Validate()
	if err != nil {
		return err
	}
	downtime := data.Params.SlashFractionDowntime
	if downtime.IsNegative() || downtime.GT(sdk.OneDec()) {
		return fmt.Errorf("Slashing fraction downtime should be less than or equal to one and greater than zero, is %s", downtime.String())
	}

	dblSign := data.Params.SlashFractionDoubleSign
	if dblSign.IsNegative() || dblSign.GT(sdk.OneDec()) {
		return fmt.Errorf("Slashing fraction double sign should be less than or equal to one and greater than zero, is %s", dblSign.String())
	}

	minSign := data.Params.MinSignedPerWindow
	if minSign.IsNegative() || minSign.GT(sdk.OneDec()) {
		return fmt.Errorf("Min signed per window should be less than or equal to one and greater than zero, is %s", minSign.String())
	}

	maxEvidence := data.Params.MaxEvidenceAge
	if maxEvidence < 1*time.Minute {
		return fmt.Errorf("Max evidence age must be at least 1 minute, is %s", maxEvidence.String())
	}

	downtimeJail := data.Params.DowntimeJailDuration
	if downtimeJail < 1*time.Minute {
		return fmt.Errorf("Downtime unblond duration must be at least 1 minute, is %s", downtimeJail.String())
	}

	signedWindow := data.Params.SignedBlocksWindow
	if signedWindow < 10 {
		return fmt.Errorf("Signed blocks window must be at least 10, is %d", signedWindow)
	}

	return nil
}

func validateGenesisStateValidators(validators []types.Validator) (err error) {
	addrMap := make(map[string]bool, len(validators))
	for i := 0; i < len(validators); i++ {
		val := validators[i]
		strKey := string(val.ConsPubKey.Bytes())
		if _, ok := addrMap[strKey]; ok {
			return fmt.Errorf("duplicate validator in genesis state: address %v", val.ConsAddress())
		}
		if val.Jailed && val.IsStaked() {
			return fmt.Errorf("validator is staked and jailed in genesis state: address %v", val.ConsAddress())
		}
		if val.StakedTokens.IsZero() && !val.IsUnstaked() {
			return fmt.Errorf("staked/unstaked genesis validator cannot have zero stake, validator: %v", val)
		}
		// todo validate the minimum stake and the validators
		addrMap[strKey] = true
	}
	return
}

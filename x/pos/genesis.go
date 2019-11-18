package pos

import (
	"fmt"
	"github.com/pokt-network/posmint/x/pos/keeper"
	"time"

	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/pos/exported"
	"github.com/pokt-network/posmint/x/pos/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis sets up the module based on the genesis state
// First TM block is at height 1, so state updates applied from
// genesis.json are in block 0.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, supplyKeeper types.SupplyKeeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	stakedTokens := sdk.ZeroInt()
	ctx = ctx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)
	// set the parameters from the data
	keeper.SetParams(ctx, data.Params)
	// set the 'previous state total power' from the data
	keeper.SetPrevStateValidatorsPower(ctx, data.PrevStateTotalPower)
	for _, validator := range data.Validators {
		if validator.IsUnstaked() {
			panic(fmt.Sprintf("%v the validators must be staked or unstaking at genesis", validator))
		}
		// Call the registration hook if not exported
		if !data.Exported {
			keeper.BeforeValidatorRegistered(ctx, validator.Address)
		}
		// set the validators from the data
		keeper.SetValidator(ctx, validator)
		keeper.SetValidatorByConsAddr(ctx, validator)
		keeper.SetStakedValidator(ctx, validator)
		// Call the creation hook if not exported
		if !data.Exported {
			keeper.AfterValidatorRegistered(ctx, validator.Address)
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
	if daoPool.GetCoins().IsZero() {
		if err := daoPool.SetCoins(sdk.NewCoins(sdk.NewCoin(data.Params.StakeDenom, data.DAO.Tokens))); err != nil {
			panic(err)
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
	// add public key relationship to address
	keeper.IterateAndExecuteOverVals(ctx,
		func(index int64, validator exported.ValidatorI) bool {
			keeper.AddPubKeyRelation(ctx, validator.GetConsPubKey())
			return false
		},
	)
	// update signing information from genesis state
	for addr, info := range data.SigningInfos {
		address, err := sdk.ConsAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		keeper.SetValidatorSigningInfo(ctx, address, info)
	}
	// update missed block information from genesis state
	for addr, array := range data.MissedBlocks {
		address, err := sdk.ConsAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		for _, missed := range array {
			keeper.SetMissedBlockArray(ctx, address, missed.Index, missed.Missed)
		}
	}
	// set the params set in the keeper
	keeper.Paramstore.SetParamSet(ctx, &data.Params)
	if data.PreviousProposer != nil {
		keeper.SetPreviousProposer(ctx, data.PreviousProposer)
	}
	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	prevStateTotalPower := keeper.PrevStateValidatorsPower(ctx)
	validators := keeper.GetAllValidators(ctx)
	var prevStateValidatorPowers []types.PrevStatePowerMapping
	keeper.IterateAndExecuteOverPrevStateValsByPower(ctx, func(addr sdk.ValAddress, power int64) (stop bool) {
		prevStateValidatorPowers = append(prevStateValidatorPowers, types.PrevStatePowerMapping{Address: addr, Power: power})
		return false
	})
	signingInfos := make(map[string]types.ValidatorSigningInfo)
	missedBlocks := make(map[string][]types.MissedBlock)
	keeper.IterateAndExecuteOverValSigningInfo(ctx, func(address sdk.ConsAddress, info types.ValidatorSigningInfo) (stop bool) {
		bechAddr := address.String()
		signingInfos[bechAddr] = info
		localMissedBlocks := []types.MissedBlock{}

		keeper.IterateAndExecuteOverMissedArray(ctx, address, func(index int64, missed bool) (stop bool) {
			localMissedBlocks = append(localMissedBlocks, types.MissedBlock{index, missed})
			return false
		})
		missedBlocks[bechAddr] = localMissedBlocks

		return false
	})
	daoTokens := keeper.GetDAOTokens(ctx)
	daoPool := types.DAOPool{Tokens: daoTokens}
	prevProposer := keeper.GetPreviousProposer(ctx)

	return types.GenesisState{
		Params:                   params,
		PrevStateTotalPower:      prevStateTotalPower,
		PrevStateValidatorPowers: prevStateValidatorPowers,
		Validators:               validators,
		Exported:                 true,
		DAO:                      daoPool,
		SigningInfos:             signingInfos,
		MissedBlocks:             missedBlocks,
		PreviousProposer:         prevProposer,
	}
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data types.GenesisState) error {
	err := validateGenesisStateValidators(data.Validators, sdk.NewInt(data.Params.StakeMinimum))
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

func validateGenesisStateValidators(validators []types.Validator, minimumStake sdk.Int) (err error) {
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
		addrMap[strKey] = true
		if !val.IsUnstaked() && val.StakedTokens.LTE(minimumStake) {
			return fmt.Errorf("validator has less than minimum stake: %v", val)
		}
	}
	return
}

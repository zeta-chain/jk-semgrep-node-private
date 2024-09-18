package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	// ModuleName defines the module name
	ModuleName                       = "emissions"
	UndistributedObserverRewardsPool = ModuleName + "Observers"
	UndistributedTssRewardsPool      = ModuleName + "Tss"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey              = "mem_emissions"
	WithdrawableEmissionsKey = "WithdrawableEmissions-value-"
	ParamsKey                = "Params-value-"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	EmissionsTrackerKey              = "EmissionsTracker-value-"
	ParamValidatorEmissionPercentage = "ValidatorEmissionPercentage"
	ParamObserverEmissionPercentage  = "ObserverEmissionPercentage"
	ParamTssSignerEmissionPercentage = "SignerEmissionPercentage"
	ParamObserverSlashAmount         = "ObserverSlashAmount"
)

var (
	EmissionsModuleAddress                  = authtypes.NewModuleAddress(ModuleName)
	UndistributedObserverRewardsPoolAddress = authtypes.NewModuleAddress(UndistributedObserverRewardsPool)
	UndistributedTssRewardsPoolAddress      = authtypes.NewModuleAddress(UndistributedTssRewardsPool)
	// BlockReward is an initial block reward amount when emissions module was initialized.
	// The current value can be obtained from by querying the params
	BlockReward = sdk.MustNewDecFromStr("9620949074074074074.074070733466756687")
	// ObserverSlashAmount is the amount of tokens to be slashed from observer in case of incorrect vote
	// by default it is set to 0.1 ZETA
	ObserverSlashAmount = sdkmath.NewInt(100000000000000000)

	// BallotMaturityBlocks is amount of blocks needed for ballot to mature
	// by default is set to 100
	BallotMaturityBlocks = 100
)

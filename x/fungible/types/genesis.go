package types

import (
	"fmt"
)

// DefaultGenesis returns the default fungible genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ForeignCoinsList: []ForeignCoins{},
		SystemContract:   nil,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in foreignCoins
	foreignCoinsIndexMap := make(map[string]struct{})

	for _, elem := range gs.ForeignCoinsList {
		index := string(ForeignCoinsKey(elem.Zrc20ContractAddress))
		if _, ok := foreignCoinsIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for foreignCoins")
		}
		foreignCoinsIndexMap[index] = struct{}{}
	}

	return nil
}

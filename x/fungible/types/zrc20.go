package types

import ethcommon "github.com/ethereum/go-ethereum/common"

// DefaultLiquidityCap is the default value set for the liquidity cap of a new ZRC20 when deployed
// for security reason, this value is low. An arbitrary value should be set during the process of deploying a new ZRC20
// The value is represented in the base unit of the ZRC20, final value is calculated by multiplying this value by 10^decimals
const DefaultLiquidityCap = uint64(1000)

// ZRC20Data represents the ZRC4 token details used to map
// the token to a Cosmos Coin
type ZRC20Data struct {
	Name     string
	Symbol   string
	Decimals uint8
}

// ZRC20StringResponse defines the string value from the call response
type ZRC20StringResponse struct {
	Value string
}

// ZRC20Uint8Response defines the uint8 value from the call response
type ZRC20Uint8Response struct {
	Value uint8
}

// ZRC20BoolResponse defines the bool value from the call response
type ZRC20BoolResponse struct {
	Value bool
}

// UniswapV2FactoryByte32Response defines the string value from the call response
type UniswapV2FactoryByte32Response struct {
	Value [32]byte
}

// SystemAddressResponse defines the address value from the call response
type SystemAddressResponse struct {
	Value ethcommon.Address
}

// NewZRC20Data creates a new ZRC20Data instance
func NewZRC20Data(name, symbol string, decimals uint8) ZRC20Data {
	return ZRC20Data{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}
}

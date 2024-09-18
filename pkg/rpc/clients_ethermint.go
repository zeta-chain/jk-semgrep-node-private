package rpc

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	feemarkettypes "github.com/zeta-chain/ethermint/x/feemarket/types"
)

// GetBaseGasPrice returns the base gas price
func (c *Clients) GetBaseGasPrice(ctx context.Context) (int64, error) {
	resp, err := c.EthermintFeeMarket.Params(ctx, &feemarkettypes.QueryParamsRequest{})
	if err != nil {
		return 0, errors.Wrap(err, "failed to get base gas price")
	}

	if resp.Params.BaseFee.IsNil() {
		return 0, fmt.Errorf("base fee is nil")
	}

	return resp.Params.BaseFee.Int64(), nil
}

package runner

import (
	"fmt"
	"math/big"
	"time"

	"github.com/cenkalti/backoff/v4"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/evm/zetaconnector.eth.sol"
	connectorzevm "github.com/zeta-chain/protocol-contracts/v1/pkg/contracts/zevm/zetaconnectorzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/retry"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func (r *E2ERunner) WaitForBlocks(n int64) {
	height, err := r.CctxClient.LastZetaHeight(r.Ctx, &types.QueryLastZetaHeightRequest{})
	if err != nil {
		return
	}
	call := func() error {
		return retry.Retry(r.waitForBlock(height.Height + n))
	}

	bo := backoff.NewConstantBackOff(time.Second * 5)
	boWithMaxRetries := backoff.WithMaxRetries(bo, 10)
	err = retry.DoWithBackoff(call, boWithMaxRetries)
	require.NoError(r, err, "failed to wait for %d blocks", n)
}
func (r *E2ERunner) waitForBlock(n int64) error {
	height, err := r.CctxClient.LastZetaHeight(r.Ctx, &types.QueryLastZetaHeightRequest{})
	if err != nil {
		return err
	}
	if height.Height < n {
		return fmt.Errorf("waiting for %d blocks, current height %d", n, height.Height)
	}
	return nil
}

// WaitForTxReceiptOnZEVM waits for a tx receipt on ZEVM
func (r *E2ERunner) WaitForTxReceiptOnZEVM(tx *ethtypes.Transaction) {
	r.Lock()
	defer r.Unlock()

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)
}

// WaitForMinedCCTX waits for a cctx to be mined from a tx
func (r *E2ERunner) WaitForMinedCCTX(txHash ethcommon.Hash) {
	r.Lock()
	defer r.Unlock()

	cctx := utils.WaitCctxMinedByInboundHash(
		r.Ctx,
		txHash.Hex(),
		r.CctxClient,
		r.Logger,
		r.CctxTimeout,
	)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)
}

// WaitForMinedCCTXFromIndex waits for a cctx to be mined from its index
func (r *E2ERunner) WaitForMinedCCTXFromIndex(index string) {
	r.Lock()
	defer r.Unlock()

	cctx := utils.WaitCCTXMinedByIndex(r.Ctx, index, r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)
}

// SendZetaOnEvm sends ZETA to an address on EVM
// this allows the ZETA contract deployer to funds other accounts on EVM
func (r *E2ERunner) SendZetaOnEvm(address ethcommon.Address, zetaAmount int64) *ethtypes.Transaction {
	// the deployer might be sending ZETA in different goroutines
	r.Lock()
	defer r.Unlock()

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(zetaAmount))
	tx, err := r.ZetaEth.Transfer(r.EVMAuth, address, amount)
	require.NoError(r, err)

	return tx
}

// DepositZeta deposits ZETA on ZetaChain from the ZETA smart contract on EVM
func (r *E2ERunner) DepositZeta() ethcommon.Hash {
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Zeta

	return r.DepositZetaWithAmount(r.EVMAddress(), amount)
}

// DepositZetaWithAmount deposits ZETA on ZetaChain from the ZETA smart contract on EVM with the specified amount
func (r *E2ERunner) DepositZetaWithAmount(to ethcommon.Address, amount *big.Int) ethcommon.Hash {
	tx, err := r.ZetaEth.Approve(r.EVMAuth, r.ConnectorEthAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "approve")
	r.requireTxSuccessful(receipt, "approve tx failed")

	// query the chain ID using zevm client
	zetaChainID, err := r.ZEVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	tx, err = r.ConnectorEth.Send(r.EVMAuth, zetaconnectoreth.ZetaInterfacesSendInput{
		// TODO: allow user to specify destination chain id
		// https://github.com/zeta-chain/node-private/issues/41
		DestinationChainId:  zetaChainID,
		DestinationAddress:  to.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("Send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "send")
	r.requireTxSuccessful(receipt, "send tx failed")

	r.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := r.ConnectorEth.ParseZetaSent(*log)
		if err == nil {
			r.Logger.Info("    Connector: %s", r.ConnectorEthAddr.String())
			r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			r.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			r.Logger.Info("    Block Num: %d", log.BlockNumber)
		}
	}

	return tx.Hash()
}

// DepositAndApproveWZeta deposits and approves WZETA on ZetaChain from the ZETA smart contract on ZEVM
func (r *E2ERunner) DepositAndApproveWZeta(amount *big.Int) {
	r.ZEVMAuth.Value = amount
	tx, err := r.WZeta.Deposit(r.ZEVMAuth)
	require.NoError(r, err)

	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta deposit")
	r.requireTxSuccessful(receipt, "deposit failed")

	tx, err = r.WZeta.Approve(r.ZEVMAuth, r.ConnectorZEVMAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("wzeta approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta approve")
	r.requireTxSuccessful(receipt, "approve failed, logs: %+v", receipt.Logs)
}

// WithdrawZeta withdraws ZETA from ZetaChain to the ZETA smart contract on EVM
// waitReceipt specifies whether to wait for the tx receipt and check if the tx was successful
func (r *E2ERunner) WithdrawZeta(amount *big.Int, waitReceipt bool) *ethtypes.Transaction {
	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	tx, err := r.ConnectorZEVM.Send(r.ZEVMAuth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  r.EVMAddress().Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     amount,
		ZetaParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("send tx hash: %s", tx.Hash().Hex())

	if waitReceipt {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.Logger.EVMReceipt(*receipt, "send")
		r.requireTxSuccessful(receipt, "send failed, logs: %+v", receipt.Logs)

		r.Logger.Info("  Logs:")
		for _, log := range receipt.Logs {
			sentLog, err := r.ConnectorZEVM.ParseZetaSent(*log)
			if err == nil {
				r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
				r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
				r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
				r.Logger.Info("    Zeta Value: %d", sentLog.ZetaValueAndGas)
			}
		}
	}

	return tx
}

// WithdrawEther withdraws Ether from ZetaChain to the ZETA smart contract on EVM
func (r *E2ERunner) WithdrawEther(amount *big.Int) *ethtypes.Transaction {
	// withdraw
	tx, err := r.ETHZRC20.Withdraw(r.ZEVMAuth, r.EVMAddress().Bytes(), amount)
	require.NoError(r, err)

	r.Logger.EVMTransaction(*tx, "withdraw")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, "withdraw failed")

	r.Logger.EVMReceipt(*receipt, "withdraw")
	r.Logger.ZRC20Withdrawal(r.ETHZRC20, *receipt, "withdraw")

	return tx
}

// WithdrawERC20 withdraws an ERC20 token from ZetaChain to the ZETA smart contract on EVM
func (r *E2ERunner) WithdrawERC20(amount *big.Int) *ethtypes.Transaction {
	tx, err := r.ERC20ZRC20.Withdraw(r.ZEVMAuth, r.EVMAddress().Bytes(), amount)
	require.NoError(r, err)

	r.Logger.EVMTransaction(*tx, "withdraw")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := r.ERC20ZRC20.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		r.Logger.Info(
			"  logs: from %s, to %x, value %d, gasfee %d",
			event.From.Hex(),
			event.To,
			event.Value,
			event.GasFee,
		)
	}

	return tx
}

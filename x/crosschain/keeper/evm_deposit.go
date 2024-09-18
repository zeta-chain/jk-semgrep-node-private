package keeper

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

const InCCTXIndexKey = "inCctxIndex"

// HandleEVMDeposit handles a deposit from an inbound tx
// returns (isContractReverted, err)
// (true, non-nil) means CallEVM() reverted
func (k Keeper) HandleEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx) (bool, error) {
	to := ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
	sender := ethcommon.HexToAddress(cctx.InboundParams.Sender)
	var ethTxHash ethcommon.Hash
	inboundAmount := cctx.GetInboundParams().Amount.BigInt()
	inboundSender := cctx.GetInboundParams().Sender
	inboundSenderChainID := cctx.GetInboundParams().SenderChainId
	inboundCoinType := cctx.InboundParams.CoinType

	if len(ctx.TxBytes()) > 0 {
		// add event for tendermint transaction hash format
		hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
		ethTxHash = ethcommon.BytesToHash(hash)
		cctx.GetCurrentOutboundParam().Hash = ethTxHash.String()
		// #nosec G115 always positive
		cctx.GetCurrentOutboundParam().ObservedExternalHeight = uint64(ctx.BlockHeight())
	}

	if inboundCoinType == coin.CoinType_Zeta {
		// In case of an error
		// 	- Return true will revert the cctx and create a revert cctx with status PendingRevert
		// 	- Return false will abort the cctx
		indexBytes, err := cctx.GetCCTXIndexBytes()
		if err != nil {
			return false, err
		}
		data, err := base64.StdEncoding.DecodeString(cctx.RelayedMessage)
		if err != nil {
			return true, errors.Wrap(types.ErrUnableToDecodeMessageString, err.Error())
		}
		// if coin type is Zeta, this is a deposit ZETA to zEVM cctx.
		evmTxResponse, err := k.fungibleKeeper.ZETADepositAndCallContract(
			ctx,
			sender,
			to,
			inboundSenderChainID,
			inboundAmount,
			data,
			indexBytes,
		)
		if fungibletypes.IsContractReverted(evmTxResponse, err) || errShouldRevertCctx(err) {
			return true, err
		} else if err != nil {
			return false, err
		}
	} else {
		var (
			message []byte
			err     error
		)

		// in protocol version 1, the destination of the deposit is the first 20 bytes of the message when the message is not empty
		// in protocol version 2, the destination of the deposit is always the to address, the message is the data to be sent to the contract
		if cctx.ProtocolContractVersion == types.ProtocolContractVersion_V1 {
			var parsedAddress ethcommon.Address
			parsedAddress, message, err = chains.ParseAddressAndData(cctx.RelayedMessage)
			if err != nil {
				return false, errors.Wrap(types.ErrUnableToParseAddress, err.Error())
			}
			if parsedAddress != (ethcommon.Address{}) {
				to = parsedAddress
			}
		} else if cctx.ProtocolContractVersion == types.ProtocolContractVersion_V2 {
			if len(cctx.RelayedMessage) > 0 {
				message, err = hex.DecodeString(cctx.RelayedMessage)
				if err != nil {
					return false, errors.Wrap(types.ErrUnableToDecodeMessageString, err.Error())
				}
			}
		}

		from, err := chains.DecodeAddressFromChainID(inboundSenderChainID, inboundSender, k.GetAuthorityKeeper().GetAdditionalChainList(ctx))
		if err != nil {
			return false, fmt.Errorf("HandleEVMDeposit: unable to decode address: %s", err.Error())
		}

		evmTxResponse, contractCall, err := k.fungibleKeeper.ZRC20DepositAndCallContract(
			ctx,
			from,
			to,
			inboundAmount,
			inboundSenderChainID,
			message,
			inboundCoinType,
			cctx.InboundParams.Asset,
			cctx.ProtocolContractVersion,
		)
		if fungibletypes.IsContractReverted(evmTxResponse, err) || errShouldRevertCctx(err) {
			return true, err
		} else if err != nil {
			return false, err
		}

		// non-empty msg.Message means this is a contract call; therefore the logs should be processed.
		// a withdrawal event in the logs could generate cctxs for outbound transactions.
		if !evmTxResponse.Failed() && contractCall {
			logs := evmtypes.LogsToEthereum(evmTxResponse.Logs)
			if len(logs) > 0 {
				ctx = ctx.WithValue(InCCTXIndexKey, cctx.Index)
				txOrigin := cctx.InboundParams.TxOrigin
				if txOrigin == "" {
					txOrigin = inboundSender
				}

				err = k.ProcessLogs(ctx, logs, to, txOrigin)
				if err != nil {
					// ProcessLogs should not error; error indicates exception, should abort
					return false, errors.Wrap(types.ErrCannotProcessWithdrawal, err.Error())
				}
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(sdk.EventTypeMessage,
						sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
						sdk.NewAttribute("action", "DepositZRC20AndCallContract"),
						sdk.NewAttribute("contract", to.String()),
						sdk.NewAttribute("data", hex.EncodeToString(message)),
						sdk.NewAttribute("cctxIndex", cctx.Index),
					),
				)
			}
		}
	}
	return false, nil
}

// errShouldRevertCctx returns true if the cctx should revert from the error of the deposit
// we revert the cctx if a non-contract is tried to be called, if the liquidity cap is reached, or if the zrc20 is paused
func errShouldRevertCctx(err error) bool {
	return errors.Is(err, fungibletypes.ErrForeignCoinCapReached) ||
		errors.Is(err, fungibletypes.ErrCallNonContract) ||
		errors.Is(err, fungibletypes.ErrPausedZRC20)
}

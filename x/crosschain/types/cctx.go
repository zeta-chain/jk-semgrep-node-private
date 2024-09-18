package types

import (
	"encoding/hex"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"

	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// GetEVMRevertAddress returns the EVM revert address
// If a revert address is specified in the revert options, it returns the address
// Otherwise returns sender address
func (m CrossChainTx) GetEVMRevertAddress() ethcommon.Address {
	addr, valid := m.RevertOptions.GetEVMRevertAddress()
	if valid {
		return addr
	}
	return ethcommon.HexToAddress(m.InboundParams.Sender)
}

// GetEVMAbortAddress returns the EVM abort address
// If an abort address is specified in the revert options, it returns the address
// Otherwise returns sender address
func (m CrossChainTx) GetEVMAbortAddress() ethcommon.Address {
	addr, valid := m.RevertOptions.GetEVMAbortAddress()
	if valid {
		return addr
	}
	return ethcommon.HexToAddress(m.InboundParams.Sender)
}

// GetCurrentOutboundParam returns the current outbound params.
// There can only be one active outbound.
// OutboundParams[0] is the original outbound, if it reverts, then
// OutboundParams[1] is the new outbound.
func (m CrossChainTx) GetCurrentOutboundParam() *OutboundParams {
	if len(m.OutboundParams) == 0 {
		return &OutboundParams{}
	}
	return m.OutboundParams[len(m.OutboundParams)-1]
}

// IsCurrentOutboundRevert returns true if the current outbound is the revert tx.
func (m CrossChainTx) IsCurrentOutboundRevert() bool {
	return len(m.OutboundParams) >= 2
}

// OriginalDestinationChainID returns the original destination of the outbound tx, reverted or not
// If there is no outbound tx, return -1
func (m CrossChainTx) OriginalDestinationChainID() int64 {
	if len(m.OutboundParams) == 0 {
		return -1
	}
	return m.OutboundParams[0].ReceiverChainId
}

// Validate checks if the CCTX is valid.
func (m CrossChainTx) Validate() error {
	if m.InboundParams == nil {
		return fmt.Errorf("inbound tx params cannot be nil")
	}
	if m.OutboundParams == nil {
		return fmt.Errorf("outbound tx params cannot be nil")
	}
	if m.CctxStatus == nil {
		return fmt.Errorf("cctx status cannot be nil")
	}
	if len(m.OutboundParams) > 2 {
		return fmt.Errorf("outbound tx params cannot be more than 2")
	}
	if m.Index != "" {
		err := ValidateCCTXIndex(m.Index)
		if err != nil {
			return err
		}
	}
	err := m.InboundParams.Validate()
	if err != nil {
		return err
	}
	for _, outboundParam := range m.OutboundParams {
		err = outboundParam.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// AddRevertOutbound does the following things in one function:
//  1. create a new OutboundTxParams for the revert
//  2. append the new OutboundTxParams to the current OutboundTxParams
//  3. update the TxFinalizationStatus of the current OutboundTxParams to Executed.
func (m *CrossChainTx) AddRevertOutbound(gasLimit uint64) error {
	if m.IsCurrentOutboundRevert() {
		return fmt.Errorf("cannot revert a revert tx")
	}
	if len(m.OutboundParams) == 0 {
		return fmt.Errorf("cannot revert before trying to process an outbound tx")
	}

	// in protocol contract V2, developers can specify a specific address to receive the revert
	// if not specified, the sender address is used
	// note: this option is current only support for EVM type chains
	revertReceiver := m.InboundParams.Sender
	if m.ProtocolContractVersion == ProtocolContractVersion_V2 {
		revertAddress, valid := m.RevertOptions.GetEVMRevertAddress()
		if valid {
			revertReceiver = revertAddress.Hex()
		}
	}

	revertTxParams := &OutboundParams{
		Receiver:        revertReceiver,
		ReceiverChainId: m.InboundParams.SenderChainId,
		Amount:          m.GetCurrentOutboundParam().Amount,
		GasLimit:        gasLimit,
		TssPubkey:       m.GetCurrentOutboundParam().TssPubkey,
	}
	// The original outbound has been finalized, the new outbound is pending
	m.GetCurrentOutboundParam().TxFinalizationStatus = TxFinalizationStatus_Executed
	m.OutboundParams = append(m.OutboundParams, revertTxParams)
	return nil
}

// AddOutbound adds a new outbound tx to the CCTX.
func (m *CrossChainTx) AddOutbound(
	ctx sdk.Context,
	msg MsgVoteOutbound,
	ballotStatus observertypes.BallotStatus,
) error {
	if ballotStatus != observertypes.BallotStatus_BallotFinalized_FailureObservation {
		if !msg.ValueReceived.Equal(m.GetCurrentOutboundParam().Amount) {
			ctx.Logger().Error(fmt.Sprintf("VoteOutbound: Mint mismatch: %s value received vs %s cctx amount",
				msg.ValueReceived,
				m.GetCurrentOutboundParam().Amount))
			return cosmoserrors.Wrap(
				sdkerrors.ErrInvalidRequest,
				fmt.Sprintf(
					"ValueReceived %s does not match sent value %s",
					msg.ValueReceived,
					m.GetCurrentOutboundParam().Amount,
				),
			)
		}
	}
	// Update CCTX values
	m.GetCurrentOutboundParam().Hash = msg.ObservedOutboundHash
	m.GetCurrentOutboundParam().GasUsed = msg.ObservedOutboundGasUsed
	m.GetCurrentOutboundParam().EffectiveGasPrice = msg.ObservedOutboundEffectiveGasPrice
	m.GetCurrentOutboundParam().EffectiveGasLimit = msg.ObservedOutboundEffectiveGasLimit
	m.GetCurrentOutboundParam().ObservedExternalHeight = msg.ObservedOutboundBlockHeight
	return nil
}

// SetAbort sets the CCTX status to Aborted with the given error message.
func (m CrossChainTx) SetAbort(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_Aborted, message)
}

// SetPendingRevert sets the CCTX status to PendingRevert with the given error message.
func (m CrossChainTx) SetPendingRevert(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_PendingRevert, message)
}

// SetPendingOutbound sets the CCTX status to PendingOutbound with the given error message.
func (m CrossChainTx) SetPendingOutbound(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_PendingOutbound, message)
}

// SetOutboundMined sets the CCTX status to OutboundMined with the given error message.
func (m CrossChainTx) SetOutboundMined(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_OutboundMined, message)
}

// SetReverted sets the CCTX status to Reverted with the given error message.
func (m CrossChainTx) SetReverted(message string) {
	m.CctxStatus.ChangeStatus(CctxStatus_Reverted, message)
}

func (m CrossChainTx) GetCCTXIndexBytes() ([32]byte, error) {
	sendHash := [32]byte{}
	if len(m.Index) < 2 {
		return [32]byte{}, fmt.Errorf("decode CCTX %s index too short", m.Index)
	}
	decodedIndex, err := hex.DecodeString(m.Index[2:]) // remove the leading 0x
	if err != nil || len(decodedIndex) != 32 {
		return [32]byte{}, fmt.Errorf("decode CCTX %s error", m.Index)
	}
	copy(sendHash[:32], decodedIndex[:32])
	return sendHash, nil
}

// SetOutboundBallotIndex sets the outbound ballot index for the most recent outbound.
func (m CrossChainTx) SetOutboundBallotIndex(index string) {
	m.GetCurrentOutboundParam().BallotIndex = index
}

// GetCctxIndexFromBytes returns the CCTX index from a byte array.
func GetCctxIndexFromBytes(sendHash [32]byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(sendHash[:]))
}

// NewCCTX creates a new CCTX from a MsgVoteInbound message and a TSS pubkey.
// It also validates the created cctx
func NewCCTX(ctx sdk.Context, msg MsgVoteInbound, tssPubkey string) (CrossChainTx, error) {
	index := msg.Digest()

	if msg.TxOrigin == "" {
		msg.TxOrigin = msg.Sender
	}
	inboundParams := &InboundParams{
		Sender:                 msg.Sender,
		SenderChainId:          msg.SenderChainId,
		TxOrigin:               msg.TxOrigin,
		Asset:                  msg.Asset,
		Amount:                 msg.Amount,
		ObservedHash:           msg.InboundHash,
		ObservedExternalHeight: msg.InboundBlockHeight,
		FinalizedZetaHeight:    0,
		BallotIndex:            index,
		CoinType:               msg.CoinType,
	}

	outboundParams := &OutboundParams{
		Receiver:               msg.Receiver,
		ReceiverChainId:        msg.ReceiverChain,
		Hash:                   "",
		TssNonce:               0,
		GasLimit:               msg.GasLimit,
		GasPrice:               "",
		GasPriorityFee:         "",
		BallotIndex:            "",
		ObservedExternalHeight: 0,
		Amount:                 sdkmath.ZeroUint(),
		TssPubkey:              tssPubkey,
		CoinType:               msg.CoinType,
	}
	status := &Status{
		Status:              CctxStatus_PendingInbound,
		StatusMessage:       "",
		CreatedTimestamp:    ctx.BlockHeader().Time.Unix(),
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
		IsAbortRefunded:     false,
	}
	cctx := CrossChainTx{
		Creator:                 msg.Creator,
		Index:                   index,
		ZetaFees:                sdkmath.ZeroUint(),
		RelayedMessage:          msg.Message,
		CctxStatus:              status,
		InboundParams:           inboundParams,
		OutboundParams:          []*OutboundParams{outboundParams},
		ProtocolContractVersion: msg.ProtocolContractVersion,
		RevertOptions:           msg.RevertOptions,
	}

	// TODO: remove this validate call
	// https://github.com/zeta-chain/node/issues/2236
	err := cctx.Validate()
	if err != nil {
		return CrossChainTx{}, err
	}

	return cctx, nil
}

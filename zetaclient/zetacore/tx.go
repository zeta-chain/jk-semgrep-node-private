package zetacore

import (
	"context"
	"strings"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/proofs"
	"github.com/zeta-chain/node/x/crosschain/types"
	clientauthz "github.com/zeta-chain/node/zetaclient/authz"
	clientcommon "github.com/zeta-chain/node/zetaclient/common"
)

// GetInboundVoteMessage returns a new MsgVoteInbound
// TODO(revamp): move to a different file
func GetInboundVoteMessage(
	sender string,
	senderChain int64,
	txOrigin string,
	receiver string,
	receiverChain int64,
	amount math.Uint,
	message string,
	inboundHash string,
	inBlockHeight uint64,
	gasLimit uint64,
	coinType coin.CoinType,
	asset string,
	signerAddress string,
	eventIndex uint,
) *types.MsgVoteInbound {
	msg := types.NewMsgVoteInbound(
		signerAddress,
		sender,
		senderChain,
		txOrigin,
		receiver,
		receiverChain,
		amount,
		message,
		inboundHash,
		inBlockHeight,
		gasLimit,
		coinType,
		asset,
		eventIndex,
		types.ProtocolContractVersion_V1,
	)
	return msg
}

// GasPriceMultiplier returns the gas price multiplier for the given chain
func GasPriceMultiplier(chain chains.Chain) float64 {
	switch chain.Consensus {
	case chains.Consensus_ethereum:
		return clientcommon.EVMOutboundGasPriceMultiplier
	case chains.Consensus_bitcoin:
		return clientcommon.BTCOutboundGasPriceMultiplier
	default:
		return clientcommon.DefaultGasPriceMultiplier
	}
}

// WrapMessageWithAuthz wraps a message with an authz message
// used since a hotkey is used to broadcast the transactions, instead of the operator
func WrapMessageWithAuthz(msg sdk.Msg) (sdk.Msg, clientauthz.Signer, error) {
	msgURL := sdk.MsgTypeURL(msg)

	// verify message validity
	if err := msg.ValidateBasic(); err != nil {
		return nil, clientauthz.Signer{}, errors.Wrapf(err, "invalid message %q", msgURL)
	}

	authzSigner := clientauthz.GetSigner(msgURL)
	authzMessage := authz.NewMsgExec(authzSigner.GranteeAddress, []sdk.Msg{msg})

	return &authzMessage, authzSigner, nil
}

// AddOutboundTracker adds an outbound tracker
// TODO(revamp): rename to PostAddOutboundTracker
func (c *Client) AddOutboundTracker(
	ctx context.Context,
	chainID int64,
	nonce uint64,
	txHash string,
	proof *proofs.Proof,
	blockHash string,
	txIndex int64,
) (string, error) {
	// don't report if the tracker already contains the txHash
	tracker, err := c.GetOutboundTracker(ctx, chains.Chain{ChainId: chainID}, nonce)
	if err == nil {
		for _, hash := range tracker.HashList {
			if strings.EqualFold(hash.TxHash, txHash) {
				return "", nil
			}
		}
	}

	signerAddress := c.keys.GetOperatorAddress().String()
	msg := types.NewMsgAddOutboundTracker(signerAddress, chainID, nonce, txHash, proof, blockHash, txIndex)

	authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
	if err != nil {
		return "", err
	}

	zetaTxHash, err := c.Broadcast(ctx, AddOutboundTrackerGasLimit, authzMsg, authzSigner)
	if err != nil {
		return "", err
	}

	return zetaTxHash, nil
}

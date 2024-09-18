package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// GetAllAuthzZetaclientTxTypes returns all the authz types for required for zetaclient
func GetAllAuthzZetaclientTxTypes() []string {
	return []string{
		sdk.MsgTypeURL(&MsgVoteGasPrice{}),
		sdk.MsgTypeURL(&MsgVoteInbound{}),
		sdk.MsgTypeURL(&MsgVoteOutbound{}),
		sdk.MsgTypeURL(&MsgAddOutboundTracker{}),
		sdk.MsgTypeURL(&observertypes.MsgVoteTSS{}),
		sdk.MsgTypeURL(&observertypes.MsgVoteBlame{}),
		sdk.MsgTypeURL(&observertypes.MsgVoteBlockHeader{}),
	}
}

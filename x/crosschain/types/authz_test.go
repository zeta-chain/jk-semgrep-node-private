package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestGetAllAuthzZetaclientTxTypes(t *testing.T) {
	require.Equal(t, []string{"/zetachain.zetacore.crosschain.MsgVoteGasPrice",
		"/zetachain.zetacore.crosschain.MsgVoteInbound",
		"/zetachain.zetacore.crosschain.MsgVoteOutbound",
		"/zetachain.zetacore.crosschain.MsgAddOutboundTracker",
		"/zetachain.zetacore.observer.MsgVoteTSS",
		"/zetachain.zetacore.observer.MsgVoteBlame",
		"/zetachain.zetacore.observer.MsgVoteBlockHeader"},
		crosschaintypes.GetAllAuthzZetaclientTxTypes())
}

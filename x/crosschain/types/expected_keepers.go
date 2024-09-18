package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/proofs"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

type StakingKeeper interface {
	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) types.ModuleAccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

type ObserverKeeper interface {
	GetObserverSet(ctx sdk.Context) (val observertypes.ObserverSet, found bool)
	GetBallot(ctx sdk.Context, index string) (val observertypes.Ballot, found bool)
	GetChainParamsByChainID(ctx sdk.Context, chainID int64) (params *observertypes.ChainParams, found bool)
	GetNodeAccount(ctx sdk.Context, address string) (nodeAccount observertypes.NodeAccount, found bool)
	GetAllNodeAccount(ctx sdk.Context) (nodeAccounts []observertypes.NodeAccount)
	SetNodeAccount(ctx sdk.Context, nodeAccount observertypes.NodeAccount)
	IsInboundEnabled(ctx sdk.Context) (found bool)
	GetCrosschainFlags(ctx sdk.Context) (val observertypes.CrosschainFlags, found bool)
	GetKeygen(ctx sdk.Context) (val observertypes.Keygen, found bool)
	SetKeygen(ctx sdk.Context, keygen observertypes.Keygen)
	SetCrosschainFlags(ctx sdk.Context, crosschainFlags observertypes.CrosschainFlags)
	SetLastObserverCount(ctx sdk.Context, lbc *observertypes.LastObserverCount)
	AddVoteToBallot(
		ctx sdk.Context,
		ballot observertypes.Ballot,
		address string,
		observationType observertypes.VoteType,
	) (observertypes.Ballot, error)
	CheckIfFinalizingVote(ctx sdk.Context, ballot observertypes.Ballot) (observertypes.Ballot, bool)
	IsNonTombstonedObserver(ctx sdk.Context, address string) bool
	FindBallot(
		ctx sdk.Context,
		index string,
		chain chains.Chain,
		observationType observertypes.ObservationType,
	) (ballot observertypes.Ballot, isNew bool, err error)
	AddBallotToList(ctx sdk.Context, ballot observertypes.Ballot)
	CheckIfTssPubkeyHasBeenGenerated(ctx sdk.Context, tssPubkey string) (observertypes.TSS, bool)
	GetAllTSS(ctx sdk.Context) (list []observertypes.TSS)
	GetTSS(ctx sdk.Context) (val observertypes.TSS, found bool)
	SetTSS(ctx sdk.Context, tss observertypes.TSS)
	SetTSSHistory(ctx sdk.Context, tss observertypes.TSS)
	GetTssAddress(
		goCtx context.Context,
		req *observertypes.QueryGetTssAddressRequest,
	) (*observertypes.QueryGetTssAddressResponse, error)

	SetFundMigrator(ctx sdk.Context, fm observertypes.TssFundMigratorInfo)
	GetFundMigrator(ctx sdk.Context, chainID int64) (val observertypes.TssFundMigratorInfo, found bool)
	GetAllTssFundMigrators(ctx sdk.Context) (fms []observertypes.TssFundMigratorInfo)
	RemoveAllExistingMigrators(ctx sdk.Context)
	SetChainNonces(ctx sdk.Context, chainNonces observertypes.ChainNonces)
	GetChainNonces(ctx sdk.Context, chainID int64) (val observertypes.ChainNonces, found bool)
	GetAllChainNonces(ctx sdk.Context) (list []observertypes.ChainNonces)
	SetNonceToCctx(ctx sdk.Context, nonceToCctx observertypes.NonceToCctx)
	GetNonceToCctx(ctx sdk.Context, tss string, chainID int64, nonce int64) (val observertypes.NonceToCctx, found bool)
	GetAllPendingNonces(ctx sdk.Context) (list []observertypes.PendingNonces, err error)
	GetPendingNonces(ctx sdk.Context, tss string, chainID int64) (val observertypes.PendingNonces, found bool)
	SetPendingNonces(ctx sdk.Context, pendingNonces observertypes.PendingNonces)
	SetTssAndUpdateNonce(ctx sdk.Context, tss observertypes.TSS)
	RemoveFromPendingNonces(ctx sdk.Context, tss string, chainID int64, nonce int64)
	GetAllNonceToCctx(ctx sdk.Context) (list []observertypes.NonceToCctx)
	VoteOnInboundBallot(
		ctx sdk.Context,
		senderChainID int64,
		receiverChainID int64,
		coinType coin.CoinType,
		voter string,
		ballotIndex string,
		inboundHash string,
	) (bool, bool, error)
	VoteOnOutboundBallot(
		ctx sdk.Context,
		ballotIndex string,
		outTxChainID int64,
		receiveStatus chains.ReceiveStatus,
		voter string,
	) (bool, bool, observertypes.Ballot, string, error)
	GetSupportedChainFromChainID(ctx sdk.Context, chainID int64) (chains.Chain, bool)
	GetSupportedChains(ctx sdk.Context) []chains.Chain
}

type FungibleKeeper interface {
	GetForeignCoins(ctx sdk.Context, zrc20Addr string) (val fungibletypes.ForeignCoins, found bool)
	GetAllForeignCoins(ctx sdk.Context) (list []fungibletypes.ForeignCoins)
	GetAllForeignCoinMap(ctx sdk.Context) map[int64]map[string]fungibletypes.ForeignCoins
	SetForeignCoins(ctx sdk.Context, foreignCoins fungibletypes.ForeignCoins)
	GetAllForeignCoinsForChain(ctx sdk.Context, foreignChainID int64) (list []fungibletypes.ForeignCoins)
	GetForeignCoinFromAsset(ctx sdk.Context, asset string, chainID int64) (fungibletypes.ForeignCoins, bool)
	GetGasCoinForForeignCoin(ctx sdk.Context, chainID int64) (fungibletypes.ForeignCoins, bool)
	GetSystemContract(ctx sdk.Context) (val fungibletypes.SystemContract, found bool)
	QuerySystemContractGasCoinZRC20(ctx sdk.Context, chainID *big.Int) (ethcommon.Address, error)
	GetUniswapV2Router02Address(ctx sdk.Context) (ethcommon.Address, error)
	QueryUniswapV2RouterGetZetaAmountsIn(
		ctx sdk.Context,
		amountOut *big.Int,
		outZRC4 ethcommon.Address,
	) (*big.Int, error)
	QueryUniswapV2RouterGetZRC4ToZRC4AmountsIn(
		ctx sdk.Context,
		amountOut *big.Int,
		inZRC4, outZRC4 ethcommon.Address,
	) (*big.Int, error)
	QueryGasLimit(ctx sdk.Context, contract ethcommon.Address) (*big.Int, error)
	QueryProtocolFlatFee(ctx sdk.Context, contract ethcommon.Address) (*big.Int, error)
	SetGasPrice(ctx sdk.Context, chainID *big.Int, gasPrice *big.Int) (uint64, error)
	DepositCoinZeta(ctx sdk.Context, to ethcommon.Address, amount *big.Int) error
	DepositZRC20(
		ctx sdk.Context,
		contract ethcommon.Address,
		to ethcommon.Address,
		amount *big.Int,
	) (*evmtypes.MsgEthereumTxResponse, error)
	ZRC20DepositAndCallContract(
		ctx sdk.Context,
		from []byte,
		to ethcommon.Address,
		amount *big.Int,
		senderChainID int64,
		data []byte,
		coinType coin.CoinType,
		asset string,
		protocolContractVersion ProtocolContractVersion,
	) (*evmtypes.MsgEthereumTxResponse, bool, error)
	ProcessV2RevertDeposit(
		ctx sdk.Context,
		amount *big.Int,
		chainID int64,
		coinType coin.CoinType,
		asset string,
		revertAddress ethcommon.Address,
		callOnRevert bool,
		revertMessage []byte,
	) error
	CallUniswapV2RouterSwapExactTokensForTokens(
		ctx sdk.Context,
		sender ethcommon.Address,
		to ethcommon.Address,
		amountIn *big.Int,
		inZRC4,
		outZRC4 ethcommon.Address,
		noEthereumTxEvent bool,
	) (ret []*big.Int, err error)
	CallUniswapV2RouterSwapExactETHForToken(
		ctx sdk.Context,
		sender ethcommon.Address,
		to ethcommon.Address,
		amountIn *big.Int,
		outZRC4 ethcommon.Address,
		noEthereumTxEvent bool,
	) ([]*big.Int, error)
	CallZRC20Burn(
		ctx sdk.Context,
		sender ethcommon.Address,
		zrc20address ethcommon.Address,
		amount *big.Int,
		noEthereumTxEvent bool,
	) error
	CallZRC20Approve(
		ctx sdk.Context,
		owner ethcommon.Address,
		zrc20address ethcommon.Address,
		spender ethcommon.Address,
		amount *big.Int,
		noEthereumTxEvent bool,
	) error
	DeployZRC20Contract(
		ctx sdk.Context,
		name, symbol string,
		decimals uint8,
		chainID int64,
		coinType coin.CoinType,
		erc20Contract string,
		gasLimit *big.Int,
	) (ethcommon.Address, error)
	FundGasStabilityPool(ctx sdk.Context, chainID int64, amount *big.Int) error
	WithdrawFromGasStabilityPool(ctx sdk.Context, chainID int64, amount *big.Int) error
	ZETADepositAndCallContract(ctx sdk.Context,
		sender ethcommon.Address,
		to ethcommon.Address,
		inboundSenderChainID int64,
		inboundAmount *big.Int,
		data []byte,
		indexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error)
	ZETARevertAndCallContract(ctx sdk.Context,
		sender ethcommon.Address,
		to ethcommon.Address,
		inboundSenderChainID int64,
		destinationChainID int64,
		remainingAmount *big.Int,
		data []byte,
		indexBytes [32]byte) (*evmtypes.MsgEthereumTxResponse, error)
}

type AuthorityKeeper interface {
	CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error
	GetAdditionalChainList(ctx sdk.Context) (list []chains.Chain)
}

type LightclientKeeper interface {
	VerifyProof(ctx sdk.Context, proof *proofs.Proof, chainID int64, blockHash string, txIndex int64) ([]byte, error)
}

type IBCCrosschainKeeper interface {
}

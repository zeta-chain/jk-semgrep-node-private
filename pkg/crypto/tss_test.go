package crypto

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/cosmos"
)

func TestGetTssAddrEVM(t *testing.T) {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	pk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	require.NoError(t, err)
	decompresspubkey, err := crypto.DecompressPubkey(pubKey.Bytes())
	require.NoError(t, err)
	testCases := []struct {
		name      string
		tssPubkey string
		wantAddr  ethcommon.Address
		wantErr   bool
	}{
		{
			name:      "Valid TSS pubkey",
			tssPubkey: pk,
			wantAddr:  crypto.PubkeyToAddress(*decompresspubkey),
			wantErr:   false,
		},
		{
			name:      "Invalid TSS pubkey",
			tssPubkey: "invalid",
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := GetTssAddrEVM(tc.tssPubkey)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.Equal(t, tc.wantAddr, addr)
				require.NoError(t, err)
				require.NotEmpty(t, addr)
			}
		})
	}
}

func TestGetTssAddrBTC(t *testing.T) {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	pk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	require.NoError(t, err)
	testCases := []struct {
		name          string
		tssPubkey     string
		bitcoinParams *chaincfg.Params
		wantErr       bool
	}{
		{
			name:          "Valid TSS pubkey testnet params",
			tssPubkey:     pk,
			bitcoinParams: &chaincfg.TestNet3Params,
			wantErr:       false,
		},
		{
			name:          "Invalid TSS pubkey testnet params",
			tssPubkey:     "invalid",
			bitcoinParams: &chaincfg.TestNet3Params,
			wantErr:       true,
		},
		{
			name:          "Valid TSS pubkey mainnet params",
			tssPubkey:     pk,
			bitcoinParams: &chaincfg.MainNetParams,
			wantErr:       false,
		},
		{
			name:          "Invalid TSS pubkey mainnet params",
			tssPubkey:     "invalid",
			bitcoinParams: &chaincfg.MainNetParams,
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := GetTssAddrBTC(tc.tssPubkey, tc.bitcoinParams)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				expectedAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.Bytes()), tc.bitcoinParams)
				require.NoError(t, err)
				require.NotEmpty(t, addr)
				require.Equal(t, expectedAddr.EncodeAddress(), addr)
			}
		})
	}
}

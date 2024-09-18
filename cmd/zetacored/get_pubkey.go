package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/cosmos"
	zetacrypto "github.com/zeta-chain/node/pkg/crypto"
)

func GetPubKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-pubkey [tssKeyName] [password]",
		Short: "Get the node account public key",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			tssKeyName := args[0]
			password := ""
			if len(args) > 1 {
				password = args[1]
			}
			pubKeySet, err := GetPubKeySet(clientCtx, tssKeyName, password)
			if err != nil {
				return err
			}
			fmt.Println(pubKeySet.String())
			return nil
		},
	}
	return cmd
}

func GetPubKeySet(clientctx client.Context, tssAccountName, password string) (zetacrypto.PubKeySet, error) {
	pubkeySet := zetacrypto.PubKeySet{
		Secp256k1: "",
		Ed25519:   "",
	}
	//kb, err := GetKeyringKeybase(keyringPath, tssAccountName, password)
	privKeyArmor, err := clientctx.Keyring.ExportPrivKeyArmor(tssAccountName, password)
	if err != nil {
		return pubkeySet, err
	}
	priKey, _, err := crypto.UnarmorDecryptPrivKey(privKeyArmor, password)
	if err != nil {
		return pubkeySet, fmt.Errorf("fail to unarmor private key: %w", err)
	}

	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, priKey.PubKey())
	if err != nil {
		return pubkeySet, err
	}
	pubkey, err := zetacrypto.NewPubKey(s)
	if err != nil {
		return pubkeySet, err
	}
	pubkeySet.Secp256k1 = pubkey
	return pubkeySet, nil
}

package crypto

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/stretchr/testify/require"
	. "gopkg.in/check.v1"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/cosmos"
)

type KeyDataAddr struct {
	mainnet string
	testnet string
	mocknet string
}

type KeyData struct {
	priv    string
	pub     string
	addrETH KeyDataAddr
}

type PubKeyTestSuite struct {
	keyData []KeyData
}

var _ = Suite(&PubKeyTestSuite{})

func (s *PubKeyTestSuite) SetUpSuite(_ *C) {
	s.keyData = []KeyData{
		{
			priv: "ef235aacf90d9f4aadd8c92e4b2562e1d9eb97f0df9ba3b508258739cb013db2",
			pub:  "02b4632d08485ff1df2db55b9dafd23347d1c47a457072a1e87be26896549a8737",
			addrETH: KeyDataAddr{
				mainnet: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
				testnet: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
				mocknet: "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
			},
		},
		{
			priv: "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032",
			pub:  "037db227d7094ce215c3a0f57e1bcc732551fe351f94249471934567e0f5dc1bf7",
			addrETH: KeyDataAddr{
				mainnet: "0x970e8128ab834e8eac17ab8e3812f010678cf791",
				testnet: "0x970e8128ab834e8eac17ab8e3812f010678cf791",
				mocknet: "0x970e8128ab834e8eac17ab8e3812f010678cf791",
			},
		},
		{
			priv: "e810f1d7d6691b4a7a73476f3543bd87d601f9a53e7faf670eac2c5b517d83bf",
			pub:  "03f98464e8d3fc8e275e34c6f8dc9b99aa244e37b0d695d0dfb8884712ed6d4d35",
			addrETH: KeyDataAddr{
				mainnet: "0xf6da288748ec4c77642f6c5543717539b3ae001b",
				testnet: "0xf6da288748ec4c77642f6c5543717539b3ae001b",
				mocknet: "0xf6da288748ec4c77642f6c5543717539b3ae001b",
			},
		},
		{
			priv: "a96e62ed3955e65be32703f12d87b6b5cf26039ecfa948dc5107a495418e5330",
			pub:  "02950e1cdfcb133d6024109fd489f734eeb4502418e538c28481f22bce276f248c",
			addrETH: KeyDataAddr{
				mainnet: "0xfabb9cc6ec839b1214bb11c53377a56a6ed81762",
				testnet: "0xfabb9cc6ec839b1214bb11c53377a56a6ed81762",
				mocknet: "0xfabb9cc6ec839b1214bb11c53377a56a6ed81762",
			},
		},
		{
			priv: "9294f4d108465fd293f7fe299e6923ef71a77f2cb1eb6d4394839c64ec25d5c0",
			pub:  "0238383ee4d60176d27cf46f0863bfc6aea624fe9bfc7f4273cc5136d9eb483e4a",
			addrETH: KeyDataAddr{
				mainnet: "0x1f30a82340f08177aba70e6f48054917c74d7d38",
				testnet: "0x1f30a82340f08177aba70e6f48054917c74d7d38",
				mocknet: "0x1f30a82340f08177aba70e6f48054917c74d7d38",
			},
		},
	}
}

// TestPubKey implementation
func (s *PubKeyTestSuite) TestPubKey(c *C) {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	spk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	c.Assert(err, IsNil)
	pk, err := NewPubKey(spk)
	c.Assert(err, IsNil)
	hexStr := pk.String()
	c.Assert(len(hexStr) > 0, Equals, true)
	pk1, err := NewPubKey(hexStr)
	c.Assert(err, IsNil)
	c.Assert(pk.Equals(pk1), Equals, true)

	result, err := json.Marshal(pk)
	c.Assert(err, IsNil)
	c.Log(result, Equals, fmt.Sprintf(`"%s"`, hexStr))
	var pk2 PubKey
	err = json.Unmarshal(result, &pk2)
	c.Assert(err, IsNil)
	c.Assert(pk2.Equals(pk), Equals, true)
}

func (s *PubKeyTestSuite) TestPubKeyGetAddress(c *C) {
	original := os.Getenv("NET")
	defer func() {
		c.Assert(os.Setenv("NET", original), IsNil)
	}()

	for _, d := range s.keyData {
		privB, _ := hex.DecodeString(d.priv)
		pubB, _ := hex.DecodeString(d.pub)
		priv := secp256k1.PrivKey(privB)
		pubKey := priv.PubKey()
		pubT, _ := pubKey.(secp256k1.PubKey)
		pub := pubT[:]

		c.Assert(hex.EncodeToString(pub), Equals, hex.EncodeToString(pubB))

		tmp, err := codec.FromTmPubKeyInterface(pubKey)
		c.Assert(err, IsNil)
		pubBech32, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, tmp)
		c.Assert(err, IsNil)

		pk, err := NewPubKey(pubBech32)
		c.Assert(err, IsNil)

		c.Assert(os.Setenv("NET", "mainnet"), IsNil)
		addrETH, err := pk.GetAddress(chains.Goerli)
		c.Assert(err, IsNil)
		c.Assert(addrETH.String(), Equals, d.addrETH.mainnet)

		c.Assert(os.Setenv("NET", "testnet"), IsNil)
		addrETH, err = pk.GetAddress(chains.Goerli)
		c.Assert(err, IsNil)
		c.Assert(addrETH.String(), Equals, d.addrETH.testnet)

		c.Assert(os.Setenv("NET", "mocknet"), IsNil)
		addrETH, err = pk.GetAddress(chains.Goerli)
		c.Assert(err, IsNil)
		c.Assert(addrETH.String(), Equals, d.addrETH.mocknet)

		addrETH, err = pk.GetAddress(chains.BitcoinRegtest)
		c.Assert(err, IsNil)
		c.Assert(addrETH, Equals, chains.NoAddress)
	}
}

func (s *PubKeyTestSuite) TestEquals(c *C) {
	var pk1, pk2, pk3, pk4 PubKey
	_, pubKey1, _ := testdata.KeyTestPubAddr()
	tpk1, err1 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey1)
	c.Assert(err1, IsNil)
	pk1 = PubKey(tpk1)

	_, pubKey2, _ := testdata.KeyTestPubAddr()
	tpk2, err2 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey2)
	c.Assert(err2, IsNil)
	pk2 = PubKey(tpk2)

	_, pubKey3, _ := testdata.KeyTestPubAddr()
	tpk3, err3 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey3)
	c.Assert(err3, IsNil)
	pk3 = PubKey(tpk3)

	_, pubKey4, _ := testdata.KeyTestPubAddr()
	tpk4, err4 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey4)
	c.Assert(err4, IsNil)
	pk4 = PubKey(tpk4)

	c.Assert(PubKeys{
		pk1, pk2,
	}.Equals(nil), Equals, false)

	c.Assert(PubKeys{
		pk1, pk2, pk3,
	}.Equals(PubKeys{
		pk1, pk2,
	}), Equals, false)

	c.Assert(PubKeys{
		pk1, pk2, pk3, pk4,
	}.Equals(PubKeys{
		pk4, pk3, pk2, pk1,
	}), Equals, true)

	c.Assert(PubKeys{
		pk1, pk2, pk3, pk4,
	}.Equals(PubKeys{
		pk1, pk2, pk3, pk4,
	}), Equals, true)

	c.Assert(PubKeys{
		pk1, pk2, pk3,
	}.Equals(PubKeys{
		pk1, pk2, pk4,
	}), Equals, false)
}

func TestNewPubKey(t *testing.T) {
	t.Run("should create new pub key from string", func(t *testing.T) {
		_, pubKey, _ := testdata.KeyTestPubAddr()
		spk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
		require.NoError(t, err)
		pk, err := NewPubKey(spk)
		require.NoError(t, err)
		require.Equal(t, PubKey(spk), pk)
	})

	t.Run("should return empty pub key from empty string", func(t *testing.T) {
		pk, err := NewPubKey("")
		require.NoError(t, err)
		require.Equal(t, EmptyPubKey, pk)
	})

	t.Run("should fail if not bech32 encoded string", func(t *testing.T) {
		pk, err := NewPubKey("invalid")
		require.Error(t, err)
		require.Equal(t, EmptyPubKey, pk)
	})
}

func TestGetAddressFromPubkeyString(t *testing.T) {
	t.Run("should get address from pubkey string", func(t *testing.T) {
		_, pubKey, _ := testdata.KeyTestPubAddr()
		spk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
		require.NoError(t, err)
		_, err = GetAddressFromPubkeyString(spk)
		require.NoError(t, err)
	})

	t.Run("should get address from nonbech32 string", func(t *testing.T) {
		_, err := GetAddressFromPubkeyString("invalid")
		require.Error(t, err)
	})
}

func TestPubKeys(t *testing.T) {
	t.Run("should valid if all are valid", func(t *testing.T) {
		_, pubKey1, _ := testdata.KeyTestPubAddr()
		spk1, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey1)
		pk1, _ := NewPubKey(spk1)

		_, pubKey2, _ := testdata.KeyTestPubAddr()
		spk2, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey2)
		pk2, _ := NewPubKey(spk2)

		_, pubKey3, _ := testdata.KeyTestPubAddr()
		spk3, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey3)
		pk3, _ := NewPubKey(spk3)

		pubKeys := PubKeys{
			pk1, pk2, pk3,
		}

		require.Nil(t, pubKeys.Valid())
	})

	t.Run("should invalid if one is invalid", func(t *testing.T) {
		_, pubKey1, _ := testdata.KeyTestPubAddr()
		spk1, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey1)
		pk1, _ := NewPubKey(spk1)

		_, pubKey2, _ := testdata.KeyTestPubAddr()
		spk2, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey2)
		pk2, _ := NewPubKey(spk2)

		_, pubKey3, _ := testdata.KeyTestPubAddr()
		spk3, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey3)
		pk3, _ := NewPubKey(spk3)

		pubKeys := PubKeys{
			pk1, pk2, pk3, PubKey("invalid"),
		}

		require.NotNil(t, pubKeys.Valid())
	})

	t.Run("contains", func(t *testing.T) {
		_, pubKey1, _ := testdata.KeyTestPubAddr()
		spk1, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey1)
		pk1, _ := NewPubKey(spk1)

		_, pubKey2, _ := testdata.KeyTestPubAddr()
		spk2, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey2)
		pk2, _ := NewPubKey(spk2)

		_, pubKey3, _ := testdata.KeyTestPubAddr()
		spk3, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey3)
		pk3, _ := NewPubKey(spk3)

		pubKeys := PubKeys{
			pk1, pk2,
		}

		require.True(t, pubKeys.Contains(pk1))
		require.True(t, pubKeys.Contains(pk2))
		require.False(t, pubKeys.Contains(pk3))
	})

	t.Run("should concat string", func(t *testing.T) {
		_, pubKey1, _ := testdata.KeyTestPubAddr()
		spk1, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey1)
		pk1, _ := NewPubKey(spk1)

		_, pubKey2, _ := testdata.KeyTestPubAddr()
		spk2, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey2)
		pk2, _ := NewPubKey(spk2)

		_, pubKey3, _ := testdata.KeyTestPubAddr()
		spk3, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey3)
		pk3, _ := NewPubKey(spk3)

		pubKeys := PubKeys{
			pk1, pk2, pk3,
		}

		require.Equal(t, fmt.Sprintf("%s, %s, %s", pk1.String(), pk2.String(), pk3.String()), pubKeys.String())
	})

	t.Run("should return strings array", func(t *testing.T) {
		_, pubKey1, _ := testdata.KeyTestPubAddr()
		spk1, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey1)
		pk1, _ := NewPubKey(spk1)

		_, pubKey2, _ := testdata.KeyTestPubAddr()
		spk2, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey2)
		pk2, _ := NewPubKey(spk2)

		_, pubKey3, _ := testdata.KeyTestPubAddr()
		spk3, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey3)
		pk3, _ := NewPubKey(spk3)

		pubKeys := PubKeys{
			pk1, pk2, pk3,
		}

		require.Equal(t, []string{pk1.String(), pk2.String(), pk3.String()}, pubKeys.Strings())
	})
}

func TestGetEVMAddress(t *testing.T) {
	t.Run("should return empty if pubkey is empty", func(t *testing.T) {
		pubKey := PubKey("")
		e, err := pubKey.GetEVMAddress()
		require.NoError(t, err)
		require.Equal(t, chains.NoAddress, e)
	})

	t.Run("should return addr from pubkey", func(t *testing.T) {
		_, pubKey, _ := testdata.KeyTestPubAddr()
		spk, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
		pk, _ := NewPubKey(spk)

		_, err := pk.GetEVMAddress()
		require.NoError(t, err)
	})

	t.Run("should error if non bech32", func(t *testing.T) {
		pk := PubKey("invalid")
		e, err := pk.GetEVMAddress()
		require.ErrorContains(t, err, "decoding bech32 failed")
		require.Equal(t, chains.NoAddress, e)
	})
}

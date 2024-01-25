package xdr

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeAccount(t *testing.T, hexKey string) string {
	b, err := hex.DecodeString(hexKey)
	require.NoError(t, err)
	var key Uint256
	copy(key[:], b)
	a, err := NewAccountId(PublicKeyTypePublicKeyTypeEd25519, key)
	require.NoError(t, err)
	addr, err := a.GetAddress()
	require.NoError(t, err)
	return addr
}

func TestNewPoolId(t *testing.T) {
	testGetPoolID := func(x, y Asset, expectedHex string) {
		// TODO: Require x y to be sorted.
		// require.Less(t, x, y);

		id, err := NewPoolId(x, y, LiquidityPoolFeeV18)
		if assert.NoError(t, err) {
			idHex := hex.EncodeToString(id[:])
			assert.Equal(t, expectedHex, idHex)
		}
	}

	acc1 := makeAccount(t, "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	acc2 := makeAccount(t, "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789")

	t.Run("NATIVE and ALPHANUM4 (short and full length)", func(t *testing.T) {
		testGetPoolID(
			MustNewNativeAsset(), MustNewCreditAsset("AbC", acc1),
			"c17f36fbd210e43dca1cda8edc5b6c0f825fcb72b39f0392fd6309844d77ff7d")
		testGetPoolID(
			MustNewNativeAsset(), MustNewCreditAsset("AbCd", acc1),
			"80e0c5dc79ed76bb7e63681f6456136762f0d01ede94bb379dbc793e66db35e6")
	})

	t.Run("NATIVE and ALPHANUM12 (short and full length)", func(t *testing.T) {
		testGetPoolID(
			MustNewNativeAsset(), MustNewCreditAsset("AbCdEfGhIjK", acc1),
			"d2306c6e8532f99418e9d38520865e1c1059cddb6793da3cc634224f2ffb5bd4")
		testGetPoolID(
			MustNewNativeAsset(), MustNewCreditAsset("AbCdEfGhIjKl", acc1),
			"807e9e66653b5fda4dd4e672ff64a929fc5fdafe152eeadc07bb460c4849d711")
	})

	t.Run("ALPHANUM4 and ALPHANUM4 (short and full length)", func(t *testing.T) {
		testGetPoolID(
			MustNewCreditAsset("AbC", acc1), MustNewCreditAsset("aBc", acc1),
			"0239013ab016985fc3ab077d165a9b21b822efa013fdd422381659e76aec505b")
		testGetPoolID(
			MustNewCreditAsset("AbCd", acc1), MustNewCreditAsset("aBc", acc1),
			"cadb490d15b4333890377cd17400acf7681e14d6d949869ffa1fbbad7a6d2fde")
		testGetPoolID(
			MustNewCreditAsset("AbC", acc1), MustNewCreditAsset("aBcD", acc1),
			"a938f8f346f3aff41d2e03b05137ef1955a723861802a4042f51f0f816e0db36")
		testGetPoolID(
			MustNewCreditAsset("AbCd", acc1), MustNewCreditAsset("aBcD", acc1),
			"c89646bb6db726bfae784ab66041abbf54747cf4b6b16dff2a5c05830ad9c16b")
	})

	t.Run("ALPHANUM12 and ALPHANUM12 (short and full length)", func(t *testing.T) {
		testGetPoolID(
			MustNewCreditAsset("AbCdEfGhIjK", acc1), MustNewCreditAsset("aBcDeFgHiJk", acc1),
			"88dc054dd0f8146bac0e691095ce2b90cd902b499761d22b1c94df120ca0b060")
		testGetPoolID(
			MustNewCreditAsset("AbCdEfGhIjKl", acc1), MustNewCreditAsset("aBcDeFgHiJk", acc1),
			"09672910d891e658219d2f33a8885a542b2a5a09e9f486461201bd278a3e92a4")
		testGetPoolID(
			MustNewCreditAsset("AbCdEfGhIjK", acc1), MustNewCreditAsset("aBcDeFgHiJkl", acc1),
			"63501addf8a5a6522eac996226069190b5226c71cfdda22347022418af1948a0")
		testGetPoolID(
			MustNewCreditAsset("AbCdEfGhIjKl", acc1), MustNewCreditAsset("aBcDeFgHiJkl", acc1),
			"e851197a0148e949bdc03d52c53821b9afccc0fadfdc41ae01058c14c252e03b")
	})

	t.Run("ALPHANUM4 same code different issuer (short and full length)", func(t *testing.T) {
		testGetPoolID(
			MustNewCreditAsset("aBc", acc1), MustNewCreditAsset("aBc", acc2),
			"5d7188454299529856586e81ea385d2c131c6afdd9d58c82e9aa558c16522fea")
		testGetPoolID(
			MustNewCreditAsset("aBcD", acc1), MustNewCreditAsset("aBcD", acc2),
			"00d152f5f6b7e46eaf558576512207ea835a332f17ca777fba3cb835ef7dc1ef")
	})

	t.Run("ALPHANUM12 same code different issuer (short and full length)", func(t *testing.T) {
		testGetPoolID(
			MustNewCreditAsset("aBcDeFgHiJk", acc1), MustNewCreditAsset("aBcDeFgHiJk", acc2),
			"cad65154300f087e652981fa5f76aa469b43ad53e9a5d348f1f93da57193d022")
		testGetPoolID(
			MustNewCreditAsset("aBcDeFgHiJkL", acc1), MustNewCreditAsset("aBcDeFgHiJkL", acc2),
			"93fa82ecaabe987461d1e3c8e0fd6510558b86ac82a41f7c70b112281be90c71")
	})

	t.Run("ALPHANUM4 before ALPHANUM12 (short and full length) doesn't depend on issuer or code", func(t *testing.T) {
		testGetPoolID(
			MustNewCreditAsset("aBc", acc1), MustNewCreditAsset("aBcDeFgHiJk", acc2),
			"c0d4c87bbaade53764b904fde2901a0353af437e9d3a976f1252670b85a36895")
		testGetPoolID(
			MustNewCreditAsset("aBcD", acc1), MustNewCreditAsset("aBcDeFgHiJk", acc2),
			"1ee5aa0f0e6b8123c2da6592389481f64d816bfe3c3c06be282b0cdb0971f840")
		testGetPoolID(
			MustNewCreditAsset("aBc", acc1), MustNewCreditAsset("aBcDeFgHiJkL", acc2),
			"a87bc151b119c1ea289905f0cb3cf95be7b0f096a0b6685bf2dcae70f9515d53")
		testGetPoolID(
			MustNewCreditAsset("aBcD", acc1), MustNewCreditAsset("aBcDeFgHiJkL", acc2),
			"3caf78118d6cabd42618eef47bbc2da8abe7fe42539b4b502f08766485592a81")
		testGetPoolID(
			MustNewCreditAsset("aBc", acc2), MustNewCreditAsset("aBcDeFgHiJk", acc1),
			"befb7f966ae63adcfde6a6670478bb7d936c29849e25e3387bb9e74566e3a29f")
		testGetPoolID(
			MustNewCreditAsset("aBcD", acc2), MustNewCreditAsset("aBcDeFgHiJk", acc1),
			"593cc996c3f0d32e165fcbee9fdc5dba6ab05140a4a9254e08ad8cb67fe657a1")
		testGetPoolID(
			MustNewCreditAsset("aBc", acc2), MustNewCreditAsset("aBcDeFgHiJkL", acc1),
			"d66af9b7417547c3dc000617533405349d1f622015daf3e9bad703ea34ee1d17")
		testGetPoolID(
			MustNewCreditAsset("aBcD", acc2), MustNewCreditAsset("aBcDeFgHiJkL", acc1),
			"c1c7a4b9db6e3754cae3017f72b6b7c93198f593182c541bcab3795c6413a677")
	})
}

func TestNewPoolIdRejectsIncorrectOrder(t *testing.T) {
	acc1 := makeAccount(t, "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	_, err := NewPoolId(MustNewCreditAsset("AbC", acc1), MustNewNativeAsset(), LiquidityPoolFeeV18)
	assert.EqualError(t, err, "AssetA must be < AssetB")
}

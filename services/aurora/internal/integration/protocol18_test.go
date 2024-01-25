package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/protocols/aurora/effects"
	"github.com/diamcircle/go/protocols/aurora/operations"
	"github.com/diamcircle/go/services/aurora/internal/test/integration"
	"github.com/diamcircle/go/txnbuild"
	"github.com/diamcircle/go/xdr"
)

func NewProtocol18Test(t *testing.T) *integration.Test {
	config := integration.Config{ProtocolVersion: 18}
	return integration.NewTest(t, config)
}

func TestProtocol18Basics(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol18Test(t)
	master := itest.Master()

	t.Run("Sanity", func(t *testing.T) {
		root, err := itest.Client().Root()
		tt.NoError(err)
		tt.LessOrEqual(int32(18), root.CoreSupportedProtocolVersion)
		tt.Equal(int32(18), root.CurrentProtocolVersion)

		// Submit a simple tx
		op := txnbuild.Payment{
			Destination: master.Address(),
			Amount:      "10",
			Asset:       txnbuild.NativeAsset{},
		}

		txResp := itest.MustSubmitOperations(itest.MasterAccount(), master, &op)
		tt.Equal(master.Address(), txResp.Account)
		tt.Equal("1", txResp.AccountSequence)
	})
}

func TestLiquidityPoolHappyPath(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol18Test(t)
	master := itest.Master()

	keys, accounts := itest.CreateAccounts(2, "1000")
	shareKeys, shareAccount := keys[0], accounts[0]
	tradeKeys, tradeAccount := keys[1], accounts[1]

	itest.MustSubmitMultiSigOperations(shareAccount, []*keypair.Full{shareKeys, master},
		&txnbuild.ChangeTrust{
			Line: txnbuild.ChangeTrustAssetWrapper{
				Asset: txnbuild.CreditAsset{
					Code:   "USD",
					Issuer: master.Address(),
				},
			},
			Limit: txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.ChangeTrust{
			Line: txnbuild.LiquidityPoolShareChangeTrustAsset{
				LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
					AssetA: txnbuild.NativeAsset{},
					AssetB: txnbuild.CreditAsset{
						Code:   "USD",
						Issuer: master.Address(),
					},
					Fee: 30,
				},
			},
			Limit: txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.Payment{
			SourceAccount: master.Address(),
			Destination:   shareAccount.GetAccountID(),
			Asset: txnbuild.CreditAsset{
				Code:   "USD",
				Issuer: master.Address(),
			},
			Amount: "1000",
		},
	)

	poolID, err := xdr.NewPoolId(
		xdr.MustNewNativeAsset(),
		xdr.MustNewCreditAsset("USD", master.Address()),
		30,
	)
	tt.NoError(err)
	poolIDHexString := xdr.Hash(poolID).HexString()

	pools, err := itest.Client().LiquidityPools(auroraclient.LiquidityPoolsRequest{})
	tt.NoError(err)
	tt.Len(pools.Embedded.Records, 1)

	pool := pools.Embedded.Records[0]
	tt.Equal(poolIDHexString, pool.ID)
	tt.Equal(uint32(30), pool.FeeBP)
	tt.Equal("constant_product", pool.Type)
	tt.Equal("0.0000000", pool.TotalShares)
	tt.Equal(uint64(1), pool.TotalTrustlines)

	tt.Equal("0.0000000", pool.Reserves[0].Amount)
	tt.Equal("native", pool.Reserves[0].Asset)
	tt.Equal("0.0000000", pool.Reserves[1].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), pool.Reserves[1].Asset)

	itest.MustSubmitOperations(shareAccount, shareKeys,
		&txnbuild.LiquidityPoolDeposit{
			LiquidityPoolID: [32]byte(poolID),
			MaxAmountA:      "400",
			MaxAmountB:      "777",
			MinPrice:        "0.5",
			MaxPrice:        "2",
		},
	)

	pool, err = itest.Client().LiquidityPoolDetail(auroraclient.LiquidityPoolRequest{
		LiquidityPoolID: poolIDHexString,
	})
	tt.NoError(err)

	tt.Equal(poolIDHexString, pool.ID)
	tt.Equal(uint64(1), pool.TotalTrustlines)

	tt.Equal("400.0000000", pool.Reserves[0].Amount)
	tt.Equal("native", pool.Reserves[0].Asset)
	tt.Equal("777.0000000", pool.Reserves[1].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), pool.Reserves[1].Asset)

	itest.MustSubmitOperations(tradeAccount, tradeKeys,
		&txnbuild.ChangeTrust{
			Line: txnbuild.ChangeTrustAssetWrapper{
				Asset: txnbuild.CreditAsset{
					Code:   "USD",
					Issuer: master.Address(),
				},
			},
			Limit: txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.PathPaymentStrictReceive{
			SendAsset: txnbuild.NativeAsset{},
			DestAsset: txnbuild.CreditAsset{
				Code:   "USD",
				Issuer: master.Address(),
			},
			SendMax:     "1000",
			DestAmount:  "2",
			Destination: tradeKeys.Address(),
		},
	)

	account, err := itest.Client().AccountDetail(auroraclient.AccountRequest{
		AccountID: shareKeys.Address(),
	})
	tt.NoError(err)
	tt.Len(account.Balances, 3)

	liquidityPoolBalance := account.Balances[0]
	tt.Equal("liquidity_pool_shares", liquidityPoolBalance.Asset.Type)
	tt.Equal(poolIDHexString, liquidityPoolBalance.LiquidityPoolId)
	tt.Equal("557.4943945", liquidityPoolBalance.Balance)

	usdBalance := account.Balances[1]
	tt.Equal("credit_alphanum4", usdBalance.Asset.Type)
	tt.Equal("USD", usdBalance.Asset.Code)
	tt.Equal(master.Address(), usdBalance.Asset.Issuer)
	tt.Equal("223.0000000", usdBalance.Balance)

	nativeBalance := account.Balances[2]
	tt.Equal("native", nativeBalance.Asset.Type)

	stats, err := itest.Client().Assets(auroraclient.AssetRequest{})
	tt.NoError(err)
	tt.Len(stats.Embedded.Records, 1)

	stat := stats.Embedded.Records[0]
	tt.Equal("credit_alphanum4", stat.Asset.Type)
	tt.Equal("USD", stat.Asset.Code)
	tt.Equal(master.Address(), stat.Asset.Issuer)
	tt.Equal(int32(2), stat.NumAccounts)
	tt.Equal("225.0000000", stat.Amount)
	tt.Equal(int32(1), stat.NumLiquidityPools)
	tt.Equal("775.0000000", stat.LiquidityPoolsAmount)

	itest.MustSubmitOperations(shareAccount, shareKeys,
		&txnbuild.LiquidityPoolWithdraw{
			LiquidityPoolID: [32]byte(poolID),
			Amount:          pool.TotalShares,
			MinAmountA:      "10",
			MinAmountB:      "20",
		},
	)

	itest.MustSubmitOperations(shareAccount, shareKeys,
		// Clear trustline...
		&txnbuild.Payment{
			Asset: txnbuild.CreditAsset{
				Code:   "USD",
				Issuer: master.Address(),
			},
			Amount:      "998",
			Destination: master.Address(),
		},
		// ...and remove it. It should also remove LP.
		&txnbuild.ChangeTrust{
			Line: txnbuild.LiquidityPoolShareChangeTrustAsset{
				LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
					AssetA: txnbuild.NativeAsset{},
					AssetB: txnbuild.CreditAsset{
						Code:   "USD",
						Issuer: master.Address(),
					},
					Fee: 30,
				},
			},
			Limit: "0",
		},
	)

	account, err = itest.Client().AccountDetail(auroraclient.AccountRequest{
		AccountID: shareKeys.Address(),
	})
	tt.NoError(err)
	tt.Len(account.Balances, 2)

	// Shouldn't contain liquidity_pool_shares balance
	usdBalance = account.Balances[0]
	tt.Equal("credit_alphanum4", usdBalance.Asset.Type)
	tt.Equal("USD", usdBalance.Asset.Code)
	tt.Equal(master.Address(), usdBalance.Asset.Issuer)
	tt.Equal("0.0000000", usdBalance.Balance)

	nativeBalance = account.Balances[1]
	tt.Equal("native", nativeBalance.Asset.Type)

	ops, err := itest.Client().Operations(auroraclient.OperationRequest{
		ForLiquidityPool: poolIDHexString,
	})
	tt.NoError(err)

	// We expect the following ops for this liquidity pool:
	// 1. change_trust creating a trust to LP.
	// 2. liquidity_pool_deposit.
	// 3. path_payment
	// 4. liquidity_pool_withdraw.
	// 5. change_trust removing a trust to LP.
	tt.Len(ops.Embedded.Records, 5)

	op1 := (ops.Embedded.Records[0]).(operations.ChangeTrust)
	tt.Equal("change_trust", op1.Type)
	tt.Equal("liquidity_pool_shares", op1.Asset.Type)
	tt.Equal(poolIDHexString, op1.LiquidityPoolID)
	tt.Equal("922337203685.4775807", op1.Limit)

	op2 := (ops.Embedded.Records[1]).(operations.LiquidityPoolDeposit)
	tt.Equal("liquidity_pool_deposit", op2.Type)
	tt.Equal(poolIDHexString, op2.LiquidityPoolID)
	tt.Equal("0.5000000", op2.MinPrice)
	tt.Equal("2.0000000", op2.MaxPrice)
	tt.Equal("native", op2.ReservesDeposited[0].Asset)
	tt.Equal("400.0000000", op2.ReservesDeposited[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), op2.ReservesDeposited[1].Asset)
	tt.Equal("777.0000000", op2.ReservesDeposited[1].Amount)
	tt.Equal("native", op2.ReservesMax[0].Asset)
	tt.Equal("400.0000000", op2.ReservesMax[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), op2.ReservesMax[1].Asset)
	tt.Equal("777.0000000", op2.ReservesMax[1].Amount)
	tt.Equal("557.4943945", op2.SharesReceived)

	op3 := (ops.Embedded.Records[2]).(operations.PathPayment)
	tt.Equal("path_payment_strict_receive", op3.Payment.Base.Type)
	tt.Equal("2.0000000", op3.Amount)
	tt.Equal("1.0353642", op3.SourceAmount)
	tt.Equal("1000.0000000", op3.SourceMax)
	tt.Equal("native", op3.SourceAssetType)
	tt.Equal("credit_alphanum4", op3.Payment.Asset.Type)
	tt.Equal("USD", op3.Payment.Asset.Code)
	tt.Equal(master.Address(), op3.Payment.Asset.Issuer)

	op4 := (ops.Embedded.Records[3]).(operations.LiquidityPoolWithdraw)
	tt.Equal("liquidity_pool_withdraw", op4.Type)
	tt.Equal(poolIDHexString, op4.LiquidityPoolID)

	tt.Equal("native", op4.ReservesMin[0].Asset)
	tt.Equal("10.0000000", op4.ReservesMin[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), op4.ReservesMin[1].Asset)
	tt.Equal("20.0000000", op4.ReservesMin[1].Amount)

	tt.Equal("native", op4.ReservesReceived[0].Asset)
	tt.Equal("401.0353642", op4.ReservesReceived[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), op4.ReservesReceived[1].Asset)
	tt.Equal("775.0000000", op4.ReservesReceived[1].Amount)

	tt.Equal("557.4943945", op4.Shares)

	op5 := (ops.Embedded.Records[4]).(operations.ChangeTrust)
	tt.Equal("change_trust", op5.Type)
	tt.Equal("liquidity_pool_shares", op5.Asset.Type)
	tt.Equal(poolIDHexString, op5.LiquidityPoolID)
	tt.Equal("0.0000000", op5.Limit)

	effs, err := itest.Client().Effects(auroraclient.EffectRequest{
		ForLiquidityPool: poolIDHexString,
	})
	tt.NoError(err)

	// We expect the following effects for this liquidity pool:
	// 1. trustline_created creating liquidity_pool_shares trust_line
	// 2. liquidity_pool_created
	// 3. liquidity_pool_deposited
	// 4. account_credited - connected to trade
	// 5. account_debited - connected to trade
	// 6. liquidity_pool_trade
	// 7. liquidity_pool_withdrew
	// 8. trustline_removed removing liquidity_pool_shares trust_line
	// 9. liquidity_pool_removed
	tt.Len(effs.Embedded.Records, 9)

	ef1 := (effs.Embedded.Records[0]).(effects.TrustlineCreated)
	tt.Equal(shareKeys.Address(), ef1.Account)
	tt.Equal("trustline_created", ef1.Type)
	tt.Equal("liquidity_pool_shares", ef1.Asset.Type)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef1.LiquidityPoolID)
	tt.Equal("922337203685.4775807", ef1.Limit)

	ef2 := (effs.Embedded.Records[1]).(effects.LiquidityPoolCreated)
	tt.Equal(shareKeys.Address(), ef2.Account)
	tt.Equal("liquidity_pool_created", ef2.Type)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef2.LiquidityPool.ID)
	tt.Equal("constant_product", ef2.LiquidityPool.Type)
	tt.Equal(uint32(30), ef2.LiquidityPool.FeeBP)
	tt.Equal("0.0000000", ef2.LiquidityPool.TotalShares)
	tt.Equal(uint64(1), ef2.LiquidityPool.TotalTrustlines)
	tt.Equal("native", ef2.LiquidityPool.Reserves[0].Asset)
	tt.Equal("0.0000000", ef2.LiquidityPool.Reserves[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef2.LiquidityPool.Reserves[1].Asset)
	tt.Equal("0.0000000", ef2.LiquidityPool.Reserves[1].Amount)

	ef3 := (effs.Embedded.Records[2]).(effects.LiquidityPoolDeposited)
	tt.Equal("liquidity_pool_deposited", ef3.Type)
	tt.Equal(shareKeys.Address(), ef3.Account)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef3.LiquidityPool.ID)
	tt.Equal("constant_product", ef3.LiquidityPool.Type)
	tt.Equal(uint32(30), ef3.LiquidityPool.FeeBP)
	tt.Equal("557.4943945", ef3.LiquidityPool.TotalShares)
	tt.Equal(uint64(1), ef3.LiquidityPool.TotalTrustlines)

	tt.Equal("native", ef3.LiquidityPool.Reserves[0].Asset)
	tt.Equal("400.0000000", ef3.LiquidityPool.Reserves[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef3.LiquidityPool.Reserves[1].Asset)
	tt.Equal("777.0000000", ef3.LiquidityPool.Reserves[1].Amount)

	tt.Equal("native", ef3.ReservesDeposited[0].Asset)
	tt.Equal("400.0000000", ef3.ReservesDeposited[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef3.ReservesDeposited[1].Asset)
	tt.Equal("777.0000000", ef3.ReservesDeposited[1].Amount)

	tt.Equal("557.4943945", ef3.SharesReceived)

	ef4 := (effs.Embedded.Records[3]).(effects.AccountCredited)
	tt.Equal("account_credited", ef4.Base.Type)
	// TODO - is it really LP effect?

	ef5 := (effs.Embedded.Records[4]).(effects.AccountDebited)
	tt.Equal("account_debited", ef5.Base.Type)
	// TODO - is it really LP effect?

	ef6 := (effs.Embedded.Records[5]).(effects.LiquidityPoolTrade)
	tt.Equal("liquidity_pool_trade", ef6.Type)
	tt.Equal(tradeKeys.Address(), ef6.Account)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef6.LiquidityPool.ID)
	tt.Equal("constant_product", ef6.LiquidityPool.Type)
	tt.Equal(uint32(30), ef6.LiquidityPool.FeeBP)
	tt.Equal("557.4943945", ef3.LiquidityPool.TotalShares)
	tt.Equal(uint64(1), ef6.LiquidityPool.TotalTrustlines)
	tt.Equal("native", ef6.LiquidityPool.Reserves[0].Asset)
	tt.Equal("401.0353642", ef6.LiquidityPool.Reserves[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef6.LiquidityPool.Reserves[1].Asset)
	tt.Equal("775.0000000", ef6.LiquidityPool.Reserves[1].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef6.Sold.Asset)
	tt.Equal("2.0000000", ef6.Sold.Amount)
	tt.Equal("native", ef6.Bought.Asset)
	tt.Equal("1.0353642", ef6.Bought.Amount)

	ef7 := (effs.Embedded.Records[6]).(effects.LiquidityPoolWithdrew)
	tt.Equal("liquidity_pool_withdrew", ef7.Type)
	tt.Equal(shareKeys.Address(), ef7.Account)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef7.LiquidityPool.ID)
	tt.Equal("constant_product", ef7.LiquidityPool.Type)
	tt.Equal(uint32(30), ef7.LiquidityPool.FeeBP)
	tt.Equal("0.0000000", ef7.LiquidityPool.TotalShares)
	tt.Equal(uint64(1), ef7.LiquidityPool.TotalTrustlines)

	tt.Equal("native", ef7.LiquidityPool.Reserves[0].Asset)
	tt.Equal("0.0000000", ef7.LiquidityPool.Reserves[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef7.LiquidityPool.Reserves[1].Asset)
	tt.Equal("0.0000000", ef7.LiquidityPool.Reserves[1].Amount)

	tt.Equal("native", ef7.ReservesReceived[0].Asset)
	tt.Equal("401.0353642", ef7.ReservesReceived[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef7.ReservesReceived[1].Asset)
	tt.Equal("775.0000000", ef7.ReservesReceived[1].Amount)

	tt.Equal("557.4943945", ef7.SharesRedeemed)

	ef8 := (effs.Embedded.Records[7]).(effects.TrustlineRemoved)
	tt.Equal("trustline_removed", ef8.Type)
	tt.Equal(shareKeys.Address(), ef8.Account)
	tt.Equal("liquidity_pool_shares", ef8.Asset.Type)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef8.LiquidityPoolID)
	tt.Equal("0.0000000", ef8.Limit)

	ef9 := (effs.Embedded.Records[8]).(effects.LiquidityPoolRemoved)
	tt.Equal("liquidity_pool_removed", ef9.Type)
	tt.Equal(shareKeys.Address(), ef9.Account)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef9.LiquidityPoolID)

	trades, err := itest.Client().Trades(auroraclient.TradeRequest{})
	tt.NoError(err)

	tt.Len(trades.Embedded.Records, 1)

	trade1 := trades.Embedded.Records[0]
	tt.Equal("liquidity_pool", trade1.TradeType)

	tt.Equal(poolIDHexString, trade1.BaseLiquidityPoolID)
	tt.Equal(uint32(30), trade1.LiquidityPoolFeeBP)
	tt.Equal("2.0000000", trade1.BaseAmount)
	tt.Equal("credit_alphanum4", trade1.BaseAssetType)
	tt.Equal("USD", trade1.BaseAssetCode)
	tt.Equal(master.Address(), trade1.BaseAssetIssuer)

	tt.Equal(tradeKeys.Address(), trade1.CounterAccount)
	tt.Equal("1.0353642", trade1.CounterAmount)
	tt.Equal("native", trade1.CounterAssetType)

	tt.Equal(int64(10353642), trade1.Price.N)
	tt.Equal(int64(20000000), trade1.Price.D)
}

func TestLiquidityPoolRevoke(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol18Test(t)
	master := itest.Master()

	keys, accounts := itest.CreateAccounts(2, "1000")
	shareKeys, shareAccount := keys[0], accounts[0]

	poolID, err := xdr.NewPoolId(
		xdr.MustNewNativeAsset(),
		xdr.MustNewCreditAsset("USD", master.Address()),
		30,
	)
	tt.NoError(err)
	poolIDHexString := xdr.Hash(poolID).HexString()

	itest.MustSubmitMultiSigOperations(shareAccount, []*keypair.Full{shareKeys, master},
		&txnbuild.SetOptions{
			SourceAccount: master.Address(),
			SetFlags: []txnbuild.AccountFlag{
				txnbuild.AuthRevocable,
			},
		},
		&txnbuild.ChangeTrust{
			Line: txnbuild.ChangeTrustAssetWrapper{
				Asset: txnbuild.CreditAsset{
					Code:   "USD",
					Issuer: master.Address(),
				},
			},
			Limit: txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.ChangeTrust{
			Line: txnbuild.LiquidityPoolShareChangeTrustAsset{
				LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
					AssetA: txnbuild.NativeAsset{},
					AssetB: txnbuild.CreditAsset{
						Code:   "USD",
						Issuer: master.Address(),
					},
					Fee: 30,
				},
			},
			Limit: txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.Payment{
			SourceAccount: master.Address(),
			Destination:   shareAccount.GetAccountID(),
			Asset: txnbuild.CreditAsset{
				Code:   "USD",
				Issuer: master.Address(),
			},
			Amount: "1000",
		},
		&txnbuild.LiquidityPoolDeposit{
			LiquidityPoolID: [32]byte(poolID),
			MaxAmountA:      "400",
			MaxAmountB:      "777",
			MinPrice:        "0.5",
			MaxPrice:        "2",
		},
		&txnbuild.SetTrustLineFlags{
			SourceAccount: master.Address(),
			Trustor:       shareKeys.Address(),
			Asset: txnbuild.CreditAsset{
				Code:   "USD",
				Issuer: master.Address(),
			},
			ClearFlags: []txnbuild.TrustLineFlag{
				txnbuild.TrustLineAuthorized,
			},
		},
	)

	// Check if claimable balances have been created
	claimableBalances, err := itest.Client().ClaimableBalances(auroraclient.ClaimableBalanceRequest{})
	tt.NoError(err)
	tt.Len(claimableBalances.Embedded.Records, 2)

	// The list is sorted by ID and preimage consists of Account ID which can
	// differ between test runs. Flip the order if the first one is no native.
	if claimableBalances.Embedded.Records[0].Asset != "native" {
		claimableBalances.Embedded.Records[0], claimableBalances.Embedded.Records[1] =
			claimableBalances.Embedded.Records[1], claimableBalances.Embedded.Records[0]
	}

	cb1 := claimableBalances.Embedded.Records[0]
	tt.Equal("native", cb1.Asset)
	tt.Equal("400.0000000", cb1.Amount)
	tt.Equal(shareAccount.GetAccountID(), cb1.Claimants[0].Destination)
	tt.Equal(xdr.ClaimPredicateTypeClaimPredicateUnconditional, cb1.Claimants[0].Predicate.Type)

	cb2 := claimableBalances.Embedded.Records[1]
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), cb2.Asset)
	tt.Equal("777.0000000", cb2.Amount)
	tt.Equal(shareAccount.GetAccountID(), cb2.Claimants[0].Destination)
	tt.Equal(xdr.ClaimPredicateTypeClaimPredicateUnconditional, cb2.Claimants[0].Predicate.Type)

	ops, err := itest.Client().Operations(auroraclient.OperationRequest{
		ForLiquidityPool: poolIDHexString,
	})
	tt.NoError(err)

	// We expect the following ops for this liquidity pool:
	// 1. change_trust creating a trust to LP.
	// 2. liquidity_pool_deposit.
	// 3. set_trust_line_flags revoking assets from LP.
	tt.Len(ops.Embedded.Records, 3)

	op1 := (ops.Embedded.Records[0]).(operations.ChangeTrust)
	tt.Equal("change_trust", op1.Type)
	tt.Equal("liquidity_pool_shares", op1.Asset.Type)
	tt.Equal(poolIDHexString, op1.LiquidityPoolID)
	tt.Equal("922337203685.4775807", op1.Limit)

	op2 := (ops.Embedded.Records[1]).(operations.LiquidityPoolDeposit)
	tt.Equal("liquidity_pool_deposit", op2.Type)
	tt.Equal(poolIDHexString, op2.LiquidityPoolID)
	tt.Equal("0.5000000", op2.MinPrice)
	tt.Equal("2.0000000", op2.MaxPrice)
	tt.Equal("native", op2.ReservesDeposited[0].Asset)
	tt.Equal("400.0000000", op2.ReservesDeposited[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), op2.ReservesDeposited[1].Asset)
	tt.Equal("777.0000000", op2.ReservesDeposited[1].Amount)
	tt.Equal("557.4943945", op2.SharesReceived)

	op3 := (ops.Embedded.Records[2]).(operations.SetTrustLineFlags)
	tt.Equal("set_trust_line_flags", op3.Base.Type)
	tt.Equal("credit_alphanum4", op3.Asset.Type)
	tt.Equal("USD", op3.Asset.Code)
	tt.Equal(master.Address(), op3.Asset.Issuer)
	tt.Equal("authorized", op3.ClearFlagsS[0])

	effs, err := itest.Client().Effects(auroraclient.EffectRequest{
		ForLiquidityPool: poolIDHexString,
		Limit:            20,
	})
	tt.NoError(err)
	// We expect the following effects for this liquidity pool:
	// 1. trustline_created creating liquidity_pool_shares trust_line
	// 2. liquidity_pool_created
	// 3. liquidity_pool_deposited
	// 4. trustline_flags_updated - revoking LP assets
	// 5. claimable_balance_created - creating CB for asset A
	// 6. claimable_balance_claimant_created - claimant for CB above
	// 7. claimable_balance_created - creating CB for asset B
	// 8. claimable_balance_claimant_created - claimant for CB above
	// 9. liquidity_pool_revoked
	// 10. claimable_balance_sponsorship_created
	// 11. claimable_balance_sponsorship_created
	// 12. liquidity_pool_removed - because no more assets inside
	tt.Len(effs.Embedded.Records, 12)

	ef1 := (effs.Embedded.Records[0]).(effects.TrustlineCreated)
	tt.Equal(shareKeys.Address(), ef1.Account)
	tt.Equal("trustline_created", ef1.Type)
	tt.Equal("liquidity_pool_shares", ef1.Asset.Type)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef1.LiquidityPoolID)
	tt.Equal("922337203685.4775807", ef1.Limit)

	ef2 := (effs.Embedded.Records[1]).(effects.LiquidityPoolCreated)
	tt.Equal(shareKeys.Address(), ef2.Account)
	tt.Equal("liquidity_pool_created", ef2.Type)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef2.LiquidityPool.ID)
	tt.Equal("constant_product", ef2.LiquidityPool.Type)
	tt.Equal(uint32(30), ef2.LiquidityPool.FeeBP)
	tt.Equal("0.0000000", ef2.LiquidityPool.TotalShares)
	tt.Equal(uint64(1), ef2.LiquidityPool.TotalTrustlines)
	tt.Equal("native", ef2.LiquidityPool.Reserves[0].Asset)
	tt.Equal("0.0000000", ef2.LiquidityPool.Reserves[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef2.LiquidityPool.Reserves[1].Asset)
	tt.Equal("0.0000000", ef2.LiquidityPool.Reserves[1].Amount)

	ef3 := (effs.Embedded.Records[2]).(effects.LiquidityPoolDeposited)
	tt.Equal("liquidity_pool_deposited", ef3.Type)
	tt.Equal(shareKeys.Address(), ef3.Account)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef3.LiquidityPool.ID)
	tt.Equal("constant_product", ef3.LiquidityPool.Type)
	tt.Equal(uint32(30), ef3.LiquidityPool.FeeBP)
	tt.Equal("557.4943945", ef3.LiquidityPool.TotalShares)
	tt.Equal(uint64(1), ef3.LiquidityPool.TotalTrustlines)

	ef4 := (effs.Embedded.Records[3]).(effects.TrustlineFlagsUpdated)
	tt.Equal("trustline_flags_updated", ef4.Base.Type)
	tt.Equal(master.Address(), ef4.Account)
	tt.Equal("USD", ef4.Asset.Code)
	tt.Equal(master.Address(), ef4.Asset.Issuer)
	tt.Equal(shareAccount.GetAccountID(), ef4.Trustor)

	ef5 := (effs.Embedded.Records[4]).(effects.ClaimableBalanceCreated)
	tt.Equal("claimable_balance_created", ef5.Type)
	tt.Equal("native", ef5.Asset)
	tt.Equal("400.0000000", ef5.Amount)

	ef6 := (effs.Embedded.Records[5]).(effects.ClaimableBalanceClaimantCreated)
	tt.Equal("claimable_balance_claimant_created", ef6.Type)
	tt.Equal("native", ef6.Asset)
	tt.Equal("400.0000000", ef6.Amount)
	tt.Equal(shareKeys.Address(), ef6.Account)
	tt.Equal(xdr.ClaimPredicateTypeClaimPredicateUnconditional, ef6.Predicate.Type)

	ef7 := (effs.Embedded.Records[6]).(effects.ClaimableBalanceCreated)
	tt.Equal("claimable_balance_created", ef7.Type)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef7.Asset)
	tt.Equal("777.0000000", ef7.Amount)

	ef8 := (effs.Embedded.Records[7]).(effects.ClaimableBalanceClaimantCreated)
	tt.Equal("claimable_balance_claimant_created", ef8.Type)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef8.Asset)
	tt.Equal("777.0000000", ef8.Amount)
	tt.Equal(shareKeys.Address(), ef8.Account)
	tt.Equal(xdr.ClaimPredicateTypeClaimPredicateUnconditional, ef8.Predicate.Type)

	ef9 := (effs.Embedded.Records[8]).(effects.LiquidityPoolRevoked)
	tt.Equal("liquidity_pool_revoked", ef9.Type)
	tt.Equal(master.Address(), ef9.Account)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef9.LiquidityPool.ID)
	tt.Equal("constant_product", ef9.LiquidityPool.Type)
	tt.Equal(uint32(30), ef9.LiquidityPool.FeeBP)
	tt.Equal("557.4943945", ef9.LiquidityPool.TotalShares)
	tt.Equal(uint64(1), ef9.LiquidityPool.TotalTrustlines)
	tt.Equal("native", ef9.LiquidityPool.Reserves[0].Asset)
	tt.Equal("400.0000000", ef9.LiquidityPool.Reserves[0].Amount)
	tt.Equal(fmt.Sprintf("USD:%s", master.Address()), ef9.LiquidityPool.Reserves[1].Asset)
	tt.Equal("777.0000000", ef9.LiquidityPool.Reserves[1].Amount)

	// ef10 and ef11 are `claimable_balance_sponsorship_created` effects not
	// relevant here.

	ef12 := (effs.Embedded.Records[11]).(effects.LiquidityPoolRemoved)
	tt.Equal("liquidity_pool_removed", ef12.Type)
	tt.Equal(master.Address(), ef12.Account)
	tt.Equal("64e163b66108152665ee325cc333211446277c86bfe021b9da6bb1769b0daea1", ef12.LiquidityPoolID)
}

func TestLiquidityPoolFailedDepositAndWithdraw(t *testing.T) {
	tt := assert.New(t)
	itest := NewProtocol18Test(t)

	keys, accounts := itest.CreateAccounts(2, "1000")
	shareKeys, shareAccount := keys[0], accounts[0]

	nonExistentPoolID := [32]byte{0xca, 0xfe}

	// Failing deposit
	tx, err := itest.CreateSignedTransaction(shareAccount, []*keypair.Full{shareKeys},
		&txnbuild.LiquidityPoolDeposit{
			LiquidityPoolID: nonExistentPoolID,
			MaxAmountA:      "400",
			MaxAmountB:      "777",
			MinPrice:        "0.5",
			MaxPrice:        "2",
		},
	)
	_, err = itest.Client().SubmitTransaction(tx)
	tt.Error(err)
	hash, err := tx.HashHex(integration.StandaloneNetworkPassphrase)
	tt.NoError(err)
	opsResponse, err := itest.Client().Operations(auroraclient.OperationRequest{
		ForTransaction: hash,
	})
	tt.NoError(err)
	tt.Len(opsResponse.Embedded.Records, 1)
	deposit := (opsResponse.Embedded.Records[0]).(operations.LiquidityPoolDeposit)
	tt.Equal("liquidity_pool_deposit", deposit.Type)
	tt.Equal("cafe000000000000000000000000000000000000000000000000000000000000", deposit.LiquidityPoolID)
	tt.Equal("0.5000000", deposit.MinPrice)
	tt.Equal("2.0000000", deposit.MaxPrice)
	tt.Equal("", deposit.ReservesDeposited[0].Asset)
	tt.Equal("0.0000000", deposit.ReservesDeposited[0].Amount)
	tt.Equal("", deposit.ReservesDeposited[1].Asset)
	tt.Equal("0.0000000", deposit.ReservesDeposited[1].Amount)
	tt.Equal("", deposit.ReservesMax[0].Asset)
	tt.Equal("400.0000000", deposit.ReservesMax[0].Amount)
	tt.Equal("", deposit.ReservesMax[1].Asset)
	tt.Equal("777.0000000", deposit.ReservesMax[1].Amount)
	tt.Equal("0.0000000", deposit.SharesReceived)

	// Failing withdrawal
	tx, err = itest.CreateSignedTransaction(shareAccount, []*keypair.Full{shareKeys},
		&txnbuild.LiquidityPoolWithdraw{
			LiquidityPoolID: nonExistentPoolID,
			Amount:          amount.StringFromInt64(int64(10)),
			MinAmountA:      "10",
			MinAmountB:      "20",
		},
	)
	_, err = itest.Client().SubmitTransaction(tx)
	tt.Error(err)

	hash, err = tx.HashHex(integration.StandaloneNetworkPassphrase)
	tt.NoError(err)
	opsResponse, err = itest.Client().Operations(auroraclient.OperationRequest{
		ForTransaction: hash,
	})
	tt.NoError(err)
	tt.Len(opsResponse.Embedded.Records, 1)
	withdrawal := (opsResponse.Embedded.Records[0]).(operations.LiquidityPoolWithdraw)
	tt.Equal("liquidity_pool_withdraw", withdrawal.Type)
	tt.Equal("cafe000000000000000000000000000000000000000000000000000000000000", withdrawal.LiquidityPoolID)

	tt.Equal("", withdrawal.ReservesMin[0].Asset)
	tt.Equal("10.0000000", withdrawal.ReservesMin[0].Amount)
	tt.Equal("", withdrawal.ReservesMin[1].Asset)
	tt.Equal("20.0000000", withdrawal.ReservesMin[1].Amount)

	tt.Equal("", withdrawal.ReservesReceived[0].Asset)
	tt.Equal("0.0000000", withdrawal.ReservesReceived[0].Amount)
	tt.Equal("", withdrawal.ReservesReceived[1].Asset)
	tt.Equal("0.0000000", withdrawal.ReservesReceived[1].Amount)

	tt.Equal("0.0000010", withdrawal.Shares)
}

// TestProtocol18TradesPriceBreakingChanges checks if the breaking changes in
// `/trades` and `/trade_aggegations` (`price.n` and `price.d` change to string)
// activate after the protocol upgrade.
func TestProtocol18TradesPriceBreakingChanges(t *testing.T) {
	tt := assert.New(t)
	// Start at Protocol 17!
	config := integration.Config{ProtocolVersion: 17}
	itest := integration.NewTest(t, config)
	master := itest.Master()

	keys, accounts := itest.CreateAccounts(1, "1000")
	tradeKeys, tradeAccount := keys[0], accounts[0]

	itest.MustSubmitMultiSigOperations(tradeAccount, []*keypair.Full{tradeKeys, master},
		&txnbuild.ChangeTrust{
			Line: txnbuild.ChangeTrustAssetWrapper{
				Asset: txnbuild.CreditAsset{
					Code:   "USD",
					Issuer: master.Address(),
				},
			},
			Limit: txnbuild.MaxTrustlineLimit,
		},
		&txnbuild.ManageBuyOffer{
			Selling: txnbuild.NativeAsset{},
			Buying: txnbuild.CreditAsset{
				Code:   "USD",
				Issuer: master.Address(),
			},
			Amount: "500",
			Price:  "1",
		},
		&txnbuild.ManageBuyOffer{
			SourceAccount: master.Address(),
			Selling: txnbuild.CreditAsset{
				Code:   "USD",
				Issuer: master.Address(),
			},
			Buying: txnbuild.NativeAsset{},
			Amount: "500",
			Price:  "1",
		},
	)

	body, err := getAuroraResponseBody(itest, "/trades")
	tt.NoError(err)
	tt.Contains(body, `"offer_id": "1",`) // Contains non-empty offer_id
	tt.Contains(body, `"n": 1,`)          // Number
	tt.Contains(body, `"d": 1`)           // Number

	body, err = getAuroraResponseBody(itest,
		"/trade_aggregations?"+
			"base_asset_type=credit_alphanum4&"+
			"base_asset_code=USD&"+
			"base_asset_issuer="+master.Address()+"&"+
			"counter_asset_type=native&"+
			"start_time="+strconv.FormatInt(time.Now().Unix()-3600, 10)+"000&"+
			"end_time="+strconv.FormatInt(time.Now().Unix()+60, 10)+"000&"+
			"resolution=60000&order=desc&limit=200")
	tt.NoError(err)
	tt.Contains(body, `"N": 1,`) // Number
	tt.Contains(body, `"D": 1`)  // Number

	// Now upgrade to protocol 18 and update core info manually to check
	// responses format
	itest.UpgradeProtocol(18)
	itest.Aurora().UpdateDiamcircleCoreInfo(context.Background())

	body, err = getAuroraResponseBody(itest, "/trades")
	tt.NoError(err)
	tt.NotContains(body, `"offer_id"`) // No offer_id
	tt.Contains(body, `"n": "1",`)     // String
	tt.Contains(body, `"d": "1"`)      // String

	body, err = getAuroraResponseBody(itest,
		"/trade_aggregations?"+
			"base_asset_type=credit_alphanum4&"+
			"base_asset_code=USD&"+
			"base_asset_issuer="+master.Address()+"&"+
			"counter_asset_type=native&"+
			"start_time="+strconv.FormatInt(time.Now().Unix()-3600, 10)+"000&"+
			"end_time="+strconv.FormatInt(time.Now().Unix()+60, 10)+"000&"+
			"resolution=60000&order=desc&limit=200")
	tt.NoError(err)
	tt.Contains(body, `"n": "1",`)  // String
	tt.Contains(body, `"n": "1"`)   // String
	tt.NotContains(body, `"N": 1,`) // Not Number
	tt.NotContains(body, `"D": 1`)  // Not Number
}

func getAuroraResponseBody(itest *integration.Test, query string) (string, error) {
	query = strings.TrimLeft(query, "/")
	resp, err := http.Get(fmt.Sprintf("%s%s", itest.Client().AuroraURL, query))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(body), err
}

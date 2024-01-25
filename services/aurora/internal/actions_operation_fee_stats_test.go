package aurora

import (
	"encoding/json"
	"testing"

	hProtocol "go/protocols/aurora"
)

func TestOperationFeeTestsActions_Show(t *testing.T) {
	testCases := []struct {
		scenario            string
		lastbasefee         int64
		maxFeeMax           int
		maxFeeMin           int
		maxFeeMode          int
		maxFeeP10           int
		maxFeeP20           int
		maxFeeP30           int
		maxFeeP40           int
		maxFeeP50           int
		maxFeeP60           int
		maxFeeP70           int
		maxFeeP80           int
		maxFeeP90           int
		maxFeeP95           int
		maxFeeP99           int
		feeChargedMax       int64
		feeChargedMin       int64
		feeChargedMode      int64
		feeChargedP10       int64
		feeChargedP20       int64
		feeChargedP30       int64
		feeChargedP40       int64
		feeChargedP50       int64
		feeChargedP60       int64
		feeChargedP70       int64
		feeChargedP80       int64
		feeChargedP90       int64
		feeChargedP95       int64
		feeChargedP99       int64
		ledgerCapacityUsage float64
	}{
		// happy path
		{
			scenario:            "operation_fee_stats_1",
			lastbasefee:         100,
			maxFeeMax:           100,
			maxFeeMin:           100,
			maxFeeMode:          100,
			maxFeeP10:           100,
			maxFeeP20:           100,
			maxFeeP30:           100,
			maxFeeP40:           100,
			maxFeeP50:           100,
			maxFeeP60:           100,
			maxFeeP70:           100,
			maxFeeP80:           100,
			maxFeeP90:           100,
			maxFeeP95:           100,
			maxFeeP99:           100,
			feeChargedMax:       100,
			feeChargedMin:       100,
			feeChargedMode:      100,
			feeChargedP10:       100,
			feeChargedP20:       100,
			feeChargedP30:       100,
			feeChargedP40:       100,
			feeChargedP50:       100,
			feeChargedP60:       100,
			feeChargedP70:       100,
			feeChargedP80:       100,
			feeChargedP90:       100,
			feeChargedP95:       100,
			feeChargedP99:       100,
			ledgerCapacityUsage: 0.04,
		},
		// no transactions in last 5 ledgers
		{
			scenario:            "operation_fee_stats_2",
			ledgerCapacityUsage: 0.00,
			lastbasefee:         100,
			maxFeeMax:           100,
			maxFeeMin:           100,
			maxFeeMode:          100,
			maxFeeP10:           100,
			maxFeeP20:           100,
			maxFeeP30:           100,
			maxFeeP40:           100,
			maxFeeP50:           100,
			maxFeeP60:           100,
			maxFeeP70:           100,
			maxFeeP80:           100,
			maxFeeP90:           100,
			maxFeeP95:           100,
			maxFeeP99:           100,
			feeChargedMax:       100,
			feeChargedMin:       100,
			feeChargedMode:      100,
			feeChargedP10:       100,
			feeChargedP20:       100,
			feeChargedP30:       100,
			feeChargedP40:       100,
			feeChargedP50:       100,
			feeChargedP60:       100,
			feeChargedP70:       100,
			feeChargedP80:       100,
			feeChargedP90:       100,
			feeChargedP95:       100,
			feeChargedP99:       100,
		},
		// transactions with varying fees
		{
			scenario:            "operation_fee_stats_3",
			ledgerCapacityUsage: 0.03,
			lastbasefee:         100,
			maxFeeMax:           400,
			maxFeeMin:           200,
			maxFeeMode:          400,
			maxFeeP10:           200,
			maxFeeP20:           300,
			maxFeeP30:           400,
			maxFeeP40:           400,
			maxFeeP50:           400,
			maxFeeP60:           400,
			maxFeeP70:           400,
			maxFeeP80:           400,
			maxFeeP90:           400,
			maxFeeP95:           400,
			maxFeeP99:           400,
			feeChargedMax:       100,
			feeChargedMin:       100,
			feeChargedMode:      100,
			feeChargedP10:       100,
			feeChargedP20:       100,
			feeChargedP30:       100,
			feeChargedP40:       100,
			feeChargedP50:       100,
			feeChargedP60:       100,
			feeChargedP70:       100,
			feeChargedP80:       100,
			feeChargedP90:       100,
			feeChargedP95:       100,
			feeChargedP99:       100,
		},
	}

	for _, kase := range testCases {
		t.Run("/fee_stats", func(t *testing.T) {
			ht := StartHTTPTest(t, kase.scenario)
			defer ht.Finish()

			// Update max_tx_set_size on ledgers
			_, err := ht.auroraSession().ExecRaw(ht.Ctx, "UPDATE history_ledgers SET max_tx_set_size = 50")
			ht.Require.NoError(err)

			ht.App.UpdateFeeStatsState(ht.Ctx)

			w := ht.Get("/fee_stats")

			if ht.Assert.Equal(200, w.Code) {
				var result hProtocol.FeeStats
				err := json.Unmarshal(w.Body.Bytes(), &result)
				ht.Require.NoError(err)
				ht.Assert.Equal(kase.lastbasefee, result.LastLedgerBaseFee, "base_fee")
				ht.Assert.Equal(kase.ledgerCapacityUsage, result.LedgerCapacityUsage, "ledger_capacity_usage")

				ht.Assert.Equal(int64(kase.maxFeeMin), result.MaxFee.Min, "min")
				ht.Assert.Equal(int64(kase.maxFeeMode), result.MaxFee.Mode, "mode")
				ht.Assert.Equal(int64(kase.maxFeeP10), result.MaxFee.P10, "p10")
				ht.Assert.Equal(int64(kase.maxFeeP20), result.MaxFee.P20, "p20")
				ht.Assert.Equal(int64(kase.maxFeeP30), result.MaxFee.P30, "p30")
				ht.Assert.Equal(int64(kase.maxFeeP40), result.MaxFee.P40, "p40")
				ht.Assert.Equal(int64(kase.maxFeeP50), result.MaxFee.P50, "p50")
				ht.Assert.Equal(int64(kase.maxFeeP60), result.MaxFee.P60, "p60")
				ht.Assert.Equal(int64(kase.maxFeeP70), result.MaxFee.P70, "p70")
				ht.Assert.Equal(int64(kase.maxFeeP80), result.MaxFee.P80, "p80")
				ht.Assert.Equal(int64(kase.maxFeeP90), result.MaxFee.P90, "p90")
				ht.Assert.Equal(int64(kase.maxFeeP95), result.MaxFee.P95, "p95")
				ht.Assert.Equal(int64(kase.maxFeeP99), result.MaxFee.P99, "p99")

				ht.Assert.Equal(kase.feeChargedMax, result.FeeCharged.Max, "fee_charged_max")
				ht.Assert.Equal(kase.feeChargedMin, result.FeeCharged.Min, "fee_charged_min")
				ht.Assert.Equal(kase.feeChargedMode, result.FeeCharged.Mode, "fee_charged_mode")
				ht.Assert.Equal(kase.feeChargedP10, result.FeeCharged.P10, "fee_charged_p10")
				ht.Assert.Equal(kase.feeChargedP20, result.FeeCharged.P20, "fee_charged_p20")
				ht.Assert.Equal(kase.feeChargedP30, result.FeeCharged.P30, "fee_charged_p30")
				ht.Assert.Equal(kase.feeChargedP40, result.FeeCharged.P40, "fee_charged_p40")
				ht.Assert.Equal(kase.feeChargedP50, result.FeeCharged.P50, "fee_charged_p50")
				ht.Assert.Equal(kase.feeChargedP60, result.FeeCharged.P60, "fee_charged_p60")
				ht.Assert.Equal(kase.feeChargedP70, result.FeeCharged.P70, "fee_charged_p70")
				ht.Assert.Equal(kase.feeChargedP80, result.FeeCharged.P80, "fee_charged_p80")
				ht.Assert.Equal(kase.feeChargedP90, result.FeeCharged.P90, "fee_charged_p90")
				ht.Assert.Equal(kase.feeChargedP95, result.FeeCharged.P95, "fee_charged_p95")
				ht.Assert.Equal(kase.feeChargedP99, result.FeeCharged.P99, "fee_charged_p99")
			}
		})
	}
}

// TestOperationFeeTestsActions_ShowMultiOp tests fee stats in case transactions contain multiple operations.
// In such case, since protocol v11, we should use number of operations as the indicator of ledger capacity usage.
func TestOperationFeeTestsActions_ShowMultiOp(t *testing.T) {
	ht := StartHTTPTest(t, "operation_fee_stats_3")
	defer ht.Finish()

	// Update max_tx_set_size on ledgers
	_, err := ht.auroraSession().ExecRaw(ht.Ctx, "UPDATE history_ledgers SET max_tx_set_size = 50")
	ht.Require.NoError(err)

	// Update number of ops on each transaction
	_, err = ht.auroraSession().ExecRaw(ht.Ctx, "UPDATE history_transactions SET operation_count = operation_count * 2")
	ht.Require.NoError(err)

	ht.App.UpdateFeeStatsState(ht.Ctx)

	w := ht.Get("/fee_stats")

	if ht.Assert.Equal(200, w.Code) {
		var result hProtocol.FeeStats
		err := json.Unmarshal(w.Body.Bytes(), &result)
		ht.Require.NoError(err)
		ht.Assert.Equal(int64(100), result.LastLedgerBaseFee, "base_fee")
		ht.Assert.Equal(0.06, result.LedgerCapacityUsage, "ledger_capacity_usage")

		ht.Assert.Equal(int64(200), result.MaxFee.Max, "max_fee_max")
		ht.Assert.Equal(int64(100), result.MaxFee.Min, "max_fee_min")
		ht.Assert.Equal(int64(200), result.MaxFee.Mode, "max_fee_mode")
		ht.Assert.Equal(int64(100), result.MaxFee.P10, "max_fee_p10")
		ht.Assert.Equal(int64(150), result.MaxFee.P20, "max_fee_p20")
		ht.Assert.Equal(int64(200), result.MaxFee.P30, "max_fee_p30")
		ht.Assert.Equal(int64(200), result.MaxFee.P40, "max_fee_p40")
		ht.Assert.Equal(int64(200), result.MaxFee.P50, "max_fee_p50")
		ht.Assert.Equal(int64(200), result.MaxFee.P60, "max_fee_p60")
		ht.Assert.Equal(int64(200), result.MaxFee.P70, "max_fee_p70")
		ht.Assert.Equal(int64(200), result.MaxFee.P80, "max_fee_p80")
		ht.Assert.Equal(int64(200), result.MaxFee.P90, "max_fee_p90")
		ht.Assert.Equal(int64(200), result.MaxFee.P95, "max_fee_p95")
		ht.Assert.Equal(int64(200), result.MaxFee.P99, "max_fee_p99")

		ht.Assert.Equal(int64(50), result.FeeCharged.Max, "fee_charged_max")
		ht.Assert.Equal(int64(50), result.FeeCharged.Min, "fee_charged_min")
		ht.Assert.Equal(int64(50), result.FeeCharged.Mode, "fee_charged_mode")
		ht.Assert.Equal(int64(50), result.FeeCharged.P10, "fee_charged_p10")
		ht.Assert.Equal(int64(50), result.FeeCharged.P20, "fee_charged_p20")
		ht.Assert.Equal(int64(50), result.FeeCharged.P30, "fee_charged_p30")
		ht.Assert.Equal(int64(50), result.FeeCharged.P40, "fee_charged_p40")
		ht.Assert.Equal(int64(50), result.FeeCharged.P50, "fee_charged_p50")
		ht.Assert.Equal(int64(50), result.FeeCharged.P60, "fee_charged_p60")
		ht.Assert.Equal(int64(50), result.FeeCharged.P70, "fee_charged_p70")
		ht.Assert.Equal(int64(50), result.FeeCharged.P80, "fee_charged_p80")
		ht.Assert.Equal(int64(50), result.FeeCharged.P90, "fee_charged_p90")
		ht.Assert.Equal(int64(50), result.FeeCharged.P95, "fee_charged_p95")
		ht.Assert.Equal(int64(50), result.FeeCharged.P99, "fee_charged_p99")
	}
}

func TestEmptyFeeStats(t *testing.T) {
	ht := StartHTTPTestWithoutScenario(t)
	defer ht.Finish()
	w := ht.Get("/fee_stats")
	ht.Assert.Equal(200, w.Code)
	var result hProtocol.FeeStats
	err := json.Unmarshal(w.Body.Bytes(), &result)
	ht.Require.NoError(err)
	ht.Assert.Equal(result, hProtocol.FeeStats{})
}

func TestOperationFeeTestsActions_NotInterpolating(t *testing.T) {
	ht := StartHTTPTest(t, "operation_fee_stats_3")
	defer ht.Finish()

	// Update max_tx_set_size on ledgers
	_, err := ht.auroraSession().ExecRaw(ht.Ctx, "UPDATE history_ledgers SET max_tx_set_size = 50")
	ht.Require.NoError(err)

	// Update one tx to a huge fee
	_, err = ht.auroraSession().ExecRaw(ht.Ctx, "UPDATE history_transactions SET max_fee = 256000, operation_count = 16 WHERE transaction_hash = '6a349e7331e93a251367287e274fb1699abaf723bde37aebe96248c76fd3071a'")
	ht.Require.NoError(err)

	ht.App.UpdateFeeStatsState(ht.Ctx)

	w := ht.Get("/fee_stats")

	if ht.Assert.Equal(200, w.Code) {
		var result hProtocol.FeeStats
		err := json.Unmarshal(w.Body.Bytes(), &result)
		ht.Require.NoError(err)
		ht.Assert.Equal(int64(100), result.LastLedgerBaseFee, "base_fee")
		ht.Assert.Equal(0.09, result.LedgerCapacityUsage, "ledger_capacity_usage")

		ht.Assert.Equal(int64(16000), result.MaxFee.Max, "max_fee_max")
		ht.Assert.Equal(int64(200), result.MaxFee.Min, "max_fee_min")
		ht.Assert.Equal(int64(400), result.MaxFee.Mode, "max_fee_mode")
		ht.Assert.Equal(int64(200), result.MaxFee.P10, "max_fee_p10")
		ht.Assert.Equal(int64(300), result.MaxFee.P20, "max_fee_p20")
		ht.Assert.Equal(int64(400), result.MaxFee.P30, "max_fee_p30")
		ht.Assert.Equal(int64(400), result.MaxFee.P40, "max_fee_p40")
		ht.Assert.Equal(int64(400), result.MaxFee.P50, "max_fee_p50")
		ht.Assert.Equal(int64(400), result.MaxFee.P60, "max_fee_p60")
		ht.Assert.Equal(int64(400), result.MaxFee.P70, "max_fee_p70")
		ht.Assert.Equal(int64(400), result.MaxFee.P80, "max_fee_p80")
		ht.Assert.Equal(int64(16000), result.MaxFee.P90, "max_fee_p90")
		ht.Assert.Equal(int64(16000), result.MaxFee.P95, "max_fee_p95")
		ht.Assert.Equal(int64(16000), result.MaxFee.P99, "max_fee_p99")

		ht.Assert.Equal(int64(100), result.FeeCharged.Max, "fee_charged_max")
		ht.Assert.Equal(int64(6), result.FeeCharged.Min, "fee_charged_min")
		ht.Assert.Equal(int64(100), result.FeeCharged.Mode, "fee_charged_mode")
		ht.Assert.Equal(int64(6), result.FeeCharged.P10, "fee_charged_p10")
		ht.Assert.Equal(int64(100), result.FeeCharged.P20, "fee_charged_p20")
		ht.Assert.Equal(int64(100), result.FeeCharged.P30, "fee_charged_p30")
		ht.Assert.Equal(int64(100), result.FeeCharged.P40, "fee_charged_p40")
		ht.Assert.Equal(int64(100), result.FeeCharged.P50, "fee_charged_p50")
		ht.Assert.Equal(int64(100), result.FeeCharged.P60, "fee_charged_p60")
		ht.Assert.Equal(int64(100), result.FeeCharged.P70, "fee_charged_p70")
		ht.Assert.Equal(int64(100), result.FeeCharged.P80, "fee_charged_p80")
		ht.Assert.Equal(int64(100), result.FeeCharged.P90, "fee_charged_p90")
		ht.Assert.Equal(int64(100), result.FeeCharged.P95, "fee_charged_p95")
		ht.Assert.Equal(int64(100), result.FeeCharged.P99, "fee_charged_p99")
	}
}

func TestOperationFeeTestsActions_FeeBump(t *testing.T) {
	ht := StartHTTPTest(t, "operation_fee_stats_3")
	defer ht.Finish()

	// Update one tx to be a fee bump
	result, err := ht.auroraSession().ExecRaw(ht.Ctx, "UPDATE history_transactions SET max_fee = 10, new_max_fee = 1000, fee_charged = 200 WHERE transaction_hash = '6a349e7331e93a251367287e274fb1699abaf723bde37aebe96248c76fd3071a'")
	ht.Require.NoError(err)
	rowsAffected, err := result.RowsAffected()
	ht.Require.NoError(err)
	ht.Require.Equal(int64(1), rowsAffected)

	ht.App.UpdateFeeStatsState(ht.Ctx)

	w := ht.Get("/fee_stats")

	if ht.Assert.Equal(200, w.Code) {
		var result hProtocol.FeeStats
		err := json.Unmarshal(w.Body.Bytes(), &result)
		ht.Require.NoError(err)
		ht.Assert.Equal(int64(100), result.LastLedgerBaseFee, "base_fee")

		ht.Assert.Equal(int64(500), result.MaxFee.Max, "max_fee_max")
		ht.Assert.Equal(int64(200), result.MaxFee.Min, "max_fee_min")
		ht.Assert.Equal(int64(400), result.MaxFee.Mode, "max_fee_mode")
		ht.Assert.Equal(int64(200), result.MaxFee.P10, "max_fee_p10")
		ht.Assert.Equal(int64(300), result.MaxFee.P20, "max_fee_p20")
		ht.Assert.Equal(int64(400), result.MaxFee.P30, "max_fee_p30")
		ht.Assert.Equal(int64(400), result.MaxFee.P40, "max_fee_p40")
		ht.Assert.Equal(int64(400), result.MaxFee.P50, "max_fee_p50")
		ht.Assert.Equal(int64(400), result.MaxFee.P60, "max_fee_p60")
		ht.Assert.Equal(int64(400), result.MaxFee.P70, "max_fee_p70")
		ht.Assert.Equal(int64(400), result.MaxFee.P80, "max_fee_p80")
		ht.Assert.Equal(int64(500), result.MaxFee.P90, "max_fee_p90")
		ht.Assert.Equal(int64(500), result.MaxFee.P95, "max_fee_p95")
		ht.Assert.Equal(int64(500), result.MaxFee.P99, "max_fee_p99")

		ht.Assert.Equal(int64(100), result.FeeCharged.Max, "fee_charged_max")
		ht.Assert.Equal(int64(100), result.FeeCharged.Min, "fee_charged_min")
		ht.Assert.Equal(int64(100), result.FeeCharged.Mode, "fee_charged_mode")
		ht.Assert.Equal(int64(100), result.FeeCharged.P10, "fee_charged_p10")
		ht.Assert.Equal(int64(100), result.FeeCharged.P20, "fee_charged_p20")
		ht.Assert.Equal(int64(100), result.FeeCharged.P30, "fee_charged_p30")
		ht.Assert.Equal(int64(100), result.FeeCharged.P40, "fee_charged_p40")
		ht.Assert.Equal(int64(100), result.FeeCharged.P50, "fee_charged_p50")
		ht.Assert.Equal(int64(100), result.FeeCharged.P60, "fee_charged_p60")
		ht.Assert.Equal(int64(100), result.FeeCharged.P70, "fee_charged_p70")
		ht.Assert.Equal(int64(100), result.FeeCharged.P80, "fee_charged_p80")
		ht.Assert.Equal(int64(100), result.FeeCharged.P90, "fee_charged_p90")
		ht.Assert.Equal(int64(100), result.FeeCharged.P95, "fee_charged_p95")
		ht.Assert.Equal(int64(100), result.FeeCharged.P99, "fee_charged_p99")
	}
}

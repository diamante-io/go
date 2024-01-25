package actions

import (
	"net/http"
	"strconv"

	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/aurora/internal/operationfeestats"
)

// FeeStatsHandler is the action handler for the /fee_stats endpoint
type FeeStatsHandler struct {
}

// GetResource fee stats resource
func (handler FeeStatsHandler) GetResource(w HeaderWriter, r *http.Request) (interface{}, error) {
	feeStats := aurora.FeeStats{}

	cur, ok := operationfeestats.CurrentState()
	feeStats.LastLedgerBaseFee = cur.LastBaseFee
	feeStats.LastLedger = cur.LastLedger

	// LedgerCapacityUsage is the empty string when operationfeestats has not had its state set
	if ok {
		capacity, err := strconv.ParseFloat(
			cur.LedgerCapacityUsage,
			64,
		)
		if err != nil {
			return nil, err
		}
		feeStats.LedgerCapacityUsage = capacity
	}

	// FeeCharged
	feeStats.FeeCharged.Max = cur.FeeChargedMax
	feeStats.FeeCharged.Min = cur.FeeChargedMin
	feeStats.FeeCharged.Mode = cur.FeeChargedMode
	feeStats.FeeCharged.P10 = cur.FeeChargedP10
	feeStats.FeeCharged.P20 = cur.FeeChargedP20
	feeStats.FeeCharged.P30 = cur.FeeChargedP30
	feeStats.FeeCharged.P40 = cur.FeeChargedP40
	feeStats.FeeCharged.P50 = cur.FeeChargedP50
	feeStats.FeeCharged.P60 = cur.FeeChargedP60
	feeStats.FeeCharged.P70 = cur.FeeChargedP70
	feeStats.FeeCharged.P80 = cur.FeeChargedP80
	feeStats.FeeCharged.P90 = cur.FeeChargedP90
	feeStats.FeeCharged.P95 = cur.FeeChargedP95
	feeStats.FeeCharged.P99 = cur.FeeChargedP99

	// MaxFee
	feeStats.MaxFee.Max = cur.MaxFeeMax
	feeStats.MaxFee.Min = cur.MaxFeeMin
	feeStats.MaxFee.Mode = cur.MaxFeeMode
	feeStats.MaxFee.P10 = cur.MaxFeeP10
	feeStats.MaxFee.P20 = cur.MaxFeeP20
	feeStats.MaxFee.P30 = cur.MaxFeeP30
	feeStats.MaxFee.P40 = cur.MaxFeeP40
	feeStats.MaxFee.P50 = cur.MaxFeeP50
	feeStats.MaxFee.P60 = cur.MaxFeeP60
	feeStats.MaxFee.P70 = cur.MaxFeeP70
	feeStats.MaxFee.P80 = cur.MaxFeeP80
	feeStats.MaxFee.P90 = cur.MaxFeeP90
	feeStats.MaxFee.P95 = cur.MaxFeeP95
	feeStats.MaxFee.P99 = cur.MaxFeeP99

	return feeStats, nil
}

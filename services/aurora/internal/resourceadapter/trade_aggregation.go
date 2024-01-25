package resourceadapter

import (
	"context"

	"go/amount"
	"go/price"
	protocol "go/protocols/aurora"
	"go/services/aurora/internal/db2/history"
)

// PopulateTradeAggregation fills out the details of a trade aggregation using a row from the trade aggregations
// table.
func PopulateTradeAggregation(
	ctx context.Context,
	dest *protocol.TradeAggregation,
	row history.TradeAggregation,
) error {
	var err error
	dest.Timestamp = row.Timestamp
	dest.TradeCount = row.TradeCount
	dest.BaseVolume, err = amount.IntStringToAmount(row.BaseVolume)
	if err != nil {
		return err
	}
	dest.CounterVolume, err = amount.IntStringToAmount(row.CounterVolume)
	if err != nil {
		return err
	}
	dest.Average = price.StringFromFloat64(row.Average)
	dest.HighR = protocol.TradePrice{
		N: row.HighN,
		D: row.HighD,
	}
	dest.High = dest.HighR.String()
	dest.LowR = protocol.TradePrice{
		N: row.LowN,
		D: row.LowD,
	}
	dest.Low = dest.LowR.String()
	dest.OpenR = protocol.TradePrice{
		N: row.OpenN,
		D: row.OpenD,
	}
	dest.Open = dest.OpenR.String()
	dest.CloseR = protocol.TradePrice{
		N: row.CloseN,
		D: row.CloseD,
	}
	dest.Close = dest.CloseR.String()
	return nil
}

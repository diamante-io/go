package ingest

import (
	"context"
	"database/sql"
	"math/rand"
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"go/exp/orderbook"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ingest/processors"
	"go/support/errors"
	"go/xdr"
)

const (
	verificationFrequency = time.Hour
	updateFrequency       = 2 * time.Second
)

// OrderBookStream updates an in memory graph to be consistent with
// offers in the aurora DB. Any offers which are created, modified, or removed
// from the aurora DB during ingestion will be applied to the in memory order book
// graph. OrderBookStream assumes that no other component will update the
// in memory graph. However, it is safe for other go routines to use the
// in memory graph for read operations.
type OrderBookStream struct {
	graph    orderbook.OBGraph
	historyQ history.IngestionQ
	// LatestLedgerGauge exposes the local (order book graph)
	// latest processed ledger
	LatestLedgerGauge prometheus.Gauge
	lastLedger        uint32
	lastVerification  time.Time
	encodingBuffer    *xdr.EncodingBuffer
}

// NewOrderBookStream constructs and initializes an OrderBookStream instance
func NewOrderBookStream(historyQ history.IngestionQ, graph orderbook.OBGraph) *OrderBookStream {
	return &OrderBookStream{
		graph:    graph,
		historyQ: historyQ,
		LatestLedgerGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "aurora", Subsystem: "order_book_stream", Name: "latest_ledger",
		}),
		lastVerification: time.Now(),
		encodingBuffer:   xdr.NewEncodingBuffer(),
	}
}

type ingestionStatus struct {
	HistoryConsistentWithState        bool
	StateInvalid                      bool
	LastIngestedLedger                uint32
	LastOfferCompactionLedger         uint32
	LastLiquidityPoolCompactionLedger uint32
}

func (o *OrderBookStream) getIngestionStatus(ctx context.Context) (ingestionStatus, error) {
	var status ingestionStatus
	var err error

	status.StateInvalid, err = o.historyQ.GetExpStateInvalid(ctx)
	if err != nil {
		return status, errors.Wrap(err, "Error from GetExpStateInvalid")
	}

	var lastHistoryLedger uint32
	lastHistoryLedger, err = o.historyQ.GetLatestHistoryLedger(ctx)
	if err != nil {
		return status, errors.Wrap(err, "Error from GetLatestHistoryLedger")
	}
	status.LastIngestedLedger, err = o.historyQ.GetLastLedgerIngestNonBlocking(ctx)
	if err != nil {
		return status, errors.Wrap(err, "Error from GetLastLedgerIngestNonBlocking")
	}
	status.LastOfferCompactionLedger, err = o.historyQ.GetOfferCompactionSequence(ctx)
	if err != nil {
		return status, errors.Wrap(err, "Error from GetOfferCompactionSequence")
	}
	status.LastLiquidityPoolCompactionLedger, err = o.historyQ.GetLiquidityPoolCompactionSequence(ctx)
	if err != nil {
		return status, errors.Wrap(err, "Error from GetLiquidityPoolCompactionSequence")
	}

	status.HistoryConsistentWithState = (status.LastIngestedLedger == lastHistoryLedger) ||
		// Running ingestion on an empty DB is a special case because we first ingest from the history archive.
		// Then, on the next iteration, we ingest TX Meta from diamcircle Core. So there is a brief
		// period where there will not be any rows in the history_ledgers table but that is ok.
		(lastHistoryLedger == 0)
	return status, nil
}

// update returns true if the order book graph was reset
func (o *OrderBookStream) update(ctx context.Context, status ingestionStatus) (bool, error) {
	reset := o.lastLedger == 0
	if status.StateInvalid {
		log.WithField("status", status).Warn("ingestion state is invalid")
		reset = true
	} else if !status.HistoryConsistentWithState {
		log.WithField("status", status).
			Info("waiting for ingestion system catchup")
		reset = true
	} else if status.LastIngestedLedger < o.lastLedger {
		log.WithField("status", status).
			WithField("last_ledger", o.lastLedger).
			Warn("ingestion is behind order book last ledger")
		reset = true
	} else if o.lastLedger > 0 && o.lastLedger < status.LastOfferCompactionLedger {
		log.WithField("status", status).
			WithField("last_ledger", o.lastLedger).
			Warn("order book is behind the last offer compaction ledger")
		reset = true
	} else if status.LastOfferCompactionLedger != status.LastLiquidityPoolCompactionLedger {
		log.WithField("status", status).
			WithField("last_ledger", o.lastLedger).
			Warn("offer compaction is not consistentwith liquidity pool compaction")
		reset = true
	}

	if reset {
		o.graph.Clear()
		o.lastLedger = 0

		// wait until offers in aurora db is valid before populating order book graph
		if status.StateInvalid || !status.HistoryConsistentWithState {
			return true, nil
		}

		defer o.graph.Discard()

		offers, err := o.historyQ.GetAllOffers(ctx)
		if err != nil {
			return true, errors.Wrap(err, "Error from GetAllOffers")
		}

		liquidityPools, err := o.historyQ.GetAllLiquidityPools(ctx)
		if err != nil {
			return true, errors.Wrap(err, "Error from GetAllLiquidityPools")
		}

		for _, offer := range offers {
			o.graph.AddOffers(offerToXDR(offer))
		}

		for _, liquidityPool := range liquidityPools {
			liquidityPoolXDR, err := liquidityPoolToXDR(liquidityPool)
			if err != nil {
				return true, errors.Wrap(err, "Invalid liquidity pool row")
			}
			o.graph.AddLiquidityPools(liquidityPoolXDR)
		}

		if err := o.graph.Apply(status.LastIngestedLedger); err != nil {
			return true, errors.Wrap(err, "Error applying changes to order book")
		}

		o.lastLedger = status.LastIngestedLedger
		o.LatestLedgerGauge.Set(float64(status.LastIngestedLedger))
		return true, nil
	}

	if status.LastIngestedLedger == o.lastLedger {
		return false, nil
	}

	defer o.graph.Discard()

	offers, err := o.historyQ.GetUpdatedOffers(ctx, o.lastLedger)
	if err != nil {
		return false, errors.Wrap(err, "Error from GetUpdatedOffers")
	}
	liquidityPools, err := o.historyQ.GetUpdatedLiquidityPools(ctx, o.lastLedger)
	if err != nil {
		return false, errors.Wrap(err, "Error from GetUpdatedLiquidityPools")
	}

	for _, offer := range offers {
		if offer.Deleted {
			o.graph.RemoveOffer(xdr.Int64(offer.OfferID))
		} else {
			o.graph.AddOffers(offerToXDR(offer))
		}
	}

	for _, liquidityPool := range liquidityPools {
		var poolXDR xdr.LiquidityPoolEntry
		poolXDR, err = liquidityPoolToXDR(liquidityPool)
		if err != nil {
			return false, errors.Wrap(err, "Error converting liquidity pool row to xdr")
		}
		if liquidityPool.Deleted {
			o.graph.RemoveLiquidityPool(poolXDR)
		} else {
			o.graph.AddLiquidityPools(poolXDR)
		}
	}

	if err = o.graph.Apply(status.LastIngestedLedger); err != nil {
		return false, errors.Wrap(err, "Error applying changes to order book")
	}

	o.lastLedger = status.LastIngestedLedger
	o.LatestLedgerGauge.Set(float64(status.LastIngestedLedger))
	return false, nil
}

func (o *OrderBookStream) verifyAllOffers(ctx context.Context, offers []xdr.OfferEntry) (bool, error) {
	ingestionOffers, err := o.historyQ.GetAllOffers(ctx)
	if err != nil {
		return false, errors.Wrap(err, "Error from GetAllOffers")
	}

	mismatch := len(offers) != len(ingestionOffers)

	if !mismatch {
		sort.Slice(offers, func(i, j int) bool {
			return offers[i].OfferId < offers[j].OfferId
		})
		sort.Slice(ingestionOffers, func(i, j int) bool {
			return ingestionOffers[i].OfferID < ingestionOffers[j].OfferID
		})

		for i, offerRow := range ingestionOffers {
			offerEntry := offers[i]
			offerRowXDR := offerToXDR(offerRow)
			offerEntryBase64, err := o.encodingBuffer.MarshalBase64(&offerEntry)
			if err != nil {
				return false, errors.Wrap(err, "Error from marshalling offerEntry")
			}
			offerRowBase64, err := o.encodingBuffer.MarshalBase64(&offerRowXDR)
			if err != nil {
				return false, errors.Wrap(err, "Error from marshalling offerRowXDR")
			}
			if offerEntryBase64 != offerRowBase64 {
				mismatch = true
				break
			}
		}
	}

	if mismatch {
		log.WithField("stream_offers", offers).
			WithField("ingestion_offers", ingestionOffers).
			Error("offers derived from order book stream does not match offers from ingestion")
		return false, nil
	}
	log.Info("offer stream verification succeeded")
	return true, nil
}

func (o *OrderBookStream) verifyAllLiquidityPools(ctx context.Context, liquidityPools []xdr.LiquidityPoolEntry) (bool, error) {
	ingestionLiquidityPools, err := o.historyQ.GetAllLiquidityPools(ctx)
	if err != nil {
		return false, errors.Wrap(err, "Error from GetAllLiquidityPools")
	}

	mismatch := len(liquidityPools) != len(ingestionLiquidityPools)

	if !mismatch {
		sort.Slice(liquidityPools, func(i, j int) bool {
			return processors.PoolIDToString(liquidityPools[i].LiquidityPoolId) <
				processors.PoolIDToString(liquidityPools[j].LiquidityPoolId)
		})
		sort.Slice(ingestionLiquidityPools, func(i, j int) bool {
			return ingestionLiquidityPools[i].PoolID < ingestionLiquidityPools[j].PoolID
		})

		for i, liquidityPoolRow := range ingestionLiquidityPools {
			liquidityPoolEntry := liquidityPools[i]
			liquidityPoolRowXDR, err := liquidityPoolToXDR(liquidityPoolRow)
			if err != nil {
				return false, errors.Wrap(err, "Error from converting liquidity pool row to xdr")
			}
			liquidityPoolEntryBase64, err := o.encodingBuffer.MarshalBase64(&liquidityPoolEntry)
			if err != nil {
				return false, errors.Wrap(err, "Error from marshalling liquidityPoolEntry")
			}
			liquidityPoolRowBase64, err := o.encodingBuffer.MarshalBase64(&liquidityPoolRowXDR)
			if err != nil {
				return false, errors.Wrap(err, "Error from marshalling liquidityPoolRowXDR")
			}
			if liquidityPoolEntryBase64 != liquidityPoolRowBase64 {
				mismatch = true
				break
			}
		}
	}

	if mismatch {
		log.WithField("stream_liquidity_pools", liquidityPools).
			WithField("ingestion_liquidity_pools", ingestionLiquidityPools).
			Error("liquidity pools derived from order book stream does not match liquidity pool from ingestion")
		return false, nil
	}
	log.Info("liquidity pool stream verification succeeded")
	return true, nil
}

// Update will query the aurora DB for offers which have been created, removed, or updated since the
// last time Update() was called. Those changes will then be applied to the in memory order book graph.
// After calling this function, the the in memory order book graph should be consistent with the
// aurora DB (assuming no error is returned).
func (o *OrderBookStream) Update(ctx context.Context) error {
	if err := o.historyQ.BeginTx(&sql.TxOptions{ReadOnly: true, Isolation: sql.LevelRepeatableRead}); err != nil {
		return errors.Wrap(err, "Error starting repeatable read transaction")
	}
	defer o.historyQ.Rollback()

	status, err := o.getIngestionStatus(ctx)
	if err != nil {
		return errors.Wrap(err, "Error obtaining ingestion status")
	}

	if reset, err := o.update(ctx, status); err != nil {
		return errors.Wrap(err, "Error updating")
	} else if reset {
		return nil
	}

	// add 15 minute jitter so that not all aurora nodes are calling
	// historyQ.GetAllOffers at the same time
	jitter := time.Duration(rand.Int63n(int64(15 * time.Minute)))
	requiresVerification := o.lastLedger > 0 &&
		time.Since(o.lastVerification) >= verificationFrequency+jitter

	if requiresVerification {
		offers, pools, err := o.graph.Verify()
		if err != nil {
			log.WithError(err).
				Error("Orderbook graph is not internally consistent")
			o.lastVerification = time.Now()
			// set last ledger to 0 so that we reset on next update
			o.lastLedger = 0
			return nil
		}

		offersOk, err := o.verifyAllOffers(ctx, offers)
		if err != nil {
			if !isCancelledError(err) {
				log.WithError(err).Info("Could not verify offers")
				return nil
			}
		}

		liquidityPoolsOK, err := o.verifyAllLiquidityPools(ctx, pools)
		if err != nil {
			if !isCancelledError(err) {
				log.WithError(err).Info("Could not verify liquidity pools")
				return nil
			}
		}
		o.lastVerification = time.Now()
		if !offersOk || !liquidityPoolsOK {
			// set last ledger to 0 so that we reset on next update
			o.lastLedger = 0
		}
	}
	return nil
}

// Run will call Update() every 30 seconds until the given context is terminated.
func (o *OrderBookStream) Run(ctx context.Context) {
	ticker := time.NewTicker(updateFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := o.Update(ctx); err != nil && !isCancelledError(err) {
				log.WithError(err).Error("could not apply updates from order book stream")
			}
		case <-ctx.Done():
			log.Info("shutting down OrderBookStream")
			return
		}
	}
}

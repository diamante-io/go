package history

import (
	"context"
	"database/sql"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/diamcircle/go/support/errors"
)

const (
	ingestVersion = "exp_ingest_version"
	// Distributed ingestion in Aurora relies on this key and it is part
	// of migration files. If you need to update the key name remember
	// to upgrade it in migration files too!
	lastLedgerKey                   = "exp_ingest_last_ledger"
	stateInvalid                    = "exp_state_invalid"
	offerCompactionSequence         = "offer_compaction_sequence"
	liquidityPoolCompactionSequence = "liquidity_pool_compaction_sequence"
)

// GetLastLedgerIngestNonBlocking works like GetLastLedgerIngest but
// it does not block the value and does not return error if the value
// has not been previously set.
// This is used in status reporting (ex. in root resource of Aurora).
func (q *Q) GetLastLedgerIngestNonBlocking(ctx context.Context) (uint32, error) {
	lastIngestedLedger, err := q.getValueFromStore(ctx, lastLedgerKey, false)
	if err != nil {
		return 0, err
	}

	if lastIngestedLedger == "" {
		return 0, nil
	} else {
		ledgerSequence, err := strconv.ParseUint(lastIngestedLedger, 10, 32)
		if err != nil {
			return 0, errors.Wrap(err, "Error converting lastIngestedLedger value")
		}

		return uint32(ledgerSequence), nil
	}
}

// GetLastLedgerIngest returns the last ledger ingested by ingest system
// in Aurora. Returns ErrKeyNotFound error if no value has been previously set.
// This is using `SELECT ... FOR UPDATE` what means it's blocking the row for all other
// transactions.This behaviour is critical in distributed ingestion so do not change
// it unless you know what you are doing.
// The value can be set using UpdateLastLedgerIngest.
func (q *Q) GetLastLedgerIngest(ctx context.Context) (uint32, error) {
	lastIngestedLedger, err := q.getValueFromStore(ctx, lastLedgerKey, true)
	if err != nil {
		return 0, err
	}

	if lastIngestedLedger == "" {
		// This key should always be in a DB (is added in migrations). Otherwise
		// locking won't work.
		return 0, errors.Errorf("`%s` key cannot be found in the key value store", ingestVersion)
	} else {
		ledgerSequence, err := strconv.ParseUint(lastIngestedLedger, 10, 32)
		if err != nil {
			return 0, errors.Wrap(err, "Error converting lastIngestedLedger value")
		}

		return uint32(ledgerSequence), nil
	}
}

// UpdateLastLedgerIngest updates the last ledger ingested by ingest system.
// Can be read using GetLastLedgerExpIngest.
func (q *Q) UpdateLastLedgerIngest(ctx context.Context, ledgerSequence uint32) error {
	return q.updateValueInStore(
		ctx,
		lastLedgerKey,
		strconv.FormatUint(uint64(ledgerSequence), 10),
	)
}

// GetIngestVersion returns the ingestion version. Returns zero
// if there is no value.
func (q *Q) GetIngestVersion(ctx context.Context) (int, error) {
	parsed, err := q.getIntValueFromStore(ctx, ingestVersion, 32)
	if err != nil {
		return 0, errors.Wrap(err, "Error converting sequence value")
	}
	return int(parsed), nil
}

// UpdateIngestVersion updates the ingestion version.
func (q *Q) UpdateIngestVersion(ctx context.Context, version int) error {
	return q.updateValueInStore(
		ctx,
		ingestVersion,
		strconv.FormatUint(uint64(version), 10),
	)
}

// GetExpStateInvalid returns true if the state was found to be invalid.
// Returns false otherwise.
func (q *Q) GetExpStateInvalid(ctx context.Context) (bool, error) {
	invalid, err := q.getValueFromStore(ctx, stateInvalid, false)
	if err != nil {
		return false, err
	}

	if invalid == "" {
		return false, nil
	} else {
		val, err := strconv.ParseBool(invalid)
		if err != nil {
			return false, errors.Wrap(err, "Error converting invalid value")
		}

		return val, nil
	}
}

// UpdateExpStateInvalid updates the state invalid value.
func (q *Q) UpdateExpStateInvalid(ctx context.Context, val bool) error {
	return q.updateValueInStore(
		ctx,
		stateInvalid,
		strconv.FormatBool(val),
	)
}

// GetOfferCompactionSequence returns the sequence number corresponding to the
// last time the offers table was compacted.
func (q *Q) GetOfferCompactionSequence(ctx context.Context) (uint32, error) {
	parsed, err := q.getIntValueFromStore(ctx, offerCompactionSequence, 32)
	if err != nil {
		return 0, errors.Wrap(err, "Error converting sequence value")
	}
	return uint32(parsed), nil
}

// GetLiquidityPoolCompactionSequence returns the sequence number corresponding to the
// last time the liquidity pools table was compacted.
func (q *Q) GetLiquidityPoolCompactionSequence(ctx context.Context) (uint32, error) {
	parsed, err := q.getIntValueFromStore(ctx, liquidityPoolCompactionSequence, 32)
	if err != nil {
		return 0, errors.Wrap(err, "Error converting sequence value")
	}

	return uint32(parsed), nil
}

func (q *Q) getIntValueFromStore(ctx context.Context, key string, bitSize int) (int64, error) {
	sequence, err := q.getValueFromStore(ctx, key, false)
	if err != nil {
		return 0, err
	}

	if sequence == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(sequence, 10, bitSize)
	if err != nil {
		return 0, errors.Wrap(err, "Error converting value")
	}
	return parsed, nil
}

// UpdateOfferCompactionSequence sets the sequence number corresponding to the
// last time the offers table was compacted.
func (q *Q) UpdateOfferCompactionSequence(ctx context.Context, sequence uint32) error {
	return q.updateValueInStore(
		ctx,
		offerCompactionSequence,
		strconv.FormatUint(uint64(sequence), 10),
	)
}

func (q *Q) UpdateLiquidityPoolCompactionSequence(ctx context.Context, sequence uint32) error {
	return q.updateValueInStore(
		ctx,
		liquidityPoolCompactionSequence,
		strconv.FormatUint(uint64(sequence), 10),
	)
}

// getValueFromStore returns a value for a given key from KV store. If value
// is not present in the key value store "" will be returned.
func (q *Q) getValueFromStore(ctx context.Context, key string, forUpdate bool) (string, error) {
	query := sq.Select("key_value_store.value").
		From("key_value_store").
		Where("key_value_store.key = ?", key)

	if forUpdate {
		query = query.Suffix("FOR UPDATE")
	}

	var value string
	if err := q.Get(ctx, &value, query); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return "", nil
		}
		return "", errors.Wrap(err, "could not get value")
	}

	return value, nil
}

// updateValueInStore updates a value for a given key in KV store
func (q *Q) updateValueInStore(ctx context.Context, key, value string) error {
	query := sq.Insert("key_value_store").
		Columns("key", "value").
		Values(key, value).
		Suffix("ON CONFLICT (key) DO UPDATE SET value=EXCLUDED.value")

	_, err := q.Exec(ctx, query)
	return err
}

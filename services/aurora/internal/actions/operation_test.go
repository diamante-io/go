package actions

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diamcircle/go/protocols/aurora/operations"
	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	"github.com/diamcircle/go/services/aurora/internal/render/problem"
	"github.com/diamcircle/go/services/aurora/internal/test"
	supportProblem "github.com/diamcircle/go/support/render/problem"
)

func TestGetOperationsWithoutFilter(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("base")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{}

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 4)
}

func TestGetOperationsExclusiveFilters(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("base")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{}

	testCases := []struct {
		desc  string
		query map[string]string
	}{
		{
			desc: "tx_id & ledger_id",
			query: map[string]string{
				"tx_id":     "1d2a4be72470658f68db50eef29ea0af3f985ce18b5c218f03461d40c47dc292",
				"ledger_id": "1",
			},
		},
		{
			desc: "tx_id & account_id",
			query: map[string]string{
				"tx_id":      "1d2a4be72470658f68db50eef29ea0af3f985ce18b5c218f03461d40c47dc292",
				"account_id": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			},
		},
		{
			desc: "account_id & ledger_id",
			query: map[string]string{
				"account_id": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
				"ledger_id":  "1",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := handler.GetResourcePage(
				httptest.NewRecorder(),
				makeRequest(
					t, tc.query, map[string]string{}, q,
				),
			)
			tt.Assert.IsType(&supportProblem.P{}, err)
			p := err.(*supportProblem.P)
			tt.Assert.Equal("bad_request", p.Type)
			tt.Assert.Equal("filters", p.Extras["invalid_field"])
			tt.Assert.Equal(
				"Use a single filter for operations, you can only use one of tx_id, account_id or ledger_id",
				p.Extras["reason"],
			)
		})
	}

}

func TestGetOperationsByLiquidityPool(t *testing.T) {

}

func TestGetOperationsFilterByAccountID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("base")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{}

	testCases := []struct {
		accountID string
		expected  int
	}{
		{
			accountID: "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
			expected:  3,
		},
		{
			accountID: "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			expected:  1,
		},
		{
			accountID: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU",
			expected:  2,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("account %s operations", tc.accountID), func(t *testing.T) {
			records, err := handler.GetResourcePage(
				httptest.NewRecorder(),
				makeRequest(
					t, map[string]string{
						"account_id": tc.accountID,
					}, map[string]string{}, q,
				),
			)
			tt.Assert.NoError(err)
			tt.Assert.Len(records, tc.expected)
		})
	}
}

func TestGetOperationsFilterByTxID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("base")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{}

	testCases := []struct {
		desc          string
		transactionID string
		expected      int
		expectedErr   string
		notFound      bool
	}{
		{
			desc:          "operations for 2374...6d4d",
			transactionID: "2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d",
			expected:      1,
		},
		{
			desc:          "operations for 164a...33b6",
			transactionID: "164a5064eba64f2cdbadb856bf3448485fc626247ada3ed39cddf0f6902133b6",
			expected:      1,
		},
		{
			desc:          "missing transaction",
			transactionID: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			expectedErr:   "sql: no rows in result set",
			notFound:      true,
		},
		{
			desc:          "uppercase tx hash not accepted",
			transactionID: "2374E99349B9EF7DBA9A5DB3339B78FDA8F34777B1AF33BA468AD5C0DF946D4D",
			expectedErr:   "Transaction hash must be a hex-encoded, lowercase SHA-256 hash",
		},
		{
			desc:          "badly formated tx hash not accepted",
			transactionID: "%00%1E4%5E%EF%BF%BD%EF%BF%BD%EF%BF%BDpVP%EF%BF%BDI&R%0BK%EF%BF%BD%1D%EF%BF%BD%EF%BF%BD=%EF%BF%BD%3F%23%EF%BF%BD%EF%BF%BDl%EF%BF%BD%1El%EF%BF%BD%EF%BF%BD",
			expectedErr:   "Transaction hash must be a hex-encoded, lowercase SHA-256 hash",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf(tc.desc), func(t *testing.T) {
			records, err := handler.GetResourcePage(
				httptest.NewRecorder(),
				makeRequest(
					t, map[string]string{
						"tx_id": tc.transactionID,
					}, map[string]string{}, q,
				),
			)

			if tc.expectedErr == "" {
				tt.Assert.NoError(err)
				tt.Assert.Len(records, tc.expected)
			} else {
				if tc.notFound {
					tt.Assert.EqualError(err, tc.expectedErr)
				} else {
					tt.Assert.IsType(&supportProblem.P{}, err)
					p := err.(*supportProblem.P)
					tt.Assert.Equal("bad_request", p.Type)
					tt.Assert.Equal("tx_id", p.Extras["invalid_field"])
					tt.Assert.Equal(
						tc.expectedErr,
						p.Extras["reason"],
					)
				}
			}
		})
	}
}

func TestGetOperationsIncludeFailed(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("failed_transactions")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{}

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"limit": "200",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)

	successful := 0
	failed := 0

	for _, record := range records {
		op := record.(operations.Operation)
		if op.IsTransactionSuccessful() {
			successful++
		} else {
			failed++
		}
	}

	tt.Assert.Equal(8, successful)
	tt.Assert.Equal(0, failed)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"include_failed": "true",
				"limit":          "200",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)

	successful = 0
	failed = 0

	for _, record := range records {
		op := record.(operations.Operation)
		if op.IsTransactionSuccessful() {
			successful++
		} else {
			failed++
		}
	}

	tt.Assert.Equal(8, successful)
	tt.Assert.Equal(1, failed)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"tx_id": "aa168f12124b7c196c0adaee7c73a64d37f99428cacb59a91ff389626845e7cf",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 1)
	for _, record := range records {
		op := record.(operations.Operation)
		tt.Assert.False(op.IsTransactionSuccessful())
	}

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"tx_id": "56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 1)
	for _, record := range records {
		op := record.(operations.Operation)
		tt.Assert.True(op.IsTransactionSuccessful())
	}

	// NULL value
	_, err = tt.AuroraSession().ExecRaw(tt.Ctx,
		`UPDATE history_transactions SET successful = NULL WHERE transaction_hash = ?`,
		"56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1",
	)
	tt.Assert.NoError(err)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"tx_id": "56e3216045d579bea40f2d35a09406de3a894ecb5be70dbda5ec9c0427a0d5a1",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 1)
	for _, record := range records {
		op := record.(operations.Operation)
		tt.Assert.True(op.IsTransactionSuccessful())
	}

	_, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"include_failed": "foo",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.Error(err)
	tt.Assert.IsType(&supportProblem.P{}, err)
	p := err.(*supportProblem.P)
	tt.Assert.Equal("bad_request", p.Type)
	tt.Assert.Equal("include_failed", p.Extras["invalid_field"])
	tt.Assert.Equal(
		"Filter should be true or false",
		p.Extras["reason"],
	)
}

func TestGetOperationsFilterByLedgerID(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("base")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{}

	testCases := []struct {
		ledgerID    string
		expected    int
		expectedErr string
		notFound    bool
	}{
		{
			ledgerID: "1",
			expected: 0,
		},
		{
			ledgerID: "2",
			expected: 3,
		},
		{
			ledgerID: "3",
			expected: 1,
		},
		{
			ledgerID:    "10000",
			expectedErr: "sql: no rows in result set",
			notFound:    true,
		},
		{
			ledgerID:    "-1",
			expectedErr: "Ledger ID must be an integer higher than 0",
		},
		{
			ledgerID:    "one",
			expectedErr: "Ledger ID must be an integer higher than 0",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("ledger %s operations", tc.ledgerID), func(t *testing.T) {
			records, err := handler.GetResourcePage(
				httptest.NewRecorder(),
				makeRequest(
					t, map[string]string{
						"ledger_id": tc.ledgerID,
					}, map[string]string{}, q,
				),
			)
			if tc.expectedErr == "" {
				tt.Assert.NoError(err)
				tt.Assert.Len(records, tc.expected)
			} else {
				if tc.notFound {
					tt.Assert.EqualError(err, tc.expectedErr)
				} else {
					tt.Assert.IsType(&supportProblem.P{}, err)
					p := err.(*supportProblem.P)
					tt.Assert.Equal("bad_request", p.Type)
					tt.Assert.Equal("ledger_id", p.Extras["invalid_field"])
					tt.Assert.Equal(
						tc.expectedErr,
						p.Extras["reason"],
					)
				}
			}
		})
	}
}

func TestGetOperationsOnlyPayments(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("base")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{
		OnlyPayments: true,
	}

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 4)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"ledger_id": "1",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 0)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"ledger_id": "3",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 1)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"account_id": "GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 1)

	tt.Scenario("pathed_payment")

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"tx_id": "b52f16ffb98c047e33b9c2ec30880330cde71f85b3443dae2c5cb86c7d4d8452",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 0)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"tx_id": "1d2a4be72470658f68db50eef29ea0af3f985ce18b5c218f03461d40c47dc292",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 1)

	record := records[0].(operations.PathPayment)
	tt.Assert.Equal("10.0000000", record.SourceAmount)
}

func TestOperation_CreatedAt(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("base")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{}

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"ledger_id": "3",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)

	l := history.Ledger{}
	tt.Assert.NoError(q.LedgerBySequence(tt.Ctx, &l, 3))

	record := records[0].(operations.Payment)

	tt.Assert.WithinDuration(l.ClosedAt, record.LedgerCloseTime, 1*time.Second)
}
func TestGetOperationsPagination(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("base")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{
		LedgerState: &ledger.State{},
	}

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"order": "asc",
				"limit": "1",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 1)

	descRecords, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"limit": "1",
				"order": "desc",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.NotEqual(records, descRecords)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"order":  "desc",
				"cursor": "12884905985",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	tt.Assert.Len(records, 3)

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"order":  "desc",
				"cursor": "0",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.Error(err)
	tt.Assert.EqualError(err, "problem: before_history")
}

func TestGetOperations_IncludeTransactions(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("failed_transactions")

	q := &history.Q{tt.AuroraSession()}
	handler := GetOperationsHandler{}

	_, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"join": "accounts",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.Error(err)
	tt.Assert.IsType(&supportProblem.P{}, err)
	p := err.(*supportProblem.P)
	tt.Assert.Equal("bad_request", p.Type)
	tt.Assert.Equal("join", p.Extras["invalid_field"])
	tt.Assert.Equal(
		"Accepted values: transactions",
		p.Extras["reason"],
	)

	records, err := handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"join":  "transactions",
				"limit": "1",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	for _, record := range records {
		op := record.(operations.CreateAccount)
		tt.Assert.NotNil(op.Transaction)
	}

	records, err = handler.GetResourcePage(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{
				"limit": "1",
			}, map[string]string{}, q,
		),
	)
	tt.Assert.NoError(err)
	for _, record := range records {
		op := record.(operations.CreateAccount)
		tt.Assert.Nil(op.Transaction)
	}
}
func TestGetOperation(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()

	handler := GetOperationByIDHandler{
		LedgerState: &ledger.State{},
	}
	handler.LedgerState.SetStatus(tt.Scenario("base"))

	record, err := handler.GetResource(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{}, map[string]string{"id": "8589938689"}, tt.AuroraSession(),
		),
	)
	tt.Assert.NoError(err)
	op := record.(operations.Operation)
	tt.Assert.Equal("8589938689", op.PagingToken())
	tt.Assert.Equal("2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d", op.GetTransactionHash())

	_, err = handler.GetResource(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{}, map[string]string{"id": "9589938689"}, tt.AuroraSession(),
		),
	)

	tt.Assert.Equal(err, sql.ErrNoRows)

	_, err = handler.GetResource(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{}, map[string]string{"id": "0"}, tt.AuroraSession(),
		),
	)
	tt.Assert.Equal(err, problem.BeforeHistory)
}

func TestOperation_IncludeTransaction(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	tt.Scenario("kahuna")

	handler := GetOperationByIDHandler{
		LedgerState: &ledger.State{},
	}
	record, err := handler.GetResource(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{}, map[string]string{"id": "261993009153"}, tt.AuroraSession(),
		),
	)

	tt.Assert.NoError(err)

	op := record.(operations.BumpSequence)
	tt.Assert.Nil(op.Transaction)

	record, err = handler.GetResource(
		httptest.NewRecorder(),
		makeRequest(
			t, map[string]string{"join": "transactions"}, map[string]string{"id": "261993009153"}, tt.AuroraSession(),
		),
	)
	op = record.(operations.BumpSequence)
	tt.Assert.NotNil(op.Transaction)
	tt.Assert.Equal(op.TransactionHash, op.Transaction.ID)
}

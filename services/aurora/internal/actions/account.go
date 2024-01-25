package actions

import (
	"context"
	"net/http"
	"strings"

	protocol "go/protocols/aurora"
	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2/history"
	"go/services/aurora/internal/ledger"
	"go/services/aurora/internal/resourceadapter"
	"go/support/errors"
	"go/support/render/hal"
	"go/support/render/problem"
	"go/xdr"
)

// AccountInfo returns the information about an account identified by addr.
func AccountInfo(ctx context.Context, hq *history.Q, addr string) (*protocol.Account, error) {
	var (
		record     history.AccountEntry
		data       []history.Data
		signers    []history.AccountSigner
		trustlines []history.TrustLine
		resouce    protocol.Account
	)

	record, err := hq.GetAccountByID(ctx, addr)
	if err != nil {
		return nil, errors.Wrap(err, "getting history account record")
	}

	data, err = hq.GetAccountDataByAccountID(ctx, addr)
	if err != nil {
		return nil, errors.Wrap(err, "getting history account data")
	}

	signers, err = hq.GetAccountSignersByAccountID(ctx, addr)
	if err != nil {
		return nil, errors.Wrap(err, "getting history signers")
	}

	trustlines, err = hq.GetSortedTrustLinesByAccountID(ctx, addr)
	if err != nil {
		return nil, errors.Wrap(err, "getting history trustlines")
	}

	ledger, err := getLedgerBySequence(ctx, hq, int32(record.LastModifiedLedger))
	if err != nil {
		return nil, err
	}

	err = resourceadapter.PopulateAccountEntry(
		ctx,
		&resouce,
		record,
		data,
		signers,
		trustlines,
		ledger,
	)
	if err != nil {
		return nil, errors.Wrap(err, "populating account entry")
	}

	return &resouce, nil
}

// AccountsQuery query struct for accounts end-point
type AccountsQuery struct {
	Signer        string `schema:"signer" valid:"accountID,optional"`
	Sponsor       string `schema:"sponsor" valid:"accountID,optional"`
	AssetFilter   string `schema:"asset" valid:"asset,optional"`
	LiquidityPool string `schema:"liquidity_pool" valid:"sha256,optional"`
}

// URITemplate returns a rfc6570 URI template the query struct
func (q AccountsQuery) URITemplate() string {
	return getURITemplate(&q, "accounts", true)
}

var invalidAccountsParams = problem.P{
	Type:   "invalid_accounts_params",
	Title:  "Invalid Accounts Parameters",
	Status: http.StatusBadRequest,
	Detail: "Exactly one filter is required. Please ensure that you are including a signer, sponsor, asset, or liquidity pool filter.",
}

// Validate runs custom validations.
func (q AccountsQuery) Validate() error {
	if q.AssetFilter == "native" {
		return problem.MakeInvalidFieldProblem(
			"asset",
			errors.New("you can't filter by asset: native"),
		)
	}

	numParams, err := countNonEmpty(q.Sponsor, q.Signer, q.Asset(), q.LiquidityPool)
	if err != nil {
		return errors.Wrap(err, "Could not count request params")
	}
	if numParams != 1 {
		return invalidAccountsParams
	}

	return nil
}

// Asset returns an xdr.Asset representing the Asset we want to find the trustees by.
func (q AccountsQuery) Asset() *xdr.Asset {
	if len(q.AssetFilter) == 0 {
		return nil
	}

	parts := strings.Split(q.AssetFilter, ":")
	asset := xdr.MustNewCreditAsset(parts[0], parts[1])

	return &asset
}

// GetAccountsHandler is the action handler for the /accounts endpoint
type GetAccountsHandler struct {
	LedgerState *ledger.State
}

// GetResourcePage returns a page containing the account records that have
// `signer` as a signer, `sponsor` as a sponsor, a trustline to the given
// `asset`, or participate in a particular `liquidity_pool`.
func (handler GetAccountsHandler) GetResourcePage(
	w HeaderWriter,
	r *http.Request,
) ([]hal.Pageable, error) {
	ctx := r.Context()
	pq, err := GetPageQuery(handler.LedgerState, r, DisableCursorValidation)
	if err != nil {
		return nil, err
	}

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}

	qp := AccountsQuery{}
	err = getParams(&qp, r)
	if err != nil {
		return nil, err
	}

	var records []history.AccountEntry

	if len(qp.Sponsor) > 0 {
		records, err = historyQ.AccountsForSponsor(ctx, qp.Sponsor, pq)
		if err != nil {
			return nil, errors.Wrap(err, "loading account records")
		}
	} else if len(qp.Signer) > 0 {
		records, err = historyQ.AccountEntriesForSigner(ctx, qp.Signer, pq)
		if err != nil {
			return nil, errors.Wrap(err, "loading account records")
		}
	} else if len(qp.LiquidityPool) > 0 {
		records, err = historyQ.AccountsForLiquidityPool(ctx, qp.LiquidityPool, pq)
		if err != nil {
			return nil, errors.Wrap(err, "loading account records")
		}
	} else {
		records, err = historyQ.AccountsForAsset(ctx, *qp.Asset(), pq)
		if err != nil {
			return nil, errors.Wrap(err, "loading account records")
		}
	}

	accounts := make([]hal.Pageable, 0, len(records))

	if len(records) == 0 {
		// early return
		return accounts, nil
	}

	accountIDs := make([]string, 0, len(records))
	for _, record := range records {
		accountIDs = append(accountIDs, record.AccountID)
	}

	signers, err := handler.loadSigners(ctx, historyQ, accountIDs)
	if err != nil {
		return nil, err
	}

	trustlines, err := handler.loadTrustlines(ctx, historyQ, accountIDs)
	if err != nil {
		return nil, err
	}

	data, err := handler.loadData(ctx, historyQ, accountIDs)
	if err != nil {
		return nil, err
	}

	ledgerCache := history.LedgerCache{}
	for _, record := range records {
		ledgerCache.Queue(int32(record.LastModifiedLedger))
	}

	if err := ledgerCache.Load(ctx, historyQ); err != nil {
		return nil, errors.Wrap(err, "failed to load ledger batch")
	}

	for _, record := range records {
		var res protocol.Account
		s := signers[record.AccountID]
		t := trustlines[record.AccountID]
		d := data[record.AccountID]
		var ledger *history.Ledger
		if l, ok := ledgerCache.Records[int32(record.LastModifiedLedger)]; ok {
			ledger = &l
		}
		resourceadapter.PopulateAccountEntry(ctx, &res, record, d, s, t, ledger)

		accounts = append(accounts, res)
	}

	return accounts, nil
}

func (handler GetAccountsHandler) loadData(ctx context.Context, historyQ *history.Q, accounts []string) (map[string][]history.Data, error) {
	data := make(map[string][]history.Data)

	records, err := historyQ.GetAccountDataByAccountsID(ctx, accounts)
	if err != nil {
		return data, errors.Wrap(err, "loading account data records by accounts id")
	}

	for _, record := range records {
		data[record.AccountID] = append(data[record.AccountID], record)
	}

	return data, nil
}

func (handler GetAccountsHandler) loadTrustlines(ctx context.Context, historyQ *history.Q, accounts []string) (map[string][]history.TrustLine, error) {
	trustLines := make(map[string][]history.TrustLine)

	records, err := historyQ.GetSortedTrustLinesByAccountIDs(ctx, accounts)
	if err != nil {
		return trustLines, errors.Wrap(err, "loading trustline records by accounts")
	}

	for _, record := range records {
		trustLines[record.AccountID] = append(trustLines[record.AccountID], record)
	}

	return trustLines, nil
}

func (handler GetAccountsHandler) loadSigners(ctx context.Context, historyQ *history.Q, accounts []string) (map[string][]history.AccountSigner, error) {
	signers := make(map[string][]history.AccountSigner)

	records, err := historyQ.SignersForAccounts(ctx, accounts)
	if err != nil {
		return signers, errors.Wrap(err, "loading account signers by account")
	}

	for _, record := range records {
		signers[record.Account] = append(signers[record.Account], record)
	}

	return signers, nil
}

func getLedgerBySequence(ctx context.Context, hq *history.Q, sequence int32) (*history.Ledger, error) {
	ledger := &history.Ledger{}
	err := hq.LedgerBySequence(ctx, ledger, sequence)
	switch {
	case hq.NoRows(err):
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return ledger, nil
	}
}

// AccountByIDQuery query struct for accounts/{account_id} end-point
type AccountByIDQuery struct {
	AccountID string `schema:"account_id" valid:"accountID,optional"`
}

// GetAccountByIDHandler is the action handler for the /accounts/{account_id} endpoint
type GetAccountByIDHandler struct{}

type Account protocol.Account

func (a Account) Equals(other StreamableObjectResponse) bool {
	otherAccount, ok := other.(Account)
	if !ok {
		return false
	}
	return a.ID == otherAccount.ID
}

func (handler GetAccountByIDHandler) GetResource(
	w HeaderWriter,
	r *http.Request,
) (StreamableObjectResponse, error) {

	historyQ, err := auroraContext.HistoryQFromRequest(r)
	if err != nil {
		return nil, err
	}

	qp := AccountByIDQuery{}
	err = getParams(&qp, r)
	if err != nil {
		return nil, err
	}
	account, err := AccountInfo(r.Context(), historyQ, qp.AccountID)
	if err != nil {
		return Account{}, err
	}
	return Account(*account), nil
}

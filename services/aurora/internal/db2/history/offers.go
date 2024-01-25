package history

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"go/support/errors"
)

// QOffers defines offer related queries.
type QOffers interface {
	GetAllOffers(ctx context.Context) ([]Offer, error)
	GetOffersByIDs(ctx context.Context, ids []int64) ([]Offer, error)
	CountOffers(ctx context.Context) (int, error)
	GetUpdatedOffers(ctx context.Context, newerThanSequence uint32) ([]Offer, error)
	UpsertOffers(ctx context.Context, offers []Offer) error
	CompactOffers(ctx context.Context, cutOffSequence uint32) (int64, error)
}

func (q *Q) CountOffers(ctx context.Context) (int, error) {
	sql := sq.Select("count(*)").Where("deleted = ?", false).From("offers")

	var count int
	if err := q.Get(ctx, &count, sql); err != nil {
		return 0, errors.Wrap(err, "could not run select query")
	}

	return count, nil
}

// GetOfferByID loads a row from the `offers` table, selected by offerid.
func (q *Q) GetOfferByID(ctx context.Context, id int64) (Offer, error) {
	var offer Offer
	sql := selectOffers.Where("deleted = ?", false).
		Where("offers.offer_id = ?", id)
	err := q.Get(ctx, &offer, sql)
	return offer, err
}

// GetOffersByIDs loads a row from the `offers` table, selected by multiple offerid.
func (q *Q) GetOffersByIDs(ctx context.Context, ids []int64) ([]Offer, error) {
	var offers []Offer
	sql := selectOffers.Where("deleted = ?", false).
		Where(map[string]interface{}{"offers.offer_id": ids})
	err := q.Select(ctx, &offers, sql)
	return offers, err
}

// GetOffers loads rows from `offers` by paging query.
func (q *Q) GetOffers(ctx context.Context, query OffersQuery) ([]Offer, error) {
	sql := selectOffers.Where("deleted = ?", false)
	sql, err := query.PageQuery.ApplyTo(sql, "offers.offer_id")

	if err != nil {
		return nil, errors.Wrap(err, "could not apply query to page")
	}

	if query.SellerID != "" {
		sql = sql.Where("offers.seller_id = ?", query.SellerID)
	}

	if query.Selling != nil {
		sql = sql.Where("offers.selling_asset = ?", query.Selling)
	}

	if query.Buying != nil {
		sql = sql.Where("offers.buying_asset = ?", query.Buying)
	}

	if query.Sponsor != "" {
		sql = sql.Where("offers.sponsor = ?", query.Sponsor)
	}

	var offers []Offer
	if err := q.Select(ctx, &offers, sql); err != nil {
		return nil, errors.Wrap(err, "could not run select query")
	}

	return offers, nil
}

// GetAllOffers loads all non deleted offers
func (q *Q) GetAllOffers(ctx context.Context) ([]Offer, error) {
	var offers []Offer
	err := q.Select(ctx, &offers, selectOffers.Where("deleted = ?", false))
	return offers, err
}

// GetUpdatedOffers returns all offers created, updated, or deleted after the given ledger sequence.
func (q *Q) GetUpdatedOffers(ctx context.Context, newerThanSequence uint32) ([]Offer, error) {
	var offers []Offer
	err := q.Select(ctx, &offers, selectOffers.Where("offers.last_modified_ledger > ?", newerThanSequence))
	return offers, err
}

// UpsertOffers upserts a batch of offers in the offers table.
// There's currently no limit of the number of offers this method can
// accept other than 2GB limit of the query string length what should be enough
// for each ledger with the current limits.
func (q *Q) UpsertOffers(ctx context.Context, offers []Offer) error {
	var sellerID, sellingAsset, buyingAsset, offerID, amount, priceN, priceD,
		price, flags, lastModifiedLedger, deleted, sponsor []interface{}

	for _, offer := range offers {
		sellerID = append(sellerID, offer.SellerID)
		offerID = append(offerID, offer.OfferID)
		sellingAsset = append(sellingAsset, offer.SellingAsset)
		buyingAsset = append(buyingAsset, offer.BuyingAsset)
		amount = append(amount, offer.Amount)
		priceN = append(priceN, offer.Pricen)
		priceD = append(priceD, offer.Priced)
		price = append(price, offer.Price)
		flags = append(flags, offer.Flags)
		lastModifiedLedger = append(lastModifiedLedger, offer.LastModifiedLedger)
		deleted = append(deleted, offer.Deleted)
		sponsor = append(sponsor, offer.Sponsor)
	}

	upsertFields := []upsertField{
		{"seller_id", "text", sellerID},
		{"offer_id", "bigint", offerID},
		{"selling_asset", "text", sellingAsset},
		{"buying_asset", "text", buyingAsset},
		{"amount", "bigint", amount},
		{"pricen", "integer", priceN},
		{"priced", "integer", priceD},
		{"price", "double precision", price},
		{"flags", "integer", flags},
		{"deleted", "bool", deleted},
		{"last_modified_ledger", "integer", lastModifiedLedger},
		{"sponsor", "text", sponsor},
	}

	return q.upsertRows(ctx, "offers", "offer_id", upsertFields)
}

// CompactOffers removes rows from the offers table which are marked for deletion.
func (q *Q) CompactOffers(ctx context.Context, cutOffSequence uint32) (int64, error) {
	sql := sq.Delete("offers").
		Where("deleted = ?", true).
		Where("last_modified_ledger <= ?", cutOffSequence)

	result, err := q.Exec(ctx, sql)
	if err != nil {
		return 0, errors.Wrap(err, "cannot delete offer rows")
	}

	if err = q.UpdateOfferCompactionSequence(ctx, cutOffSequence); err != nil {
		return 0, errors.Wrap(err, "cannot update offer compaction sequence")
	}

	return result.RowsAffected()
}

var selectOffers = sq.Select(`
	seller_id,
	offer_id,
	selling_asset,
	buying_asset,
	amount,
	pricen,
	priced,
	price,
	flags,
	deleted,
	last_modified_ledger,
	sponsor
`).From("offers")

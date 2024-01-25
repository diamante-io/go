package paths

import (
	"context"

	"go/xdr"
)

// Query is a query for paths
type Query struct {
	DestinationAsset    xdr.Asset
	DestinationAmount   xdr.Int64
	SourceAssets        []xdr.Asset
	SourceAssetBalances []xdr.Int64
	// if ValidateSourceBalance is true then we won't consider payment paths
	// which require a source asset amount which exceeds the balance present in `SourceAssetBalances`
	ValidateSourceBalance bool
	SourceAccount         *xdr.AccountId
}

// Path is the result returned by a path finder and is tied to the DestinationAmount used in the input query
type Path struct {
	Path              []string
	Source            string
	SourceAmount      xdr.Int64
	Destination       string
	DestinationAmount xdr.Int64
}

// Finder finds paths.
type Finder interface {
	// Find returns a list of payment paths and the most recent ledger
	// for a Query of a maximum length `maxLength`. The payment paths
	// are accurate and consistent with the returned ledger sequence number
	Find(ctx context.Context, q Query, maxLength uint) ([]Path, uint32, error)
	// FindFixedPaths return a list of payment paths the most recent ledger
	// Each of the payment paths start by spending `amountToSpend` of `sourceAsset` and end
	// with delivering a postive amount of `destinationAsset`.
	// The payment paths are accurate and consistent with the returned ledger sequence number
	FindFixedPaths(
		ctx context.Context,
		sourceAsset xdr.Asset,
		amountToSpend xdr.Int64,
		destinationAssets []xdr.Asset,
		maxLength uint,
	) ([]Path, uint32, error)
}

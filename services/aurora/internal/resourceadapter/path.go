package resourceadapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/aurora/internal/paths"
)

func extractAsset(asset string, t, c, i *string) error {
	if asset == "native" {
		*t = asset
		return nil
	}
	parts := strings.Split(asset, "/")
	if len(parts) != 3 {
		return fmt.Errorf("expected length to be 3 but got %v", parts)
	}
	*t = parts[0]
	*c = parts[1]
	*i = parts[2]
	return nil
}

// PopulatePath converts the paths.Path into a Path
func PopulatePath(ctx context.Context, dest *aurora.Path, p paths.Path) (err error) {
	dest.DestinationAmount = amount.String(p.DestinationAmount)
	dest.SourceAmount = amount.String(p.SourceAmount)

	err = extractAsset(
		p.Source,
		&dest.SourceAssetType,
		&dest.SourceAssetCode,
		&dest.SourceAssetIssuer)
	if err != nil {
		return
	}

	err = extractAsset(
		p.Destination,
		&dest.DestinationAssetType,
		&dest.DestinationAssetCode,
		&dest.DestinationAssetIssuer)
	if err != nil {
		return
	}

	dest.Path = make([]aurora.Asset, len(p.Path))
	for i, a := range p.Path {
		err = extractAsset(
			a,
			&dest.Path[i].Type,
			&dest.Path[i].Code,
			&dest.Path[i].Issuer)
		if err != nil {
			return
		}
	}
	return
}

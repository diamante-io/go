package tickerdb

import (
	"net/http"
	"os"
	"strings"
	"testing"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/shurcooL/httpfs/filter"

	bdata "github.com/diamcircle/go/services/ticker/internal/tickerdb/migrations"
	supportHttp "github.com/diamcircle/go/support/http"
)

func TestGeneratedAssets(t *testing.T) {
	var localAssets http.FileSystem = filter.Keep(http.Dir("migrations"), func(path string, fi os.FileInfo) bool {
		return fi.IsDir() || strings.HasSuffix(path, ".sql")
	})
	generatedAssets := &assetfs.AssetFS{
		Asset:     bdata.Asset,
		AssetDir:  bdata.AssetDir,
		AssetInfo: bdata.AssetInfo,
		Prefix:    "/migrations",
	}

	if !supportHttp.EqualFileSystems(localAssets, generatedAssets, "/") {
		t.Fatalf("generated migrations does not match local migrations")
	}
}

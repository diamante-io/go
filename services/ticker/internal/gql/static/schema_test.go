package static

import (
	"net/http"
	"os"
	"strings"
	"testing"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/shurcooL/httpfs/filter"

	supportHttp "go/support/http"
)

func TestGeneratedAssets(t *testing.T) {
	var localAssets http.FileSystem = filter.Skip(http.Dir("."), func(path string, fi os.FileInfo) bool {
		return !fi.IsDir() && strings.HasSuffix(path, ".go")
	})
	generatedAssets := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
	}

	if !supportHttp.EqualFileSystems(localAssets, generatedAssets, "/") {
		t.Fatalf("generated migrations does not match local migrations")
	}
}

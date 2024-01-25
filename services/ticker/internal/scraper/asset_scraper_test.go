package scraper

import (
	"net/url"
	"testing"

	hProtocol "github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/render/hal"
	"github.com/stretchr/testify/assert"
)

func TestShouldDiscardAsset(t *testing.T) {
	testAsset := hProtocol.AssetStat{
		Amount: "",
	}

	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount: "0.0",
	}
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount: "0",
	}
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 8,
	}
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 12,
	}
	testAsset.Code = "REMOVE"
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 100,
	}
	testAsset.Code = "SOMETHINGVALID"
	testAsset.Links.Toml.Href = ""
	assert.Equal(t, shouldDiscardAsset(testAsset, true), false)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 40,
	}
	testAsset.Code = "SOMETHINGVALID"
	testAsset.Links.Toml.Href = "http://www.diamcircle.org/.well-known/diamcircle.toml"
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 40,
	}
	testAsset.Code = "SOMETHINGVALID"
	testAsset.Links.Toml.Href = ""
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 40,
	}
	testAsset.Code = "SOMETHINGVALID"
	testAsset.Links.Toml.Href = "https://www.diamcircle.org/.well-known/diamcircle.toml"
	assert.Equal(t, shouldDiscardAsset(testAsset, true), false)
}

func TestDomainsMatch(t *testing.T) {
	tomlURL, _ := url.Parse("https://diamcircle.org/diamcircle.toml")
	orgURL, _ := url.Parse("https://diamcircle.org/")
	assert.True(t, domainsMatch(tomlURL, orgURL))

	tomlURL, _ = url.Parse("https://assets.diamcircle.org/diamcircle.toml")
	orgURL, _ = url.Parse("https://diamcircle.org/")
	assert.False(t, domainsMatch(tomlURL, orgURL))

	tomlURL, _ = url.Parse("https://diamcircle.org/diamcircle.toml")
	orgURL, _ = url.Parse("https://home.diamcircle.org/")
	assert.True(t, domainsMatch(tomlURL, orgURL))

	tomlURL, _ = url.Parse("https://diamcircle.org/diamcircle.toml")
	orgURL, _ = url.Parse("https://home.diamcircle.com/")
	assert.False(t, domainsMatch(tomlURL, orgURL))

	tomlURL, _ = url.Parse("https://diamcircle.org/diamcircle.toml")
	orgURL, _ = url.Parse("https://diamcircle.com/")
	assert.False(t, domainsMatch(tomlURL, orgURL))
}

func TestIsDomainVerified(t *testing.T) {
	tomlURL := "https://diamcircle.org/diamcircle.toml"
	orgURL := "https://diamcircle.org/"
	hasCurrency := true
	assert.True(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "https://diamcircle.org/diamcircle.toml"
	orgURL = ""
	hasCurrency = true
	assert.True(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = ""
	orgURL = ""
	hasCurrency = true
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "https://diamcircle.org/diamcircle.toml"
	orgURL = "https://diamcircle.org/"
	hasCurrency = false
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "http://diamcircle.org/diamcircle.toml"
	orgURL = "https://diamcircle.org/"
	hasCurrency = true
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "https://diamcircle.org/diamcircle.toml"
	orgURL = "http://diamcircle.org/"
	hasCurrency = true
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "https://diamcircle.org/diamcircle.toml"
	orgURL = "https://diamcircle.com/"
	hasCurrency = true
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))
}

func TestIgnoreInvalidTOMLUrls(t *testing.T) {
	invalidURL := "https:// there is something wrong here.com/diamcircle.toml"
	assetStat := hProtocol.AssetStat{}
	assetStat.Links.Toml = hal.Link{Href: invalidURL}

	_, err := fetchTOMLData(assetStat)

	urlErr, ok := errors.Cause(err).(*url.Error)
	if !ok {
		t.Fatalf("err expected to be a url.Error but was %#v", err)
	}
	assert.Equal(t, "parse", urlErr.Op)
	assert.Equal(t, "https:// there is something wrong here.com/diamcircle.toml", urlErr.URL)
	assert.EqualError(t, urlErr.Err, `invalid character " " in host name`)
}

package serve

import (
	"net/http"
	"testing"
	"time"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/stretchr/testify/require"
)

func TestAuroraClient(t *testing.T) {
	opts := Options{AuroraURL: "my-aurora.domain.com"}
	auroraClientInterface := opts.auroraClient()

	auroraClient, ok := auroraClientInterface.(*auroraclient.Client)
	require.True(t, ok)
	require.Equal(t, "my-aurora.domain.com", auroraClient.AuroraURL)

	httpClient, ok := auroraClient.HTTP.(*http.Client)
	require.True(t, ok)
	require.Equal(t, http.Client{Timeout: 30 * time.Second}, *httpClient)
}

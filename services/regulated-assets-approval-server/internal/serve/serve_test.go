package serve

import (
	"net/http"
	"testing"
	"time"

	"go/clients/auroraclient"

	"github.com/stretchr/testify/require"
)

func TestauroraClient(t *testing.T) {
	opts := Options{auroraURL: "my-aurora.domain.com"}
	auroraClientInterface := opts.auroraClient()

	auroraClient, ok := auroraClientInterface.(*auroraclient.Client)
	require.True(t, ok)
	require.Equal(t, "my-aurora.domain.com", auroraClient.auroraURL)

	httpClient, ok := auroraClient.HTTP.(*http.Client)
	require.True(t, ok)
	require.Equal(t, http.Client{Timeout: 30 * time.Second}, *httpClient)
}

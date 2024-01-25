package ledgerbackend

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/diamcircle/go/support/log"
)

func TestCloseBeforeStartOffline(t *testing.T) {
	captiveCoreToml, err := NewCaptiveCoreToml(CaptiveCoreTomlParams{})
	assert.NoError(t, err)

	runner, err := newDiamcircleCoreRunner(CaptiveCoreConfig{
		HistoryArchiveURLs: []string{"http://localhost"},
		Log:                log.New(),
		Context:            context.Background(),
		Toml:               captiveCoreToml,
	}, diamcircleCoreRunnerModeOffline)
	assert.NoError(t, err)

	tempDir := runner.storagePath
	info, err := os.Stat(tempDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	assert.NoError(t, runner.close())

	// Directory cleaned up on shutdown when reingesting to save space
	_, err = os.Stat(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestCloseBeforeStartOnline(t *testing.T) {
	captiveCoreToml, err := NewCaptiveCoreToml(CaptiveCoreTomlParams{})
	assert.NoError(t, err)

	captiveCoreToml.AddExamplePubnetValidators()

	runner, err := newDiamcircleCoreRunner(CaptiveCoreConfig{
		HistoryArchiveURLs: []string{"http://localhost"},
		Log:                log.New(),
		Context:            context.Background(),
		Toml:               captiveCoreToml,
	}, diamcircleCoreRunnerModeOnline)
	assert.NoError(t, err)

	tempDir := runner.storagePath
	info, err := os.Stat(tempDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	assert.NoError(t, runner.close())

	// Directory no longer cleaned up on shutdown (perf. bump in v2.5.0)
	_, err = os.Stat(tempDir)
	assert.NoError(t, err)
}

func TestCloseBeforeStartOnlineWithError(t *testing.T) {
	captiveCoreToml, err := NewCaptiveCoreToml(CaptiveCoreTomlParams{})
	assert.NoError(t, err)

	captiveCoreToml.AddExamplePubnetValidators()

	runner, err := newDiamcircleCoreRunner(CaptiveCoreConfig{
		HistoryArchiveURLs: []string{"http://localhost"},
		Log:                log.New(),
		Context:            context.Background(),
		Toml:               captiveCoreToml,
	}, diamcircleCoreRunnerModeOnline)
	assert.NoError(t, err)

	runner.processExitError = errors.New("some error")

	tempDir := runner.storagePath
	info, err := os.Stat(tempDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	assert.NoError(t, runner.close())

	// Directory cleaned up on shutdown with error (potentially corrupted files)
	_, err = os.Stat(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

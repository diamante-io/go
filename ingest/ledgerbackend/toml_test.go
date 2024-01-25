package ledgerbackend

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newUint(v uint) *uint {
	return &v
}

func newString(s string) *string {
	return &s
}

func TestCaptiveCoreTomlValidation(t *testing.T) {
	for _, testCase := range []struct {
		name              string
		networkPassphrase string
		appendPath        string
		httpPort          *uint
		peerPort          *uint
		logPath           *string
		expectedError     string
	}{
		{
			name:              "mismatching NETWORK_PASSPHRASE",
			networkPassphrase: "bogus passphrase",
			appendPath:        filepath.Join("testdata", "appendix-with-fields.cfg"),
			httpPort:          newUint(6789),
			peerPort:          newUint(12345),
			logPath:           nil,
			expectedError: "invalid captive core toml: NETWORK_PASSPHRASE in captive core config file: " +
				"Public Global Diamcircle Network ; September 2015 does not match Aurora network-passphrase " +
				"flag: bogus passphrase",
		},
		{
			name:              "mismatching HTTP_PORT",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "appendix-with-fields.cfg"),
			httpPort:          newUint(1161),
			peerPort:          newUint(12345),
			logPath:           nil,
			expectedError: "invalid captive core toml: HTTP_PORT in captive core config file: 6789 " +
				"does not match Aurora captive-core-http-port flag: 1161",
		},
		{
			name:              "mismatching PEER_PORT",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "appendix-with-fields.cfg"),
			httpPort:          newUint(6789),
			peerPort:          newUint(2346),
			logPath:           nil,
			expectedError: "invalid captive core toml: PEER_PORT in captive core config file: 12345 " +
				"does not match Aurora captive-core-peer-port flag: 2346",
		},
		{
			name:              "mismatching LOG_FILE_PATH",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "appendix-with-fields.cfg"),
			httpPort:          newUint(6789),
			peerPort:          newUint(12345),
			logPath:           newString("/my/test/path"),
			expectedError: "invalid captive core toml: LOG_FILE_PATH in captive core config file:  " +
				"does not match Aurora captive-core-log-path flag: /my/test/path",
		},
		{
			name:              "duplicate HOME_DOMAIN entry",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "duplicate-home-domain.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found duplicate home domain in captive " +
				"core configuration: testnet.diamcircle.org",
		},
		{
			name:              "empty HOME_DOMAIN entry",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "empty-home-domain.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found invalid home domain entry which is " +
				"missing a HOME_DOMAIN value",
		},
		{
			name:              "HOME_DOMAIN with empty QUALITY",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "empty-home-domain-quality.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found invalid home domain entry which is " +
				"missing a QUALITY value: testnet.diamcircle.org",
		},
		{
			name:              "HOME_DOMAIN with invalid QUALITY",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "invalid-home-domain-quality.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found invalid home domain entry which has an " +
				"invalid QUALITY value: testnet.diamcircle.org",
		},
		{
			name:              "duplicate VALIDATOR entry",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "duplicate-validator.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found duplicate validator in captive core " +
				"configuration: sdf_testnet_1",
		},
		{
			name:              "VALIDATOR with invalid public key",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "validator-has-invalid-public-key.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found invalid validator entry which has an invalid " +
				"PUBLIC_KEY : sdf_testnet_2",
		},
		{
			name:              "VALIDATOR with invalid quality",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "validator-has-invalid-quality.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found invalid validator entry which has an invalid " +
				"QUALITY value: sdf_testnet_2",
		},
		{
			name:              "VALIDATOR without home domain",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "validator-missing-home-domain.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found invalid validator entry which is missing a " +
				"HOME_DOMAIN value: sdf_testnet_1",
		},
		{
			name:              "VALIDATOR without name",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "validator-missing-name.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found invalid validator entry which is missing " +
				"a NAME value: ",
		},
		{
			name:              "VALIDATOR without public key",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "validator-missing-public-key.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found invalid validator entry which is missing " +
				"a PUBLIC_KEY value: sdf_testnet_1",
		},
		{
			name:              "VALIDATOR without quality",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "validator-missing-quality.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError: "invalid captive core toml: found invalid validator entry which is missing " +
				"a QUALITY value: sdf_testnet_2",
		},
		{
			name:              "field not supported by captive core",
			networkPassphrase: "Public Global Diamcircle Network ; September 2015",
			appendPath:        filepath.Join("testdata", "invalid-captive-core-field.cfg"),
			httpPort:          nil,
			peerPort:          nil,
			logPath:           nil,
			expectedError:     "could not unmarshal captive core toml: these fields are not supported by captive core: [\"CATCHUP_RECENT\"]",
		},
		{
			name:          "unexpected BUCKET_DIR_PATH",
			appendPath:    filepath.Join("testdata", "appendix-with-bucket-dir-path.cfg"),
			expectedError: "could not unmarshal captive core toml: setting BUCKET_DIR_PATH is disallowed, it can cause clashes between instances",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			params := CaptiveCoreTomlParams{
				NetworkPassphrase:  testCase.networkPassphrase,
				HistoryArchiveURLs: []string{"http://localhost:1170"},
				HTTPPort:           testCase.httpPort,
				PeerPort:           testCase.peerPort,
				LogPath:            testCase.logPath,
				Strict:             true,
			}
			_, err := NewCaptiveCoreTomlFromFile(testCase.appendPath, params)
			assert.EqualError(t, err, testCase.expectedError)
		})
	}
}

func TestGenerateConfig(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		appendPath   string
		mode         diamcircleCoreRunnerMode
		expectedPath string
		httpPort     *uint
		peerPort     *uint
		logPath      *string
	}{
		{
			name:         "offline config with no appendix",
			mode:         diamcircleCoreRunnerModeOffline,
			appendPath:   "",
			expectedPath: filepath.Join("testdata", "expected-offline-core.cfg"),
			httpPort:     newUint(6789),
			peerPort:     newUint(12345),
			logPath:      nil,
		},
		{
			name:         "offline config with no peer port",
			mode:         diamcircleCoreRunnerModeOffline,
			appendPath:   "",
			expectedPath: filepath.Join("testdata", "expected-offline-with-no-peer-port.cfg"),
			httpPort:     newUint(6789),
			peerPort:     nil,
			logPath:      newString("/var/diamcircle-core/test.log"),
		},
		{
			name:         "online config with appendix",
			mode:         diamcircleCoreRunnerModeOnline,
			appendPath:   filepath.Join("testdata", "sample-appendix.cfg"),
			expectedPath: filepath.Join("testdata", "expected-online-core.cfg"),
			httpPort:     newUint(6789),
			peerPort:     newUint(12345),
			logPath:      nil,
		},
		{
			name:         "online config with unsupported field in appendix",
			mode:         diamcircleCoreRunnerModeOnline,
			appendPath:   filepath.Join("testdata", "invalid-captive-core-field.cfg"),
			expectedPath: filepath.Join("testdata", "expected-online-core.cfg"),
			httpPort:     newUint(6789),
			peerPort:     newUint(12345),
			logPath:      nil,
		},
		{
			name:         "online config with no peer port",
			mode:         diamcircleCoreRunnerModeOnline,
			appendPath:   filepath.Join("testdata", "sample-appendix.cfg"),
			expectedPath: filepath.Join("testdata", "expected-online-with-no-peer-port.cfg"),
			httpPort:     newUint(6789),
			peerPort:     nil,
			logPath:      newString("/var/diamcircle-core/test.log"),
		},
		{
			name:         "online config with no http port",
			mode:         diamcircleCoreRunnerModeOnline,
			appendPath:   filepath.Join("testdata", "sample-appendix.cfg"),
			expectedPath: filepath.Join("testdata", "expected-online-with-no-http-port.cfg"),
			httpPort:     nil,
			peerPort:     newUint(12345),
			logPath:      nil,
		},
		{
			name:         "offline config with appendix",
			mode:         diamcircleCoreRunnerModeOffline,
			appendPath:   filepath.Join("testdata", "sample-appendix.cfg"),
			expectedPath: filepath.Join("testdata", "expected-offline-with-appendix-core.cfg"),
			httpPort:     newUint(6789),
			peerPort:     newUint(12345),
			logPath:      nil,
		},
		{
			name:         "offline config with extra fields in appendix",
			mode:         diamcircleCoreRunnerModeOffline,
			appendPath:   filepath.Join("testdata", "appendix-with-fields.cfg"),
			expectedPath: filepath.Join("testdata", "expected-offline-with-extra-fields.cfg"),
			httpPort:     newUint(6789),
			peerPort:     newUint(12345),
			logPath:      nil,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var err error
			var captiveCoreToml *CaptiveCoreToml
			params := CaptiveCoreTomlParams{
				NetworkPassphrase:  "Public Global Diamcircle Network ; September 2015",
				HistoryArchiveURLs: []string{"http://localhost:1170"},
				HTTPPort:           testCase.httpPort,
				PeerPort:           testCase.peerPort,
				LogPath:            testCase.logPath,
				Strict:             false,
			}
			if testCase.appendPath != "" {
				captiveCoreToml, err = NewCaptiveCoreTomlFromFile(testCase.appendPath, params)
			} else {
				captiveCoreToml, err = NewCaptiveCoreToml(params)
			}
			assert.NoError(t, err)

			configBytes, err := generateConfig(captiveCoreToml, testCase.mode)
			assert.NoError(t, err)

			expectedByte, err := ioutil.ReadFile(testCase.expectedPath)
			assert.NoError(t, err)

			assert.Equal(t, string(configBytes), string(expectedByte))
		})
	}
}

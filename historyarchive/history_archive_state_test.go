// Copyright 2016 diamcircle Development Foundation and contributors. Licensed
// under the Apache License, Version 2.0. See the COPYING file at the root
// of this distribution or at http://www.apache.org/licenses/LICENSE-2.0

package historyarchive

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalState(t *testing.T) {
	var jsonBlob = []byte(`{
		"version": 1,
		"server": "v0.4.0-34-g2f015f6",
		"currentLedger": 2113919,
		"currentBuckets": [
			{
				"curr": "0000000000000000000000000000000000000000000000000000000000000000",
				"next": {
					"state": 0
				},
				"snap": "0000000000000000000000000000000000000000000000000000000000000000"
			},
			{
				"curr": "0000000000000000000000000000000000000000000000000000000000000000",
				"next": {
					"state": 1,
					"output": "0000000000000000000000000000000000000000000000000000000000000000"
				},
				"snap": "0000000000000000000000000000000000000000000000000000000000000000"
			}
		 ]
	}`)

	var state HistoryArchiveState
	err := json.Unmarshal(jsonBlob, &state)
	if err != nil {
		t.Error(err)
	} else if state.CurrentLedger != 2113919 {
		t.Error(state)
	}
}

func TestHashValidation(t *testing.T) {
	// This is real bucket hash list for pubnet's ledger: 24088895
	// https://aurora.diamcircle.org/ledgers/24088895
	// http://history.diamcircle.org/prd/core-live/core_live_001/history/01/6f/91/history-016f913f.json
	var jsonBlob = []byte(`{
    "version": 1,
    "server": "v11.1.0",
    "currentLedger": 24088895,
    "currentBuckets": [
        {
            "curr": "2a4416e7f3e301c2fc1078dce0e1dd109b8ae6d3958942b91b447f24014a7b5c",
            "next": {
                "state": 0
            },
            "snap": "7ff95a98838dfd39a36858f15c8d503641560f02a52aa15335559e1183ce2ca1"
        },
        {
            "curr": "2c7e74c4c5555e41b39a5fc04e77e77852c35e7769e32b486e07a072b9b3177c",
            "next": {
                "state": 1,
                "output": "7ff95a98838dfd39a36858f15c8d503641560f02a52aa15335559e1183ce2ca1"
            },
            "snap": "5f0bc7d0bd9e8ed6530fc270339b7dd2fbcedf0d80235f5ef64daa90b84259f4"
        },
        {
            "curr": "068f2a1ece2817c98c0d21d5ac20817637c331df6793d0ff3e874e29da5d65b1",
            "next": {
                "state": 1,
                "output": "e93d50365d74d8a8dc2ff7631dfb506b7e6b2245f7f46556d407e82f197a6c59"
            },
            "snap": "875cbdf9ab03c488c075a36ee3ee1e02aef9d5fe9d253a2b1f99b92fe64598b8"
        },
        {
            "curr": "f413ff9d27e2cad12754ff84ca905f8c309ca7b68a6fbe8e9b01ecd18f5d3759",
            "next": {
                "state": 1,
                "output": "ffbb6cd3a4170dbf547ab0783fea90c1a57a28e562db7bcd3a079374f1e63464"
            },
            "snap": "5d198cdc5a2139d92fe74f6541a098c27aba61e8aee543a6a551814aae9adb5a"
        },
        {
            "curr": "1c6f9ec76b06aac2aac77e9a0da2224d120dc25c1cf10211ce33475db4d66f13",
            "next": {
                "state": 1,
                "output": "6473d4a3ff5b6448fc6dfd279ef33bf0b1524d8361b758dbde49fc84691cadbe"
            },
            "snap": "6dd30650a7c8cadad545561d732828cf55cefdf5f70c615fbdc33e01c647907b"
        },
        {
            "curr": "b3b3c9b54db9e08f3994a17a40e5be7583c546488e88523ebf8b444ee53f4aec",
            "next": {
                "state": 1,
                "output": "ed452df8b803190b7a1cf07894c27c03415029272e9c4a3171e7f3ad6e67c90a"
            },
            "snap": "7d84d34019975b37758278e858e86265323ddbb7b46f6d679433c93bbae693ee"
        },
        {
            "curr": "a6c20a247ed2afc2cea7f4dc5856efa61a51b4e4b0318877eebdf8ad47be83b7",
            "next": {
                "state": 1,
                "output": "ce9a7c779d0873ff364a9abd20007bbf7e41646ac4662eb87f89a5c39b69f70d"
            },
            "snap": "285ac930ee2bd358d3202666c545fd3b94ee973d1a0cd2569de357042ec12b3d"
        },
        {
            "curr": "2e779b37b97052a1141a65a92c4ca14a7bd28f7c2d646749b1d584f45d50fa7b",
            "next": {
                "state": 1,
                "output": "e4dba3994ad576489880eee38db2d8c0f8889585e932b7192dd7af168d79b43f"
            },
            "snap": "37094a837769dbae5783dca9831be463b895f1b07d1cd24e271966e10503fdfc"
        },
        {
            "curr": "48f435285dd96511d0822f7ae1a20e28c6c28019e385313713655fc76fe3bc03",
            "next": {
                "state": 1,
                "output": "11f8c2f8e1cb0d47576f74d9e2fa838f5f3a37180907a24a85d0ad8b647862e4"
            },
            "snap": "96e0d8bf7d7eb775299edf285b6324499a1a05122d95eed9289c6477cf6a01cb"
        },
        {
            "curr": "4100ad3b1085bd14d1c808ece3b38db97171532d0d11ed5edd57aff0e416e06a",
            "next": {
                "state": 1,
                "output": "5f351041761b45f3e725f98bb8b6713873e30ab6c8aee56ba0823d357c7ebd0d"
            },
            "snap": "23669fa3d310ca8ac8dbe9dcce7e4e4361b1c3334da1dda2fb6447a30c67422f"
        },
        {
            "curr": "14cc632ab181396418fc761503105047e3b63d0455d0a4e9480578129ea8e9dc",
            "next": {
                "state": 1,
                "output": "a4811c9ba9505e421f0015e5fcfd9f5d204ae85b584766759e844ef85db10d47"
            },
            "snap": "0000000000000000000000000000000000000000000000000000000000000000"
        }
    ]
}`)

	var state HistoryArchiveState
	err := json.Unmarshal(jsonBlob, &state)
	require.NoError(t, err)

	expectedHash, err := hex.DecodeString("fc5fe47af3f5a9b18b278f2a7edbbc641e1934bf68131d9aa5ab7aebb4aa8aa3")
	require.NoError(t, err)

	hash, err := state.BucketListHash()
	require.NoError(t, err)
	assert.Equal(t, expectedHash, hash[:])
}

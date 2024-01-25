package history

import (
	"sort"
	"testing"

	"go/services/aurora/internal/test"

	"github.com/stretchr/testify/assert"
)

func TestIsAuthRequired(t *testing.T) {
	tt := assert.New(t)

	account := AccountEntry{Flags: 1}
	tt.True(account.IsAuthRequired())

	account = AccountEntry{Flags: 0}
	tt.False(account.IsAuthRequired())
}

func TestIsAuthRevocable(t *testing.T) {
	tt := assert.New(t)

	account := AccountEntry{Flags: 2}
	tt.True(account.IsAuthRevocable())

	account = AccountEntry{Flags: 1}
	tt.False(account.IsAuthRevocable())
}
func TestIsAuthImmutable(t *testing.T) {
	tt := assert.New(t)

	account := AccountEntry{Flags: 4}
	tt.True(account.IsAuthImmutable())

	account = AccountEntry{Flags: 0}
	tt.False(account.IsAuthImmutable())
}

func assertAccountsContainAddresses(tt *test.T, accounts map[string]int64, addresses []string) {
	tt.Assert.Len(accounts, len(addresses))
	set := map[int64]bool{}
	for _, address := range addresses {
		accountID, ok := accounts[address]
		tt.Assert.True(ok)
		tt.Assert.False(set[accountID])
		set[accountID] = true
	}
}

func TestCreateAccountsSortedOrder(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	addresses := []string{
		"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		"GCYVFGI3SEQJGBNQQG7YCMFWEYOHK3XPVOVPA6C566PXWN4SN7LILZSM",
		"GBYSBDAJZMHL5AMD7QXQ3JEP3Q4GLKADWIJURAAHQALNAWD6Z5XF2RAC",
		"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
	}
	accounts, err := q.CreateAccounts(tt.Ctx, addresses, 1)
	tt.Assert.NoError(err)

	idToAddress := map[int64]string{}
	sortedIDs := []int64{}
	for address, id := range accounts {
		idToAddress[id] = address
		sortedIDs = append(sortedIDs, id)
	}

	sort.Slice(sortedIDs, func(i, j int) bool {
		return sortedIDs[i] < sortedIDs[j]
	})
	sort.Strings(addresses)

	values := []string{}
	for _, id := range sortedIDs {
		values = append(values, idToAddress[id])
	}

	tt.Assert.Equal(addresses, values)
}

func TestCreateAccounts(t *testing.T) {
	tt := test.Start(t)
	defer tt.Finish()
	test.ResetauroraDB(t, tt.auroraDB)
	q := &Q{tt.auroraSession()}

	addresses := []string{
		"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
	}
	accounts, err := q.CreateAccounts(tt.Ctx, addresses, 1)
	tt.Assert.NoError(err)
	tt.Assert.Len(accounts, 2)
	assertAccountsContainAddresses(tt, accounts, addresses)

	dupAccounts, err := q.CreateAccounts(tt.Ctx, []string{
		"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
	}, 2)
	tt.Assert.NoError(err)
	tt.Assert.Equal(accounts, dupAccounts)

	addresses = []string{
		"GAOQJGUAB7NI7K7I62ORBXMN3J4SSWQUQ7FOEPSDJ322W2HMCNWPHXFB",
		"GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		"GCYVFGI3SEQJGBNQQG7YCMFWEYOHK3XPVOVPA6C566PXWN4SN7LILZSM",
		"GBYSBDAJZMHL5AMD7QXQ3JEP3Q4GLKADWIJURAAHQALNAWD6Z5XF2RAC",
	}
	accounts, err = q.CreateAccounts(tt.Ctx, addresses, 1)
	tt.Assert.NoError(err)
	assertAccountsContainAddresses(tt, accounts, addresses)
	for address, accountID := range dupAccounts {
		id, ok := accounts[address]
		tt.Assert.True(ok)
		tt.Assert.Equal(id, accountID)
	}
}

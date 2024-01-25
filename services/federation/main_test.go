package main

import (
	"testing"

	"github.com/diamcircle/go/support/db/dbtest"
	"github.com/diamcircle/go/support/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitDriver_dialect(t *testing.T) {
	c := Config{}

	testCases := []struct {
		dbType  string
		dbDSN   string
		wantErr error
	}{
		{dbType: "", wantErr: errors.New("Invalid db type: ")},
		{dbType: "postgres", dbDSN: dbtest.Postgres(t).DSN, wantErr: nil},
		{dbType: "mysql", wantErr: errors.New("Invalid db type: mysql, mysql support is discontinued")},
		{dbType: "bogus", wantErr: errors.New("Invalid db type: bogus")},
	}

	for _, tc := range testCases {
		t.Run(tc.dbType, func(t *testing.T) {
			c.Database.Type = tc.dbType
			c.Database.DSN = tc.dbDSN
			_, err := initDriver(c)
			if tc.wantErr == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			}
		})
	}
}

package federation

import (
	"context"

	"go/support/db"
	"go/support/errors"
)

// LookupRecord implements `Driver` by performing `drv.LookupRecordQuery`
// against `drv.DB` using the provided parameters
func (drv *SQLDriver) LookupRecord(ctx context.Context, name, domain string) (*Record, error) {
	drv.initDB()
	var result Record

	err := drv.db.GetRaw(ctx, &result, drv.LookupRecordQuery, name, domain)

	if drv.db.NoRows(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "db get")
	}

	return &result, nil
}

var _ Driver = &SQLDriver{}

func (drv *SQLDriver) initDB() {
	drv.init.Do(func() {
		if drv.Dialect == "" {
			panic("no dialect specified")
		}

		drv.db = db.Wrap(drv.DB, drv.Dialect)
	})
}

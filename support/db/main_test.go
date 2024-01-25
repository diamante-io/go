package db

import (
	"testing"

	"go/support/db/dbtest"

	"github.com/stretchr/testify/assert"
)

type person struct {
	Name        string `db:"name"`
	HungerLevel string `db:"hunger_level"`

	SomethingIgnored int `db:"-"`
}

func TestGetTable(t *testing.T) {
	db := dbtest.Postgres(t).Load(testSchema)
	defer db.Close()
	sess := &Session{DB: db.Open()}
	defer sess.DB.Close()

	tbl := sess.GetTable("users")
	if assert.NotNil(t, tbl) {
		assert.Equal(t, "users", tbl.Name)
		assert.Equal(t, sess, tbl.Session)
	}

}

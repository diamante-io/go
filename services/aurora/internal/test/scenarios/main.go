package scenarios

import (
	"log"

	"github.com/jmoiron/sqlx"
)

//go:generate go run github.com/kevinburke/go-bindata/go-bindata@v3.18.0+incompatible -nometadata -ignore (go|rb)$ -pkg scenarios .

// Load executes the sql script at `path` on postgres database at `url`
func Load(url string, path string) {
	sql, err := Asset(path)
	if err != nil {
		log.Panic(err)
	}

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		log.Fatalf("could not exec open postgres connection: %v\n", err)
	}
	defer db.Close()

	// clear out existing schema before applying scenario
	// otherwise, applying the scenario will result in the following error:
	// pq: cannot drop schema public because other objects depend on it
	_, err = db.Exec("DROP SCHEMA IF EXISTS public cascade")
	if err != nil {
		log.Fatalf("could not drop public schema: %v\n", err)
	}

	_, err = db.Exec(string(sql))
	if err != nil {
		log.Fatalf("could not exec scenario %v: %v\n", path, err)
	}
}

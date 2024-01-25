package main

import (
	"fmt"
	"os"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"github.com/diamcircle/go/handlers/federation"
	"github.com/diamcircle/go/support/app"
	"github.com/diamcircle/go/support/config"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/http"
	"github.com/diamcircle/go/support/log"
)

// Config represents the configuration of a federation server
type Config struct {
	Port     int `valid:"required"`
	Database struct {
		Type string `valid:"matches(^sqlite3|postgres$)"`
		DSN  string `valid:"required"`
	} `valid:"required"`
	Queries struct {
		Federation        string `valid:"required"`
		ReverseFederation string `toml:"reverse-federation" valid:"optional"`
	} `valid:"required"`
	TLS *config.TLS `valid:"optional"`
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "federation",
		Short: "diamcircle federation server",
		Long: `
The diamcircle federation server let's you easily integrate the diamcircle federation
protocol with your organization.  This is achieved by connecting the
application to your customer database and providing the appropriate queries in
the config file.
    `,
		Run: run,
	}

	rootCmd.PersistentFlags().String("conf", "./federation.cfg", "config file path")
	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	var (
		cfg     Config
		cfgPath = cmd.PersistentFlags().Lookup("conf").Value.String()
	)
	log.SetLevel(log.InfoLevel)
	err := config.Read(cfgPath, &cfg)

	if err != nil {
		switch cause := errors.Cause(err).(type) {
		case *config.InvalidConfigError:
			log.Error("config file: ", cause)
		default:
			log.Error(err)
		}
		os.Exit(1)
	}

	driver, err := initDriver(cfg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	mux := initMux(driver)
	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)

	http.Run(http.Config{
		ListenAddr: addr,
		Handler:    mux,
		TLS:        cfg.TLS,
		OnStarting: func() {
			log.Infof("starting federation server - %s", app.Version())
			log.Infof("listening on %s", addr)
		},
	})
}

func initDriver(cfg Config) (federation.Driver, error) {
	var dialect string

	switch cfg.Database.Type {
	case "mysql":
		return nil, errors.Errorf("Invalid db type: %s, mysql support is discontinued", cfg.Database.Type)
	case "postgres":
		dialect = "postgres"
	case "sqlite3":
		dialect = "sqlite3"
	default:
		return nil, errors.Errorf("Invalid db type: %s", cfg.Database.Type)
	}

	repo, err := db.Open(dialect, cfg.Database.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "db open failed")
	}

	sqld := federation.SQLDriver{
		DB:                repo.DB.DB, // unwrap the repo to the bare *sql.DB instance,
		Dialect:           dialect,
		LookupRecordQuery: cfg.Queries.Federation,
	}

	if cfg.Queries.ReverseFederation == "" {
		return &sqld, nil
	}

	rsqld := federation.ReverseSQLDriver{
		SQLDriver: federation.SQLDriver{
			DB:                repo.DB.DB,
			Dialect:           dialect,
			LookupRecordQuery: cfg.Queries.Federation,
		},
		LookupReverseRecordQuery: cfg.Queries.ReverseFederation,
	}

	return &rsqld, nil
}

func initMux(driver federation.Driver) *chi.Mux {
	mux := http.NewAPIMux(log.DefaultLogger)

	fed := &federation.Handler{
		Driver: driver,
	}

	mux.Get("/federation", fed.ServeHTTP)
	mux.Get("/federation/", fed.ServeHTTP)

	return mux
}

package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/diamcircle/go/services/keystore"
	"github.com/diamcircle/go/support/log"

	_ "github.com/lib/pq"
)

var keystoreMigrations = &migrate.FileMigrationSource{
	Dir: "migrations",
}

const dbDriverName = "postgres"

func main() {
	ctx := context.Background()
	tlsCert := flag.String("tls-cert", "", "TLS certificate file path")
	tlsKey := flag.String("tls-key", "", "TLS private key file path")
	logFilePath := flag.String("log-file", "", "Log file file path")
	logLevel := flag.String("log-level", "info", "Log level used by logrus (debug, info, warn, error)")
	auth := flag.Bool("auth", true, "Enable authentication")
	apiType := flag.String("api-type", "REST", "Auth Forwarding API Type")

	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "too few arguments")
		os.Exit(1)
	}

	if (*tlsCert == "" && *tlsKey != "") || (*tlsCert != "" && *tlsKey == "") {
		fmt.Fprintln(os.Stderr, "TLS cert and TLS key have to be presented together")
		os.Exit(1)
	}

	if *logFilePath != "" {
		logFile, err := os.OpenFile(*logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file to log: %v\n", err)
			os.Exit(1)
		}

		log.DefaultLogger.SetOutput(logFile)

		ll, err := logrus.ParseLevel(*logLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse log-level: %v\n", err)
			os.Exit(1)
		}
		log.DefaultLogger.SetLevel(ll)
	}

	cfg := getConfig()
	if cfg.ListenerPort < 0 {
		fmt.Fprintf(os.Stderr, "Port number %d cannot be negative\n", cfg.ListenerPort)
		os.Exit(1)
	}

	db, err := sql.Open(dbDriverName, cfg.DBURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening database: %v\n", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(cfg.MaxOpenDBConns)
	db.SetMaxIdleConns(cfg.MaxIdleDBConns)

	err = db.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error accessing database: %v\n", err)
		os.Exit(1)
	}

	cmd := flag.Arg(0)
	switch cmd {
	case "serve":
		if *auth {
			if cfg.AUTHURL == "" {
				fmt.Fprintln(os.Stderr, "Auth is enabled but auth forwarding URL is not set")
				os.Exit(1)
			}
			if _, err := url.Parse(cfg.AUTHURL); err != nil {
				fmt.Fprintln(os.Stderr, "Invalid auth forwarding URL")
				os.Exit(1)
			}
		}

		aType := strings.ToUpper(*apiType)
		if aType != keystore.REST && aType != keystore.GraphQL {
			fmt.Fprintln(os.Stderr, `Auth forwarding endpoint type can only be either "REST" or "GRAPHQL"`)
			os.Exit(1)
		}

		addr := ":" + strconv.Itoa(cfg.ListenerPort)
		var authenticator *keystore.Authenticator
		if *auth {
			authenticator = &keystore.Authenticator{
				URL:     cfg.AUTHURL,
				APIType: aType,
			}
		}

		server := &http.Server{
			Addr:        addr,
			Handler:     keystore.ServeMux(keystore.NewService(ctx, db, authenticator)),
			ReadTimeout: 5 * time.Second,
		}

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error listening: %v\n", err)
			os.Exit(1)
		}

		listener = tcpKeepAliveListener{listener.(*net.TCPListener)}
		if *tlsCert != "" {
			cer, err := tls.LoadX509KeyPair(*tlsCert, *tlsKey)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error parsing TLS keypair: %v\n", err)
				os.Exit(1)
			}

			listener = tls.NewListener(listener, &tls.Config{Certificates: []tls.Certificate{cer}})
		}

		go func() {
			err := server.Serve(listener)
			if err != nil {
				panic(err)
			}
		}()
		fmt.Fprintln(os.Stdout, "Server listening at https://localhost"+addr)

		// block forever without using any resources so this process won't quit while
		// the goroutine containing ListenAndServe is still working
		select {}

	case "migrate":
		migrateCmd := flag.Arg(1)
		switch migrateCmd {
		case "up":
			n, err := migrate.Exec(db, dbDriverName, keystoreMigrations, migrate.Up)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error applying up migrations: %v\n", err)
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "Applied %d up migrations!\n", n)

		case "down":
			n, err := migrate.ExecMax(db, dbDriverName, keystoreMigrations, migrate.Down, 1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error applying down migrations: %v\n", err)
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "Applied %d down migration!\n", n)

		case "redo":
			migrations, _, err := migrate.PlanMigration(db, dbDriverName, keystoreMigrations, migrate.Down, 1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting migration data: %v\n", err)
				os.Exit(1)
			}

			if len(migrations) == 0 {
				fmt.Fprintln(os.Stdout, "Nothing to do!")
				os.Exit(0)
			}

			_, err = migrate.ExecMax(db, dbDriverName, keystoreMigrations, migrate.Down, 1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error applying the last down migration: %v\n", err)
				os.Exit(1)
			}

			_, err = migrate.ExecMax(db, dbDriverName, keystoreMigrations, migrate.Up, 1)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error applying the last up migration: %v\n", err)
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "Reapplied migration %s.\n", migrations[0].Id)

		case "status":
			unappliedMigrations := getUnappliedMigrations(db)
			if len(unappliedMigrations) > 0 {
				fmt.Fprintf(os.Stdout, "There are %d unapplied migrations:\n", len(unappliedMigrations))
				for _, id := range unappliedMigrations {
					fmt.Fprintln(os.Stdout, id)
				}
			} else {
				fmt.Fprintln(os.Stdout, "All migrations have been applied!")
			}

		default:
			fmt.Fprintf(os.Stderr, "unrecognized migration command: %q\n", migrateCmd)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "unrecognized command: %q\n", cmd)
		os.Exit(1)
	}
}

// https://github.com/golang/go/blob/c5cf6624076a644906aa7ec5c91c4e01ccd375d3/src/net/http/server.go#L3272-L3288
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func getUnappliedMigrations(db *sql.DB) []string {
	migrations, err := keystoreMigrations.FindMigrations()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting keystore migrations: %v\n", err)
		os.Exit(1)
	}

	records, err := migrate.GetMigrationRecords(db, dbDriverName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting keystore migrations records: %v\n", err)
		os.Exit(1)
	}

	unappliedMigrations := make(map[string]struct{})
	for _, m := range migrations {
		unappliedMigrations[m.Id] = struct{}{}
	}

	for _, r := range records {
		if _, ok := unappliedMigrations[r.Id]; !ok {
			fmt.Fprintf(os.Stdout, "Could not find migration file: %v\n", r.Id)
			continue
		}

		delete(unappliedMigrations, r.Id)
	}

	result := make([]string, 0, len(unappliedMigrations))
	for id := range unappliedMigrations {
		result = append(result, id)
	}

	sort.Strings(result)

	return result
}

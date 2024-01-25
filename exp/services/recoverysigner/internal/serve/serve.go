package serve

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go/exp/services/recoverysigner/internal/account"
	"go/exp/services/recoverysigner/internal/db"
	"go/exp/services/recoverysigner/internal/serve/auth"
	"go/keypair"
	"go/support/errors"
	supporthttp "go/support/http"
	supportlog "go/support/log"
	"go/support/render/health"

	firebaseauth "firebase.google.com/go/auth"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/square/go-jose.v2"
)

type Options struct {
	Logger               *supportlog.Entry
	DatabaseURL          string
	DatabaseMaxOpenConns int
	Port                 int
	NetworkPassphrase    string
	SigningKeys          string
	SEP10JWKS            string
	SEP10JWTIssuer       string
	FirebaseProjectID    string

	AdminPort        int
	MetricsNamespace string

	AllowedSourceAccounts string
}

func Serve(opts Options) {
	deps, err := getHandlerDeps(opts)
	if err != nil {
		opts.Logger.Fatalf("Error: %v", err)
		return
	}

	if opts.AdminPort != 0 {
		adminDeps := adminDeps{
			Logger:          opts.Logger,
			MetricsGatherer: deps.MetricsRegistry,
		}
		go serveAdmin(opts, adminDeps)
	}

	handler := handler(deps)

	addr := fmt.Sprintf(":%d", opts.Port)
	supporthttp.Run(supporthttp.Config{
		ListenAddr: addr,
		Handler:    handler,
		OnStarting: func() {
			deps.Logger.Infof("Starting SEP-30 Recovery Signer server on %s", addr)
		},
	})
}

type handlerDeps struct {
	Logger                *supportlog.Entry
	NetworkPassphrase     string
	SigningKeys           []*keypair.Full
	SigningAddresses      []*keypair.FromAddress
	AccountStore          account.Store
	SEP10JWKS             jose.JSONWebKeySet
	SEP10JWTIssuer        string
	FirebaseAuthClient    *firebaseauth.Client
	MetricsRegistry       *prometheus.Registry
	AllowedSourceAccounts []*keypair.FromAddress
}

func getHandlerDeps(opts Options) (handlerDeps, error) {
	// TODO: Replace this signing key with randomly generating a unique signing
	// key for each account so that it is not possible to identify which
	// accounts are recoverable via a recovery signer.
	signingKeys := []*keypair.Full{}
	signingAddresses := []*keypair.FromAddress{}
	for i, signingKeyStr := range strings.Split(opts.SigningKeys, ",") {
		signingKey, err := keypair.ParseFull(signingKeyStr)
		if err != nil {
			return handlerDeps{}, errors.Wrap(err, "parsing signing key seed")
		}
		signingKeys = append(signingKeys, signingKey)
		signingAddresses = append(signingAddresses, signingKey.FromAddress())
		opts.Logger.Info("Signing key ", i, ": ", signingKey.Address())
	}

	sep10JWKS := jose.JSONWebKeySet{}
	err := json.Unmarshal([]byte(opts.SEP10JWKS), &sep10JWKS)
	if err != nil {
		return handlerDeps{}, errors.Wrap(err, "parsing SEP-10 JSON Web Key (JWK) Set")
	}
	if len(sep10JWKS.Keys) == 0 {
		return handlerDeps{}, errors.New("no keys included in SEP-10 JSON Web Key (JWK) Set")
	}
	opts.Logger.Infof("SEP10 JWKS contains %d keys", len(sep10JWKS.Keys))

	db, err := db.Open(opts.DatabaseURL)
	if err != nil {
		return handlerDeps{}, errors.Wrap(err, "error parsing database url")
	}
	db.SetMaxOpenConns(opts.DatabaseMaxOpenConns)

	err = db.Ping()
	if err != nil {
		opts.Logger.Warn("Error pinging to Database: ", err)
	}
	accountStore := &account.DBStore{DB: db}

	firebaseAuthClient, err := auth.NewFirebaseAuthClient(opts.FirebaseProjectID)
	if err != nil {
		return handlerDeps{}, errors.Wrap(err, "error setting up firebase auth client")
	}

	metricsRegistry := prometheus.NewRegistry()

	err = metricsRegistry.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	if err != nil {
		opts.Logger.Warn("Error registering metric for process: ", err)
	}
	err = metricsRegistry.Register(prometheus.NewGoCollector())
	if err != nil {
		opts.Logger.Warn("Error registering metric for Go: ", err)
	}

	metricsRegistryNamespaced := prometheus.Registerer(metricsRegistry)
	if opts.MetricsNamespace != "" {
		metricsRegistryNamespaced = prometheus.WrapRegistererWithPrefix(opts.MetricsNamespace+"_", metricsRegistry)
	}

	err = metricsRegistryNamespaced.Register(metricAccountsCount{
		Logger:       opts.Logger,
		AccountStore: accountStore,
	}.NewCollector())
	if err != nil {
		opts.Logger.Warn("Error registering metric for accounts count: ", err)
	}

	allowedSourceAccounts := []*keypair.FromAddress{}
	for _, addressStr := range strings.Split(opts.AllowedSourceAccounts, ",") {
		accountAddress, err := keypair.ParseAddress(addressStr)
		if err != nil {
			return handlerDeps{}, errors.Wrap(err, "parsing allowed source accounts")
		}
		allowedSourceAccounts = append(allowedSourceAccounts, accountAddress)
	}

	deps := handlerDeps{
		Logger:                opts.Logger,
		NetworkPassphrase:     opts.NetworkPassphrase,
		SigningKeys:           signingKeys,
		SigningAddresses:      signingAddresses,
		AccountStore:          accountStore,
		SEP10JWKS:             sep10JWKS,
		SEP10JWTIssuer:        opts.SEP10JWTIssuer,
		FirebaseAuthClient:    firebaseAuthClient,
		MetricsRegistry:       metricsRegistry,
		AllowedSourceAccounts: allowedSourceAccounts,
	}

	return deps, nil
}

func handler(deps handlerDeps) http.Handler {
	mux := supporthttp.NewAPIMux(deps.Logger)

	mux.NotFound(errorHandler{Error: notFound}.ServeHTTP)
	mux.MethodNotAllowed(errorHandler{Error: methodNotAllowed}.ServeHTTP)

	mux.Get("/health", health.PassHandler{}.ServeHTTP)
	mux.Route("/accounts", func(mux chi.Router) {
		mux.Use(auth.SEP10Middleware(deps.SEP10JWTIssuer, deps.SEP10JWKS))
		mux.Use(auth.FirebaseMiddleware(auth.FirebaseTokenVerifierLive{AuthClient: deps.FirebaseAuthClient}))
		mux.Get("/", accountListHandler{
			Logger:           deps.Logger,
			SigningAddresses: deps.SigningAddresses,
			AccountStore:     deps.AccountStore,
		}.ServeHTTP)
		mux.Route("/{address}", func(mux chi.Router) {
			mux.Post("/", accountPostHandler{
				Logger:           deps.Logger,
				SigningAddresses: deps.SigningAddresses,
				AccountStore:     deps.AccountStore,
			}.ServeHTTP)
			mux.Put("/", accountPutHandler{
				Logger:           deps.Logger,
				SigningAddresses: deps.SigningAddresses,
				AccountStore:     deps.AccountStore,
			}.ServeHTTP)
			mux.Get("/", accountGetHandler{
				Logger:           deps.Logger,
				SigningAddresses: deps.SigningAddresses,
				AccountStore:     deps.AccountStore,
			}.ServeHTTP)
			mux.Delete("/", accountDeleteHandler{
				Logger:           deps.Logger,
				SigningAddresses: deps.SigningAddresses,
				AccountStore:     deps.AccountStore,
			}.ServeHTTP)
			signHandler := accountSignHandler{
				Logger:                deps.Logger,
				SigningKeys:           deps.SigningKeys,
				NetworkPassphrase:     deps.NetworkPassphrase,
				AccountStore:          deps.AccountStore,
				AllowedSourceAccounts: deps.AllowedSourceAccounts,
			}
			mux.Post("/sign", signHandler.ServeHTTP)
			mux.Post("/sign/{signing-address}", signHandler.ServeHTTP)
		})
	})

	return mux
}

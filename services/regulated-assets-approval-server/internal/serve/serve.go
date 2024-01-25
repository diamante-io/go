package serve

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/diamcircle/go/amount"
	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/services/regulated-assets-approval-server/internal/db"
	"github.com/diamcircle/go/services/regulated-assets-approval-server/internal/serve/kycstatus"
	"github.com/diamcircle/go/support/errors"
	supporthttp "github.com/diamcircle/go/support/http"
	"github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/support/render/health"
)

type Options struct {
	AssetCode                         string
	BaseURL                           string
	DatabaseURL                       string
	FriendbotPaymentAmount            int
	AuroraURL                        string
	IssuerAccountSecret               string
	KYCRequiredPaymentAmountThreshold string
	NetworkPassphrase                 string
	Port                              int
}

func Serve(opts Options) {
	listenAddr := fmt.Sprintf(":%d", opts.Port)
	serverConfig := supporthttp.Config{
		ListenAddr:          listenAddr,
		Handler:             handleHTTP(opts),
		TCPKeepAlive:        time.Minute * 3,
		ShutdownGracePeriod: time.Second * 50,
		ReadTimeout:         time.Second * 5,
		WriteTimeout:        time.Second * 35,
		IdleTimeout:         time.Minute * 2,
		OnStarting: func() {
			log.Info("Starting SEP-8 Approval Server")
			log.Infof("Listening on %s", listenAddr)
		},
		OnStopping: func() {
			log.Info("Stopping SEP-8 Approval Server")
		},
	}
	supporthttp.Run(serverConfig)
}

func handleHTTP(opts Options) http.Handler {
	issuerKP, err := keypair.ParseFull(opts.IssuerAccountSecret)
	if err != nil {
		log.Fatal(errors.Wrap(err, "parsing secret"))
	}
	parsedKYCRequiredPaymentThreshold, err := amount.ParseInt64(opts.KYCRequiredPaymentAmountThreshold)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "%s cannot be parsed as a Diamcircle amount", opts.KYCRequiredPaymentAmountThreshold))
	}
	db, err := db.Open(opts.DatabaseURL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error parsing database url"))
	}
	db.SetMaxOpenConns(20)
	err = db.Ping()
	if err != nil {
		log.Warn("Error pinging to Database: ", err)
	}
	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(supporthttp.LoggingMiddleware)
	mux.Use(corsHandler)

	mux.Get("/health", health.PassHandler{}.ServeHTTP)
	mux.Get("/.well-known/diamcircle.toml", diamcircleTOMLHandler{
		assetCode:         opts.AssetCode,
		issuerAddress:     issuerKP.Address(),
		networkPassphrase: opts.NetworkPassphrase,
		approvalServer:    buildURLString(opts.BaseURL, "tx-approve"),
		kycThreshold:      parsedKYCRequiredPaymentThreshold,
	}.ServeHTTP)
	mux.Get("/friendbot", friendbotHandler{
		assetCode:           opts.AssetCode,
		issuerAccountSecret: opts.IssuerAccountSecret,
		auroraClient:       opts.auroraClient(),
		auroraURL:          opts.AuroraURL,
		networkPassphrase:   opts.NetworkPassphrase,
		paymentAmount:       opts.FriendbotPaymentAmount,
	}.ServeHTTP)
	mux.Post("/tx-approve", txApproveHandler{
		assetCode:         opts.AssetCode,
		issuerKP:          issuerKP,
		auroraClient:     opts.auroraClient(),
		networkPassphrase: opts.NetworkPassphrase,
		db:                db,
		kycThreshold:      parsedKYCRequiredPaymentThreshold,
		baseURL:           opts.BaseURL,
	}.ServeHTTP)
	mux.Route("/kyc-status", func(mux chi.Router) {
		mux.Post("/{callback_id}", kycstatus.PostHandler{
			DB: db,
		}.ServeHTTP)
		mux.Get("/{diamcircle_address_or_callback_id}", kycstatus.GetDetailHandler{
			DB: db,
		}.ServeHTTP)
		mux.Delete("/{diamcircle_address}", kycstatus.DeleteHandler{
			DB: db,
		}.ServeHTTP)
	})

	return mux
}

func (opts Options) auroraClient() auroraclient.ClientInterface {
	return &auroraclient.Client{
		AuroraURL: opts.AuroraURL,
		HTTP:       &http.Client{Timeout: 30 * time.Second},
	}
}

func buildURLString(baseURL, endpoint string) string {
	URL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "Unable to parse URL: %s", baseURL))
	}
	URL.Path = path.Join(URL.Path, endpoint)
	return URL.String()
}

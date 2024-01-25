package main

import (
	"database/sql"
	"fmt"
	stdhttp "net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"github.com/diamcircle/go/services/friendbot/internal"
	"github.com/diamcircle/go/support/app"
	"github.com/diamcircle/go/support/config"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/http"
	"github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/support/render/problem"
)

// Config represents the configuration of a friendbot server
type Config struct {
	Port                   int         `toml:"port" valid:"required"`
	FriendbotSecret        string      `toml:"friendbot_secret" valid:"required"`
	NetworkPassphrase      string      `toml:"network_passphrase" valid:"required"`
	AuroraURL             string      `toml:"aurora_url" valid:"required"`
	StartingBalance        string      `toml:"starting_balance" valid:"required"`
	TLS                    *config.TLS `valid:"optional"`
	NumMinions             int         `toml:"num_minions" valid:"optional"`
	BaseFee                int64       `toml:"base_fee" valid:"optional"`
	MinionBatchSize        int         `toml:"minion_batch_size" valid:"optional"`
	SubmitTxRetriesAllowed int         `toml:"submit_tx_retries_allowed" valid:"optional"`
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "friendbot",
		Short: "friendbot for the Diamcircle Test Network",
		Long:  "Client-facing API server for the friendbot service on the Diamcircle Test Network",
		Run:   run,
	}

	rootCmd.PersistentFlags().String("conf", "./friendbot.cfg", "config file path")
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

	fb, err := initFriendbot(cfg.FriendbotSecret, cfg.NetworkPassphrase, cfg.AuroraURL, cfg.StartingBalance,
		cfg.NumMinions, cfg.BaseFee, cfg.MinionBatchSize, cfg.SubmitTxRetriesAllowed)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	router := initRouter(fb)
	registerProblems()

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)

	http.Run(http.Config{
		ListenAddr: addr,
		Handler:    router,
		TLS:        cfg.TLS,
		OnStarting: func() {
			log.Infof("starting friendbot server - %s", app.Version())
			log.Infof("listening on %s", addr)
		},
	})
}

func initRouter(fb *internal.Bot) *chi.Mux {
	mux := http.NewAPIMux(log.DefaultLogger)

	handler := &internal.FriendbotHandler{Friendbot: fb}
	mux.Get("/", handler.Handle)
	mux.Post("/", handler.Handle)
	mux.NotFound(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		problem.Render(r.Context(), w, problem.NotFound)
	}))

	return mux
}

func registerProblems() {
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)
}

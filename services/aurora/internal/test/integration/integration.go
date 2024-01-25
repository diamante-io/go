//lint:file-ignore U1001 Ignore all unused code, this is only used in tests.
package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	sdk "go/clients/auroraclient"
	"go/clients/diamcirclecore"
	"go/keypair"
	proto "go/protocols/aurora"
	aurora "go/services/aurora/internal"
	"go/support/db/dbtest"
	"go/support/errors"
	"go/txnbuild"
	"go/xdr"
)

const (
	StandaloneNetworkPassphrase    = "Standalone Network ; February 2017"
	diamcircleCorePostgresPassword = "mysecretpassword"
	adminPort                      = 6060
	diamcircleCorePort             = 11626
	diamcircleCorePostgresPort     = 5641
	historyArchivePort             = 1570
)

var (
	RunWithCaptiveCore = os.Getenv("aurora_INTEGRATION_ENABLE_CAPTIVE_CORE") != ""
)

type Config struct {
	PostgresURL           string
	ProtocolVersion       uint32
	SkipContainerCreation bool
	CoreDockerImage       string

	// Weird naming here because bools default to false, but we want to start
	// aurora by default.
	SkipauroraStart bool

	// If you want to override the default parameters passed to aurora, you can
	// set this map accordingly. All of them are passed along as --k=v, but if
	// you pass an empty value, the parameter will be dropped. (Note that you
	// should exclude the prepending `--` from keys; this is for compatibility
	// with the constant names in flags.go)
	//
	// You can also control the environmental variables in a similar way, but
	// note that CLI args take precedence over envvars, so set the corresponding
	// CLI arg empty.
	auroraParameters  map[string]string
	auroraEnvironment map[string]string
}

type CaptiveConfig struct {
	binaryPath string
	configPath string
}

type Test struct {
	t *testing.T

	composePath string

	config       Config
	coreConfig   CaptiveConfig
	auroraConfig aurora.Config
	environment  *EnvironmentManager

	auroraClient *sdk.Client
	coreClient   *diamcirclecore.Client

	app           *aurora.App
	appStopped    chan struct{}
	shutdownOnce  sync.Once
	shutdownCalls []func()
	masterKey     *keypair.Full
	passPhrase    string
}

func NewTestForRemoteaurora(t *testing.T, auroraURL string, passPhrase string, masterKey *keypair.Full) *Test {
	return &Test{
		t:            t,
		auroraClient: &sdk.Client{auroraURL: auroraURL},
		masterKey:    masterKey,
		passPhrase:   passPhrase,
	}
}

// NewTest starts a new environment for integration test at a given
// protocol version and blocks until aurora starts ingesting.
//
// Skips the test if aurora_INTEGRATION_TESTS env variable is not set.
//
// WARNING: This requires Docker Compose installed.
func NewTest(t *testing.T, config Config) *Test {
	if os.Getenv("aurora_INTEGRATION_TESTS") == "" {
		t.Skip("skipping integration test: aurora_INTEGRATION_TESTS not set")
	}

	composePath := findDockerComposePath()
	i := &Test{
		t:           t,
		config:      config,
		composePath: composePath,
		passPhrase:  StandaloneNetworkPassphrase,
		environment: NewEnvironmentManager(),
	}

	i.configureCaptiveCore()

	// Only run diamcircle Core container and its dependencies.
	i.runComposeCommand("up", "--detach", "--quiet-pull", "--no-color", "core")
	i.prepareShutdownHandlers()
	i.coreClient = &diamcirclecore.Client{URL: "http://localhost:" + strconv.Itoa(diamcircleCorePort)}
	i.waitForCore()

	if !config.SkipauroraStart {
		if innerErr := i.Startaurora(); innerErr != nil {
			t.Fatalf("Failed to start aurora: %v", innerErr)
		}

		i.WaitForaurora()
	}

	return i
}

func (i *Test) configureCaptiveCore() {
	// We either test Captive Core through environment variables or through
	// custom aurora parameters.
	if RunWithCaptiveCore {
		composePath := findDockerComposePath()
		i.coreConfig.binaryPath = os.Getenv("CAPTIVE_CORE_BIN")
		i.coreConfig.configPath = filepath.Join(composePath, "captive-core-integration-tests.cfg")
	}

	if value := i.getParameter(
		aurora.diamcircleCoreBinaryPathName,
		"diamcircle_CORE_BINARY_PATH",
	); value != "" {
		i.coreConfig.binaryPath = value
	}
	if value := i.getParameter(
		aurora.CaptiveCoreConfigPathName,
		"CAPTIVE_CORE_CONFIG_PATH",
	); value != "" {
		i.coreConfig.configPath = value
	}
}

func (i *Test) getParameter(argName, envName string) string {
	if value, ok := i.config.auroraEnvironment[envName]; ok {
		return value
	}
	if value, ok := i.config.auroraParameters[argName]; ok {
		return value
	}
	return ""
}

// Runs a docker-compose command applied to the above configs
func (i *Test) runComposeCommand(args ...string) {
	integrationYaml := filepath.Join(i.composePath, "docker-compose.integration-tests.yml")

	cmdline := append([]string{"-f", integrationYaml}, args...)
	cmd := exec.Command("docker-compose", cmdline...)
	if i.config.CoreDockerImage != "" {
		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("CORE_IMAGE=%s", i.config.CoreDockerImage),
		)
	}

	i.t.Log("Running", cmd.Env, cmd.Args)
	out, innerErr := cmd.Output()
	if exitErr, ok := innerErr.(*exec.ExitError); ok {
		fmt.Printf("stdout:\n%s\n", string(out))
		fmt.Printf("stderr:\n%s\n", string(exitErr.Stderr))
	}

	if innerErr != nil {
		i.t.Fatalf("Compose command failed: %v", innerErr)
	}
}

func (i *Test) prepareShutdownHandlers() {
	i.shutdownCalls = append(i.shutdownCalls,
		func() {
			if i.app != nil {
				i.app.Close()
			}
			i.runComposeCommand("rm", "-fvs", "core")
			i.runComposeCommand("rm", "-fvs", "core-postgres")
		},
		i.environment.Restore,
	)

	// Register cleanup handlers (on panic and ctrl+c) so the containers are
	// stopped even if ingestion or testing fails.
	i.t.Cleanup(i.Shutdown)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		i.Shutdown()
		os.Exit(int(syscall.SIGTERM))
	}()
}

func (i *Test) Restartaurora() error {
	i.Stopaurora()

	if err := i.Startaurora(); err != nil {
		return err
	}

	i.WaitForaurora()
	return nil
}

func (i *Test) GetauroraConfig() aurora.Config {
	return i.auroraConfig
}

// Shutdown stops the integration tests and destroys all its associated
// resources. It will be implicitly called when the calling test (i.e. the
// `testing.Test` passed to `New()`) is finished if it hasn't been explicitly
// called before.
func (i *Test) Shutdown() {
	i.shutdownOnce.Do(func() {
		// run them in the opposite order in which they where added
		for callI := len(i.shutdownCalls) - 1; callI >= 0; callI-- {
			i.shutdownCalls[callI]()
		}
	})
}

func (i *Test) Startaurora() error {
	auroraPostgresURL := i.config.PostgresURL
	if auroraPostgresURL == "" {
		postgres := dbtest.Postgres(i.t)
		i.shutdownCalls = append(i.shutdownCalls, func() {
			// FIXME: Unfortunately, aurora leaves open sessions behind,
			//        leading to a "database is being accessed by other users"
			//        error when trying to drop it.
			// postgres.Close()
		})
		auroraPostgresURL = postgres.DSN
	}

	config, configOpts := aurora.Flags()
	cmd := &cobra.Command{
		Use:   "aurora",
		Short: "Client-facing API server for the diamcircle network",
		Long:  "Client-facing API server for the diamcircle network.",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			i.app, err = aurora.NewAppFromFlags(config, configOpts)
			if err != nil {
				// Explicitly exit here as that's how these tests are structured for now.
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	// To facilitate custom runs of aurora, we merge a default set of
	// parameters with the tester-supplied ones (if any).
	//
	// TODO: Ideally, we'd be pulling host/port information from the Docker
	//       Compose YAML file itself rather than hardcoding it.
	hostname := "localhost"
	coreBinaryPath := i.coreConfig.binaryPath
	captiveCoreConfigPath := i.coreConfig.configPath

	defaultArgs := map[string]string{
		"diamcircle-core-url": i.coreClient.URL,
		"diamcircle-core-db-url": fmt.Sprintf(
			"postgres://postgres:%s@%s:%d/diamcircle?sslmode=disable",
			diamcircleCorePostgresPassword,
			hostname,
			diamcircleCorePostgresPort,
		),
		"diamcircle-core-binary-path":   coreBinaryPath,
		"captive-core-config-path":      captiveCoreConfigPath,
		"captive-core-http-port":        "21626",
		"enable-captive-core-ingestion": strconv.FormatBool(len(coreBinaryPath) > 0),
		"ingest":                        "true",
		"history-archive-urls":          fmt.Sprintf("http://%s:%d", hostname, historyArchivePort),
		"db-url":                        auroraPostgresURL,
		"network-passphrase":            i.passPhrase,
		"apply-migrations":              "true",
		"admin-port":                    strconv.Itoa(i.AdminPort()),
		"port":                          "8000",
		// due to ARTIFICIALLY_ACCELERATE_TIME_FOR_TESTING
		"checkpoint-frequency": "8",
		"per-hour-rate-limit":  "0", // disable rate limiting
	}

	merged := MergeMaps(defaultArgs, i.config.auroraParameters)
	args := mapToFlags(merged)

	// initialize core arguments
	i.t.Log("aurora command line:", args)
	var env strings.Builder
	for key, value := range i.config.auroraEnvironment {
		env.WriteString(fmt.Sprintf("%s=%s ", key, value))
	}
	i.t.Logf("aurora environmental variables: %s\n", env.String())

	// prepare env
	cmd.SetArgs(args)
	for key, value := range i.config.auroraEnvironment {
		innerErr := i.environment.Add(key, value)
		if innerErr != nil {
			return errors.Wrap(innerErr, fmt.Sprintf(
				"failed to set envvar (%s=%s)", key, value))
		}
	}

	var err error
	if err = configOpts.Init(cmd); err != nil {
		return errors.Wrap(err, "cannot initialize params")
	}

	if err = cmd.Execute(); err != nil {
		return errors.Wrap(err, "cannot initialize aurora")
	}

	auroraPort := "8000"
	if port, ok := merged["--port"]; ok {
		auroraPort = port
	}
	i.auroraConfig = *config
	i.auroraClient = &sdk.Client{
		auroraURL: fmt.Sprintf("http://%s:%s", hostname, auroraPort),
	}

	if err = i.app.Ingestion().BuildGenesisState(); err != nil {
		return errors.Wrap(err, "cannot build genesis state")
	}

	done := make(chan struct{})
	go func() {
		i.app.Serve()
		close(done)
	}()
	i.appStopped = done

	return nil
}

// Wait for core to be up and manually close the first ledger
func (i *Test) waitForCore() {
	i.t.Log("Waiting for core to be up...")
	for t := 30 * time.Second; t >= 0; t -= time.Second {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := i.coreClient.Info(ctx)
		cancel()
		if err != nil {
			i.t.Logf("could not obtain info response: %v", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	i.UpgradeProtocol(i.config.ProtocolVersion)

	for t := 0; t < 5; t++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		info, err := i.coreClient.Info(ctx)
		cancel()
		if err != nil || !info.IsSynced() {
			i.t.Logf("Core is still not synced: %v %v", err, info)
			time.Sleep(time.Second)
			continue
		}
		i.t.Log("Core is up.")
		return
	}
	i.t.Fatal("Core could not sync after 30s")
}

// UpgradeProtocol arms Core with upgrade and blocks until protocol is upgraded.
func (i *Test) UpgradeProtocol(version uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	err := i.coreClient.Upgrade(ctx, int(version))
	cancel()
	if err != nil {
		i.t.Fatalf("could not upgrade protocol: %v", err)
	}

	for t := 0; t < 10; t++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		info, err := i.coreClient.Info(ctx)
		cancel()
		if err != nil {
			i.t.Logf("could not obtain info response: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if info.Info.Ledger.Version == int(version) {
			i.t.Logf("Protocol upgraded to: %d", info.Info.Ledger.Version)
			return
		}
		time.Sleep(time.Second)
	}

	i.t.Fatalf("could not upgrade protocol in 10s")
}

func (i *Test) WaitForaurora() {
	for t := 60; t >= 0; t -= 1 {
		time.Sleep(time.Second)

		i.t.Log("Waiting for ingestion and protocol upgrade...")
		root, err := i.auroraClient.Root()
		if err != nil {
			i.t.Logf("could not obtain root response %v", err)
			continue
		}

		if root.auroraSequence < 3 ||
			int(root.auroraSequence) != int(root.IngestSequence) {
			i.t.Logf("aurora ingesting... %v", root)
			continue
		}

		if uint32(root.CurrentProtocolVersion) == i.config.ProtocolVersion {
			i.t.Logf("aurora protocol version matches... %v", root)
			return
		}
	}

	i.t.Fatal("aurora not ingesting...")
}

// Client returns aurora.Client connected to started aurora instance.
func (i *Test) Client() *sdk.Client {
	return i.auroraClient
}

// aurora returns the aurora.App instance for the current integration test
func (i *Test) aurora() *aurora.App {
	return i.app
}

// Stopaurora shuts down the running aurora process
func (i *Test) Stopaurora() {
	i.app.CloseDB()
	i.app.Close()

	// Wait for aurora to shut down completely.
	<-i.appStopped

	i.app = nil
}

// AdminPort returns aurora admin port.
func (i *Test) AdminPort() int {
	return adminPort
}

// Metrics URL returns aurora metrics URL.
func (i *Test) MetricsURL() string {
	return fmt.Sprintf("http://localhost:%d/metrics", i.AdminPort())
}

// Master returns a keypair of the network masterKey account.
func (i *Test) Master() *keypair.Full {
	if i.masterKey != nil {
		return i.masterKey
	}
	return keypair.Master(i.passPhrase).(*keypair.Full)
}

func (i *Test) MasterAccount() txnbuild.Account {
	master, client := i.Master(), i.Client()
	request := sdk.AccountRequest{AccountID: master.Address()}
	account, err := client.AccountDetail(request)
	panicIf(err)
	return &account
}

func (i *Test) CurrentTest() *testing.T {
	return i.t
}

/* Utility functions for easier test case creation. */

// Creates new accounts via the master account.
//
// It funds each account with the given balance and then queries the API to
// find the randomized sequence number for future operations.
//
// Returns: The slice of created keypairs and account objects.
//
// Note: panics on any errors, since we assume that tests cannot proceed without
// this method succeeding.
func (i *Test) CreateAccounts(count int, initialBalance string) ([]*keypair.Full, []txnbuild.Account) {
	client := i.Client()
	master := i.Master()

	pairs := make([]*keypair.Full, count)
	ops := make([]txnbuild.Operation, count)

	// Two paths here: either caller already did some stuff with the master
	// account so we should retrieve the sequence number, or caller hasn't and
	// we start from scratch.
	seq := int64(0)
	request := sdk.AccountRequest{AccountID: master.Address()}
	account, err := client.AccountDetail(request)
	if err == nil {
		seq, err = strconv.ParseInt(account.Sequence, 10, 64) // str -> bigint
		panicIf(err)
	}

	masterAccount := txnbuild.SimpleAccount{
		AccountID: master.Address(),
		Sequence:  seq,
	}

	for i := 0; i < count; i++ {
		pair, _ := keypair.Random()
		pairs[i] = pair

		ops[i] = &txnbuild.CreateAccount{
			SourceAccount: masterAccount.AccountID,
			Destination:   pair.Address(),
			Amount:        initialBalance,
		}
	}

	// Submit transaction, then retrieve new account details.
	_ = i.MustSubmitOperations(&masterAccount, master, ops...)

	accounts := make([]txnbuild.Account, count)
	for i, kp := range pairs {
		request := sdk.AccountRequest{AccountID: kp.Address()}
		account, err := client.AccountDetail(request)
		panicIf(err)

		accounts[i] = &account
	}

	for _, keys := range pairs {
		i.t.Logf("Funded %s (%s) with %s XLM.\n",
			keys.Seed(), keys.Address(), initialBalance)
	}

	return pairs, accounts
}

// Panics on any error establishing a trustline.
func (i *Test) MustEstablishTrustline(
	truster *keypair.Full, account txnbuild.Account, asset txnbuild.Asset,
) (resp proto.Transaction) {
	txResp, err := i.EstablishTrustline(truster, account, asset)
	panicIf(err)
	return txResp
}

// EstablishTrustline works on a given asset for a particular account.
func (i *Test) EstablishTrustline(
	truster *keypair.Full, account txnbuild.Account, asset txnbuild.Asset,
) (proto.Transaction, error) {
	if asset.IsNative() {
		return proto.Transaction{}, nil
	}
	line, err := asset.ToChangeTrustAsset()
	if err != nil {
		return proto.Transaction{}, err
	}
	return i.SubmitOperations(account, truster, &txnbuild.ChangeTrust{
		Line:  line,
		Limit: "2000",
	})
}

// MustCreateClaimableBalance panics on any error creating a claimable balance.
func (i *Test) MustCreateClaimableBalance(
	source *keypair.Full, asset txnbuild.Asset, amount string,
	claimants ...txnbuild.Claimant,
) (claim proto.ClaimableBalance) {
	account := i.MustGetAccount(source)
	_ = i.MustSubmitOperations(&account, source,
		&txnbuild.CreateClaimableBalance{
			Destinations: claimants,
			Asset:        asset,
			Amount:       amount,
		},
	)

	// Ensure it exists in the global list
	balances, err := i.Client().ClaimableBalances(sdk.ClaimableBalanceRequest{})
	panicIf(err)

	claims := balances.Embedded.Records
	if len(claims) == 0 {
		panic(-1)
	}

	claim = claims[len(claims)-1] // latest one
	i.t.Logf("Created claimable balance w/ id=%s", claim.BalanceID)
	return
}

// MustGetAccount panics on any error retrieves an account's details from its
// key. This means it must have previously been funded.
func (i *Test) MustGetAccount(source *keypair.Full) proto.Account {
	client := i.Client()
	account, err := client.AccountDetail(sdk.AccountRequest{AccountID: source.Address()})
	panicIf(err)
	return account
}

// MustSubmitOperations submits a signed transaction from an account with
// standard options.
//
// Namely, we set the standard fee, time bounds, etc. to "non-production"
// defaults that work well for tests.
//
// Most transactions only need one signer, so see the more verbose
// `MustSubmitOperationsWithSigners` below for multi-sig transactions.
//
// Note: We assume that transaction will be successful here so we panic in case
// of all errors. To allow failures, use `SubmitOperations`.
func (i *Test) MustSubmitOperations(
	source txnbuild.Account, signer *keypair.Full, ops ...txnbuild.Operation,
) proto.Transaction {
	tx, err := i.SubmitOperations(source, signer, ops...)
	panicIf(err)
	return tx
}

func (i *Test) SubmitOperations(
	source txnbuild.Account, signer *keypair.Full, ops ...txnbuild.Operation,
) (proto.Transaction, error) {
	return i.SubmitMultiSigOperations(source, []*keypair.Full{signer}, ops...)
}

func (i *Test) SubmitMultiSigOperations(
	source txnbuild.Account, signers []*keypair.Full, ops ...txnbuild.Operation,
) (proto.Transaction, error) {
	tx, err := i.CreateSignedTransaction(source, signers, ops...)
	if err != nil {
		return proto.Transaction{}, err
	}
	return i.Client().SubmitTransaction(tx)
}

func (i *Test) MustSubmitMultiSigOperations(
	source txnbuild.Account, signers []*keypair.Full, ops ...txnbuild.Operation,
) proto.Transaction {
	tx, err := i.SubmitMultiSigOperations(source, signers, ops...)
	panicIf(err)
	return tx
}

func (i *Test) CreateSignedTransaction(
	source txnbuild.Account, signers []*keypair.Full, ops ...txnbuild.Operation,
) (*txnbuild.Transaction, error) {
	txParams := txnbuild.TransactionParams{
		SourceAccount:        source,
		Operations:           ops,
		BaseFee:              txnbuild.MinBaseFee,
		Timebounds:           txnbuild.NewInfiniteTimeout(),
		IncrementSequenceNum: true,
		EnableMuxedAccounts:  true,
	}

	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return nil, err
	}

	for _, signer := range signers {
		tx, err = tx.Sign(i.passPhrase, signer)
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}

func (i *Test) GetCurrentCoreLedgerSequence() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	info, err := i.coreClient.Info(ctx)
	if err != nil {
		return 0, err
	}
	return info.Info.Ledger.Num, nil
}

// LogFailedTx is a convenience function to provide verbose information about a
// failing transaction to the test output log, if it's expected to succeed.
func (i *Test) LogFailedTx(txResponse proto.Transaction, auroraResult error) {
	t := i.CurrentTest()
	assert.NoErrorf(t, auroraResult, "Submitting the transaction failed")
	if prob := sdk.GetError(auroraResult); prob != nil {
		t.Logf("  problem: %s\n", prob.Problem.Detail)
		t.Logf("  extras: %s\n", prob.Problem.Extras["result_codes"])
		return
	}

	var txResult xdr.TransactionResult
	err := xdr.SafeUnmarshalBase64(txResponse.ResultXdr, &txResult)
	assert.NoErrorf(t, err, "Unmarshalling transaction failed.")
	assert.Equalf(t, xdr.TransactionResultCodeTxSuccess, txResult.Result.Code,
		"Transaction doesn't have success code.")
}

func (i *Test) GetPassPhrase() string {
	return i.passPhrase
}

// Cluttering code with if err != nil is absolute nonsense.
func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// findDockerComposePath performs a best-effort attempt to find the project's
// Docker Compose files.
func findDockerComposePath() string {
	// Lets you check if a particular directory contains a file.
	directoryContainsFilename := func(dir string, filename string) bool {
		files, innerErr := ioutil.ReadDir(dir)
		panicIf(innerErr)

		for _, file := range files {
			if file.Name() == filename {
				return true
			}
		}

		return false
	}

	current, err := os.Getwd()
	panicIf(err)

	//
	// We have a primary and backup attempt for finding the necessary docker
	// files: via $GOPATH and via local directory traversal.
	//

	if gopath := os.Getenv("GOPATH"); gopath != "" {
		monorepo := filepath.Join(gopath, "src", "github.com", "diamcircle", "go")
		if _, err = os.Stat(monorepo); !os.IsNotExist(err) {
			current = monorepo
		}
	}

	// In either case, we try to walk up the tree until we find "go.mod",
	// which we hope is the root directory of the project.
	for !directoryContainsFilename(current, "go.mod") {
		current, err = filepath.Abs(filepath.Join(current, ".."))

		// FIXME: This only works on *nix-like systems.
		if err != nil || filepath.Base(current)[0] == filepath.Separator {
			fmt.Println("Failed to establish project root directory.")
			panic(err)
		}
	}

	// Directly jump down to the folder that should contain the configs
	return filepath.Join(current, "services", "aurora", "docker")
}

// MergeMaps returns a new map which contains the keys and values of *all* input
// maps, overwriting earlier values with later values on duplicate keys.
func MergeMaps(maps ...map[string]string) map[string]string {
	merged := map[string]string{}
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

// mapToFlags will convert a map of parameters into an array of CLI args (i.e.
// in the form --key=value). Note that an empty value for a key means to drop
// the parameter.
func mapToFlags(params map[string]string) []string {
	args := make([]string, 0, len(params))
	for key, value := range params {
		if value == "" {
			continue
		}

		args = append(args, fmt.Sprintf("--%s=%s", key, value))
	}
	return args
}

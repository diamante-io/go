package actions

import (
	"net/http"
	"net/url"

	"go/protocols/aurora"
	"go/services/aurora/internal/ledger"
	"go/services/aurora/internal/resourceadapter"
)

type GetRootHandler struct {
	LedgerState *ledger.State
	CoreStateGetter
	NetworkPassphrase string
	FriendbotURL      *url.URL
	auroraVersion     string
}

func (handler GetRootHandler) GetResource(w HeaderWriter, r *http.Request) (interface{}, error) {
	var res aurora.Root
	templates := map[string]string{
		"accounts":           AccountsQuery{}.URITemplate(),
		"claimableBalances":  ClaimableBalancesQuery{}.URITemplate(),
		"liquidityPools":     LiquidityPoolsQuery{}.URITemplate(),
		"offers":             OffersQuery{}.URITemplate(),
		"strictReceivePaths": StrictReceivePathsQuery{}.URITemplate(),
		"strictSendPaths":    FindFixedPathsQuery{}.URITemplate(),
	}
	coreState := handler.GetCoreState()
	resourceadapter.PopulateRoot(
		r.Context(),
		&res,
		handler.LedgerState.CurrentStatus(),
		handler.auroraVersion,
		coreState.CoreVersion,
		handler.NetworkPassphrase,
		coreState.CurrentProtocolVersion,
		coreState.CoreSupportedProtocolVersion,
		handler.FriendbotURL,
		templates,
	)
	return res, nil
}

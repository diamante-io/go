package actions

import (
	"net/http"
	"net/url"

	"github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	"github.com/diamcircle/go/services/aurora/internal/resourceadapter"
)

type GetRootHandler struct {
	LedgerState *ledger.State
	CoreStateGetter
	NetworkPassphrase string
	FriendbotURL      *url.URL
	AuroraVersion    string
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
		handler.AuroraVersion,
		coreState.CoreVersion,
		handler.NetworkPassphrase,
		coreState.CurrentProtocolVersion,
		coreState.CoreSupportedProtocolVersion,
		handler.FriendbotURL,
		templates,
	)
	return res, nil
}

package serve

import (
	"net/http"

	"github.com/diamcircle/go/exp/services/recoverysigner/internal/account"
	"github.com/diamcircle/go/exp/services/recoverysigner/internal/serve/auth"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/support/http/httpdecode"
	supportlog "github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/support/render/httpjson"
	"github.com/diamcircle/go/txnbuild"
)

type accountSignHandler struct {
	Logger                *supportlog.Entry
	SigningKeys           []*keypair.Full
	NetworkPassphrase     string
	AccountStore          account.Store
	AllowedSourceAccounts []*keypair.FromAddress
}

type accountSignRequest struct {
	Address        *keypair.FromAddress `path:"address"`
	SigningAddress *keypair.FromAddress `path:"signing-address"`
	Transaction    string               `json:"transaction" form:"transaction"`
}

type accountSignResponse struct {
	Signature         string `json:"signature"`
	NetworkPassphrase string `json:"network_passphrase"`
}

func (h accountSignHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check that the client is authenticated in some bare minimum way.
	claims, _ := auth.FromContext(ctx)
	if claims.Address == "" && claims.PhoneNumber == "" && claims.Email == "" {
		unauthorized.Render(w)
		return
	}

	// Decode request.
	req := accountSignRequest{}
	err := httpdecode.Decode(r, &req)
	if err != nil || req.Address == nil || req.SigningAddress == nil {
		badRequest.Render(w)
		return
	}

	l := h.Logger.Ctx(ctx).
		WithField("account", req.Address.Address())
	if req.SigningAddress != nil {
		l = l.WithField("signingaddress", req.SigningAddress.Address())
	}

	l.Info("Request to sign transaction.")

	var signingKey *keypair.Full
	for _, sk := range h.SigningKeys {
		if req.SigningAddress.Address() == sk.Address() {
			signingKey = sk
			break
		}
	}
	if signingKey == nil {
		l.Info("Signing key not found.")
		notFound.Render(w)
		return
	}

	// Find the account that the request is for.
	acc, err := h.AccountStore.Get(req.Address.Address())
	if err == account.ErrNotFound {
		l.Info("Account not found.")
		notFound.Render(w)
		return
	} else if err != nil {
		l.Error(err)
		serverError.Render(w)
		return
	}

	// Authorized if authenticated as the account.
	authorized := claims.Address == req.Address.Address()
	l.Infof("Authorized with self: %v.", authorized)

	// Authorized if authenticated as an identity registered with the account.
	for _, i := range acc.Identities {
		for _, m := range i.AuthMethods {
			if m.Value != "" && ((m.Type == account.AuthMethodTypeAddress && m.Value == claims.Address) ||
				(m.Type == account.AuthMethodTypePhoneNumber && m.Value == claims.PhoneNumber) ||
				(m.Type == account.AuthMethodTypeEmail && m.Value == claims.Email)) {
				authorized = true
				l.Infof("Authorized with %s.", m.Type)
				break
			}
		}
	}

	l.Infof("Authorized: %v.", authorized)
	if !authorized {
		notFound.Render(w)
		return
	}

	// Decode the request transaction.
	parsed, err := txnbuild.TransactionFromXDR(req.Transaction)
	if err != nil {
		l.WithField("transaction", req.Transaction).
			Info("Parsing transaction failed.")
		badRequest.Render(w)
		return
	}
	tx, ok := parsed.Transaction()
	if !ok {
		l.Info("Transaction is not a simple transaction.")
		badRequest.Render(w)
		return
	}
	hashHex, err := tx.HashHex(h.NetworkPassphrase)
	if err != nil {
		l.Error("Error hashing transaction:", err)
		serverError.Render(w)
		return
	}

	l = l.WithField("transaction_hash", hashHex)

	l.Info("Signing transaction.")

	// Check that the transaction's source account and any operations it
	// contains references only to this account.
	if tx.SourceAccount().AccountID != req.Address.Address() {
		l.Info("Transaction's source account is not the account in the request.")
		badRequest.Render(w)
		return
	}
	for _, op := range tx.Operations() {
		opSourceAccount := op.GetSourceAccount()
		if opSourceAccount == "" {
			continue
		}

		if op.GetSourceAccount() != req.Address.Address() {
			var opHasAllowedAccount bool
			for _, sa := range h.AllowedSourceAccounts {
				if sa.Address() == op.GetSourceAccount() {
					opHasAllowedAccount = true
					break
				}
			}

			if !opHasAllowedAccount {
				l.Info("Operation's source account is not the account in the request and not any account that is configured to be allowed.")
				badRequest.Render(w)
				return
			}
		}
	}

	// Sign the transaction.
	hash, err := tx.Hash(h.NetworkPassphrase)
	if err != nil {
		l.Error("Error hashing transaction:", err)
		serverError.Render(w)
		return
	}
	sig, err := signingKey.SignBase64(hash[:])
	if err != nil {
		l.Error("Error signing transaction:", err)
		serverError.Render(w)
		return
	}

	l.Info("Transaction signed.")

	resp := accountSignResponse{
		Signature:         sig,
		NetworkPassphrase: h.NetworkPassphrase,
	}
	httpjson.Render(w, resp, httpjson.JSON)
}

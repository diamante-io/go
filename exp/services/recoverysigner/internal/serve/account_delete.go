package serve

import (
	"net/http"

	"go/exp/services/recoverysigner/internal/account"
	"go/exp/services/recoverysigner/internal/serve/auth"
	"go/keypair"
	"go/support/http/httpdecode"
	supportlog "go/support/log"
	"go/support/render/httpjson"
)

type accountDeleteHandler struct {
	Logger           *supportlog.Entry
	SigningAddresses []*keypair.FromAddress
	AccountStore     account.Store
}

type accountDeleteRequest struct {
	Address *keypair.FromAddress `path:"address"`
}

func (h accountDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, _ := auth.FromContext(ctx)
	if claims.Address == "" && claims.PhoneNumber == "" && claims.Email == "" {
		unauthorized.Render(w)
		return
	}

	req := accountDeleteRequest{}
	err := httpdecode.Decode(r, &req)
	if err != nil || req.Address == nil {
		badRequest.Render(w)
		return
	}

	l := h.Logger.Ctx(ctx).
		WithField("account", req.Address.Address())

	l.Info("Request to delete account.")

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

	resp := accountResponse{
		Address: acc.Address,
	}
	for _, signingAddress := range h.SigningAddresses {
		resp.Signers = append(resp.Signers, accountResponseSigner{
			Key: signingAddress.Address(),
		})
	}

	// Authorized if authenticated as the account.
	authorized := claims.Address == req.Address.Address()
	l.Infof("Authorized with self: %v.", authorized)

	// Authorized if authenticated as an identity registered with the account.
	for _, i := range acc.Identities {
		respIdentity := accountResponseIdentity{
			Role: i.Role,
		}
		for _, m := range i.AuthMethods {
			if m.Value != "" && ((m.Type == account.AuthMethodTypeAddress && m.Value == claims.Address) ||
				(m.Type == account.AuthMethodTypePhoneNumber && m.Value == claims.PhoneNumber) ||
				(m.Type == account.AuthMethodTypeEmail && m.Value == claims.Email)) {
				respIdentity.Authenticated = true
				authorized = true
				l.Infof("Authorized with %s.", m.Type)
				break
			}
		}

		resp.Identities = append(resp.Identities, respIdentity)
	}

	l.Infof("Authorized: %v.", authorized)
	if !authorized {
		notFound.Render(w)
		return
	}

	l.Info("Deleting account.")

	err = h.AccountStore.Delete(req.Address.Address())
	if err == account.ErrNotFound {
		// It can happen if two authorized users are trying to delete the account at the same time.
		l.Info("Account not found.")
		notFound.Render(w)
		return
	} else if err != nil {
		l.Error("Error deleting account:", err)
		serverError.Render(w)
		return
	}

	l.Info("Deleted account.")

	httpjson.Render(w, resp, httpjson.JSON)
}

package serve

import (
	"net/http"

	"github.com/diamcircle/go/exp/services/recoverysigner/internal/account"
	"github.com/diamcircle/go/exp/services/recoverysigner/internal/serve/auth"
	"github.com/diamcircle/go/keypair"
	supportlog "github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/support/render/httpjson"
)

type accountListHandler struct {
	Logger           *supportlog.Entry
	SigningAddresses []*keypair.FromAddress
	AccountStore     account.Store
}

type accountListResponse struct {
	Accounts []accountResponse `json:"accounts"`
}

func (h accountListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, _ := auth.FromContext(ctx)
	if claims.Address == "" && claims.PhoneNumber == "" && claims.Email == "" {
		unauthorized.Render(w)
		return
	}

	l := h.Logger.Ctx(ctx)

	l.Info("Request to get accounts.")

	resp := accountListResponse{
		Accounts: []accountResponse{},
	}

	// Find accounts matching the authenticated address.
	if claims.Address != "" {
		// Find an account that has that address.
		acc, err := h.AccountStore.Get(claims.Address)
		if err == account.ErrNotFound {
			// Do nothing.
		} else if err != nil {
			l.Error(err)
			serverError.Render(w)
			return
		} else {
			accResp := accountResponse{
				Address: acc.Address,
			}
			for _, signingAddress := range h.SigningAddresses {
				accResp.Signers = append(accResp.Signers, accountResponseSigner{
					Key: signingAddress.Address(),
				})
			}
			for _, i := range acc.Identities {
				accRespIdentity := accountResponseIdentity{
					Role: i.Role,
				}
				accResp.Identities = append(accResp.Identities, accRespIdentity)
			}
			resp.Accounts = append(resp.Accounts, accResp)
			l.WithField("account", acc.Address).
				WithField("auth_method_type", account.AuthMethodTypeAddress).
				Info("Found account with auth method type as self.")
		}

		// Find accounts that have the address listed as an owner or other identity.
		accs, err := h.AccountStore.FindWithIdentityAddress(claims.Address)
		if err != nil {
			l.Error(err)
			serverError.Render(w)
			return
		}
		for _, acc := range accs {
			accResp := accountResponse{
				Address: acc.Address,
			}
			for _, signingAddress := range h.SigningAddresses {
				accResp.Signers = append(accResp.Signers, accountResponseSigner{
					Key: signingAddress.Address(),
				})
			}
			for _, i := range acc.Identities {
				accRespIdentity := accountResponseIdentity{
					Role: i.Role,
				}
				for _, m := range i.AuthMethods {
					if m.Type == account.AuthMethodTypeAddress && m.Value == claims.Address {
						accRespIdentity.Authenticated = true
						break
					}
				}
				accResp.Identities = append(accResp.Identities, accRespIdentity)
			}
			resp.Accounts = append(resp.Accounts, accResp)
			l.WithField("account", acc.Address).
				WithField("auth_method_type", account.AuthMethodTypeAddress).
				Info("Found account with auth method type as identity.")
		}
	}

	// Find accounts matching the authenticated phone number.
	if claims.PhoneNumber != "" {
		accs, err := h.AccountStore.FindWithIdentityPhoneNumber(claims.PhoneNumber)
		if err != nil {
			h.Logger.Error(err)
			serverError.Render(w)
			return
		}
		for _, acc := range accs {
			accResp := accountResponse{
				Address: acc.Address,
			}
			for _, signingAddress := range h.SigningAddresses {
				accResp.Signers = append(accResp.Signers, accountResponseSigner{
					Key: signingAddress.Address(),
				})
			}
			for _, i := range acc.Identities {
				accRespIdentity := accountResponseIdentity{
					Role: i.Role,
				}
				for _, m := range i.AuthMethods {
					if m.Type == account.AuthMethodTypePhoneNumber && m.Value == claims.PhoneNumber {
						accRespIdentity.Authenticated = true
						break
					}
				}
				accResp.Identities = append(accResp.Identities, accRespIdentity)
			}
			resp.Accounts = append(resp.Accounts, accResp)
			l.WithField("account", acc.Address).
				WithField("auth_method_type", account.AuthMethodTypePhoneNumber).
				Info("Found account with auth method type as identity.")
		}
	}

	// Find accounts matching the authenticated email.
	if claims.Email != "" {
		accs, err := h.AccountStore.FindWithIdentityEmail(claims.Email)
		if err != nil {
			h.Logger.Error(err)
			serverError.Render(w)
			return
		}
		for _, acc := range accs {
			accResp := accountResponse{
				Address: acc.Address,
			}
			for _, signingAddress := range h.SigningAddresses {
				accResp.Signers = append(accResp.Signers, accountResponseSigner{
					Key: signingAddress.Address(),
				})
			}
			for _, i := range acc.Identities {
				accRespIdentity := accountResponseIdentity{
					Role: i.Role,
				}
				for _, m := range i.AuthMethods {
					if m.Type == account.AuthMethodTypeEmail && m.Value == claims.Email {
						accRespIdentity.Authenticated = true
						break
					}
				}
				accResp.Identities = append(accResp.Identities, accRespIdentity)
			}
			resp.Accounts = append(resp.Accounts, accResp)
			l.WithField("account", acc.Address).
				WithField("auth_method_type", account.AuthMethodTypeEmail).
				Info("Found account with auth method type as identity.")
		}
	}

	httpjson.Render(w, resp, httpjson.JSON)
}

package serve

import (
	"net/http"

	"go/exp/services/recoverysigner/internal/account"
	"go/exp/services/recoverysigner/internal/serve/auth"
	"go/keypair"
	"go/support/errors"
	"go/support/http/httpdecode"
	supportlog "go/support/log"
	"go/support/render/httpjson"
)

type accountPutHandler struct {
	Logger           *supportlog.Entry
	SigningAddresses []*keypair.FromAddress
	AccountStore     account.Store
}

type accountPutRequest struct {
	Address    *keypair.FromAddress        `path:"address"`
	Identities []accountPutRequestIdentity `json:"identities" form:"identities"`
}

func (r accountPutRequest) Validate() error {
	if len(r.Identities) == 0 {
		return errors.Errorf("no identities provided but at least one is required")
	}
	for _, i := range r.Identities {
		err := i.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

type accountPutRequestIdentity struct {
	Role        string                                `json:"role" form:"role"`
	AuthMethods []accountPutRequestIdentityAuthMethod `json:"auth_methods" form:"auth_methods"`
}

func (i accountPutRequestIdentity) Validate() error {
	if i.Role == "" {
		return errors.Errorf("role is not set but required")
	}
	if len(i.AuthMethods) == 0 {
		return errors.Errorf("auth methods not provided for identity but required")
	}
	for _, am := range i.AuthMethods {
		err := am.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

type accountPutRequestIdentityAuthMethod struct {
	Type  string `json:"type" form:"type"`
	Value string `json:"value" form:"value"`
}

func (am accountPutRequestIdentityAuthMethod) Validate() error {
	if !account.AuthMethodType(am.Type).Valid() {
		return errors.Errorf("auth method type %q unrecognized", am.Type)
	}
	// TODO: Validate auth method values: diamcircle address, phone number and email.
	return nil
}

func (h accountPutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, _ := auth.FromContext(ctx)
	if claims.Address == "" && claims.PhoneNumber == "" && claims.Email == "" {
		unauthorized.Render(w)
		return
	}

	req := accountPutRequest{}
	err := httpdecode.Decode(r, &req)
	if err != nil || req.Address == nil {
		badRequest.Render(w)
		return
	}

	l := h.Logger.Ctx(ctx).
		WithField("account", req.Address.Address())

	l.Info("Request to update account.")

	if req.Validate() != nil {
		l.Info("Request validation failed.")
		badRequest.Render(w)
		return
	}

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
	if !authorized {
		notFound.Render(w)
		return
	}

	authMethodCount := 0
	accWithNewIdentiies := account.Account{
		Address:    req.Address.Address(),
		Identities: []account.Identity{},
	}
	for _, i := range req.Identities {
		accIdentity := account.Identity{
			Role: i.Role,
		}
		for _, m := range i.AuthMethods {
			accIdentity.AuthMethods = append(accIdentity.AuthMethods, account.AuthMethod{
				Type:  account.AuthMethodType(m.Type),
				Value: m.Value,
			})
			authMethodCount++
		}
		accWithNewIdentiies.Identities = append(accWithNewIdentiies.Identities, accIdentity)
	}
	l = l.
		WithField("identities_count", len(accWithNewIdentiies.Identities)).
		WithField("auth_methods_count", authMethodCount)

	err = h.AccountStore.Update(accWithNewIdentiies)
	if err == account.ErrNotFound {
		// It can happen if another authorized user is trying to delete the account at the same time.
		l.Info("Account not found.")
		notFound.Render(w)
		return
	} else if err != nil {
		h.Logger.Error(err)
		serverError.Render(w)
		return
	}

	l.Info("Account updated.")

	resp := accountResponse{
		Address: accWithNewIdentiies.Address,
	}
	for _, signingAddress := range h.SigningAddresses {
		resp.Signers = append(resp.Signers, accountResponseSigner{
			Key: signingAddress.Address(),
		})
	}
	for _, i := range accWithNewIdentiies.Identities {
		resp.Identities = append(resp.Identities, accountResponseIdentity{
			Role: i.Role,
		})
	}

	httpjson.Render(w, resp, httpjson.JSON)
}

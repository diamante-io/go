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

type accountPostHandler struct {
	Logger           *supportlog.Entry
	SigningAddresses []*keypair.FromAddress
	AccountStore     account.Store
}

type accountPostRequest struct {
	Address    *keypair.FromAddress         `path:"address"`
	Identities []accountPostRequestIdentity `json:"identities" form:"identities"`
}

func (r accountPostRequest) Validate() error {
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

type accountPostRequestIdentity struct {
	Role        string                                 `json:"role" form:"role"`
	AuthMethods []accountPostRequestIdentityAuthMethod `json:"auth_methods" form:"auth_methods"`
}

func (i accountPostRequestIdentity) Validate() error {
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

type accountPostRequestIdentityAuthMethod struct {
	Type  string `json:"type" form:"type"`
	Value string `json:"value" form:"value"`
}

func (am accountPostRequestIdentityAuthMethod) Validate() error {
	if !account.AuthMethodType(am.Type).Valid() {
		return errors.Errorf("auth method type %q unrecognized", am.Type)
	}
	// TODO: Validate auth method values: diamcircle address, phone number and email.
	return nil
}

func (h accountPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, _ := auth.FromContext(ctx)
	if claims.Address == "" {
		unauthorized.Render(w)
		return
	}

	req := accountPostRequest{}
	err := httpdecode.Decode(r, &req)
	if err != nil || req.Address == nil {
		badRequest.Render(w)
		return
	}

	l := h.Logger.Ctx(ctx).
		WithField("account", req.Address.Address())

	l.Info("Request to register account.")

	if req.Address.Address() != claims.Address {
		l.WithField("address", claims.Address).
			Info("Not authorized as self, authorized as other address.")
		unauthorized.Render(w)
		return
	}

	if req.Validate() != nil {
		l.Info("Request validation failed.")
		badRequest.Render(w)
		return
	}

	authMethodCount := 0
	acc := account.Account{
		Address: req.Address.Address(),
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
		acc.Identities = append(acc.Identities, accIdentity)
	}
	l = l.
		WithField("identities_count", len(acc.Identities)).
		WithField("auth_methods_count", authMethodCount)

	err = h.AccountStore.Add(acc)
	if err == account.ErrAlreadyExists {
		l.Info("Account already registered.")
		conflict.Render(w)
		return
	} else if err != nil {
		l.Error(err)
		serverError.Render(w)
		return
	}

	l.Info("Account registered.")

	resp := accountResponse{
		Address: acc.Address,
	}
	for _, signingAddress := range h.SigningAddresses {
		resp.Signers = append(resp.Signers, accountResponseSigner{
			Key: signingAddress.Address(),
		})
	}
	for _, i := range acc.Identities {
		respIdentity := accountResponseIdentity{
			Role: i.Role,
		}
		resp.Identities = append(resp.Identities, respIdentity)
	}
	httpjson.Render(w, resp, httpjson.JSON)
}

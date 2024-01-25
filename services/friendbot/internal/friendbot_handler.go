package internal

import (
	"net/http"
	"net/url"

	"go/protocols/aurora"
	"go/strkey"
	"go/support/render/hal"
	"go/support/render/problem"
)

// FriendbotHandler causes an account at `Address` to be created.
type FriendbotHandler struct {
	Friendbot *Bot
}

// Handle is a method that implements http.HandlerFunc
func (handler *FriendbotHandler) Handle(w http.ResponseWriter, r *http.Request) {
	accountExistsProblem := problem.BadRequest
	accountExistsProblem.Detail = ErrAccountExists.Error()
	problem.RegisterError(ErrAccountExists, accountExistsProblem)

	result, err := handler.doHandle(r)
	if err != nil {
		problem.Render(r.Context(), w, err)
		return
	}

	hal.Render(w, *result)
}

// doHandle is just a convenience method that returns the object to be rendered
func (handler *FriendbotHandler) doHandle(r *http.Request) (*aurora.Transaction, error) {
	err := r.ParseForm()
	if err != nil {
		p := problem.BadRequest
		p.Detail = "Request parameters are not escaped or incorrectly formatted."
		return nil, &p
	}

	address, err := handler.loadAddress(r)
	if err != nil {
		return nil, problem.MakeInvalidFieldProblem("addr", err)
	}
	return handler.Friendbot.Pay(address)
}

func (handler *FriendbotHandler) loadAddress(r *http.Request) (string, error) {
	address := r.Form.Get("addr")
	unescaped, err := url.QueryUnescape(address)
	if err != nil {
		return unescaped, err
	}

	_, err = strkey.Decode(strkey.VersionByteAccountID, unescaped)
	return unescaped, err
}

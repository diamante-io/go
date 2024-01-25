package federation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go/address"
	proto "go/protocols/federation"
	"go/support/log"

	"github.com/pkg/errors"
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	typ := r.URL.Query().Get("type")
	q := r.URL.Query().Get("q")

	if typ != "forward" && q == "" {
		h.writeJSON(w, ErrorResponse{
			Code:    "invalid_request",
			Message: "q parameter is blank",
		}, http.StatusBadRequest)
		return
	}

	switch typ {
	case "name":
		h.lookupByName(w, r, q)
	case "id":
		h.lookupByID(w, r, q)
	case "forward":
		h.lookupByForward(w, r.URL.Query())
	case "txid":
		h.failNotImplemented(w, "txid type queries are not supported")
	default:
		h.writeJSON(w, ErrorResponse{
			Code:    "invalid_request",
			Message: fmt.Sprintf("invalid type: '%s'", typ),
		}, http.StatusBadRequest)
	}
}

func (h *Handler) failNotFound(w http.ResponseWriter) {
	h.writeJSON(w, ErrorResponse{
		Code:    "not_found",
		Message: "Account not found",
	}, http.StatusNotFound)
}

func (h *Handler) failNotImplemented(w http.ResponseWriter, msg string) {
	h.writeJSON(w, ErrorResponse{
		Code:    "not_implemented",
		Message: msg,
	}, http.StatusNotImplemented)
}

func (h *Handler) lookupByID(w http.ResponseWriter, r *http.Request, q string) {
	rd, ok := h.Driver.(ReverseDriver)

	if !ok {
		h.failNotImplemented(w, "id type queries are not supported")
		return
	}

	// TODO: validate that `q` is a strkey encoded address

	rec, err := rd.LookupReverseRecord(r.Context(), q)
	if err != nil {
		h.writeError(w, errors.Wrap(err, "lookup record"))
		return
	}

	if rec == nil {
		h.failNotFound(w)
		return
	}

	h.writeJSON(w, proto.IDResponse{
		Address: address.New(rec.Name, rec.Domain),
	}, http.StatusOK)
}

func (h *Handler) lookupByName(w http.ResponseWriter, r *http.Request, q string) {
	name, domain, err := address.Split(q)
	if err != nil {
		h.writeJSON(w, ErrorResponse{
			Code:    "invalid_query",
			Message: "Please use an address of the form name*domain.com",
		}, http.StatusBadRequest)
		return
	}

	rec, err := h.Driver.LookupRecord(r.Context(), name, domain)
	if err != nil {
		h.writeError(w, errors.Wrap(err, "lookupByName"))
		return
	}
	if rec == nil {
		h.failNotFound(w)
		return
	}

	h.writeJSON(w, proto.NameResponse{
		AccountID: rec.AccountID,
		Memo:      proto.Memo{rec.Memo},
		MemoType:  rec.MemoType,
	}, http.StatusOK)
}

func (h *Handler) lookupByForward(w http.ResponseWriter, query url.Values) {
	fd, ok := h.Driver.(ForwardDriver)

	if !ok {
		h.failNotImplemented(w, "forward type queries are not supported")
		return
	}

	rec, err := fd.LookupForwardingRecord(query)
	if err != nil {
		h.writeError(w, errors.Wrap(err, "lookupByForward"))
		return
	}
	if rec == nil {
		h.failNotFound(w)
		return
	}

	h.writeJSON(w, proto.NameResponse{
		AccountID: rec.AccountID,
		Memo:      proto.Memo{rec.Memo},
		MemoType:  rec.MemoType,
	}, http.StatusOK)
}

func (h *Handler) writeJSON(
	w http.ResponseWriter,
	obj interface{},
	status int,
) {
	json, err := json.Marshal(obj)

	if err != nil {
		h.writeError(w, errors.Wrap(err, "response marshal"))
		return
	}

	if status == 0 {
		status = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(json)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch err := errors.Cause(err).(type) {
	case ErrorResponse:
		h.writeJSON(w, err, err.StatusCode)
	default:
		log.Error(err)
		http.Error(w, "An internal error occurred", http.StatusInternalServerError)
	}
}

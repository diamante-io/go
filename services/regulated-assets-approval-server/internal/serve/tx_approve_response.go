package serve

import (
	"net/http"

	"github.com/diamcircle/go/support/render/httpjson"
)

type txApprovalResponse struct {
	Error        string     `json:"error,omitempty"`
	Message      string     `json:"message,omitempty"`
	Status       sep8Status `json:"status"`
	StatusCode   int        `json:"-"`
	Tx           string     `json:"tx,omitempty"`
	ActionURL    string     `json:"action_url,omitempty"`
	ActionMethod string     `json:"action_method,omitempty"`
	ActionFields []string   `json:"action_fields,omitempty"`
	Timeout      *int64     `json:"timeout,omitempty"`
}

func (t *txApprovalResponse) Render(w http.ResponseWriter) {
	httpjson.RenderStatus(w, t.StatusCode, t, httpjson.JSON)
}

func NewRejectedTxApprovalResponse(errMessage string) *txApprovalResponse {
	return &txApprovalResponse{
		Status:     sep8StatusRejected,
		Error:      errMessage,
		StatusCode: http.StatusBadRequest,
	}
}

func NewRevisedTxApprovalResponse(tx string) *txApprovalResponse {
	return &txApprovalResponse{
		Status:     sep8StatusRevised,
		Tx:         tx,
		StatusCode: http.StatusOK,
		Message:    "Authorization and deauthorization operations were added.",
	}
}

func NewActionRequiredTxApprovalResponse(message, actionURL string, actionFields []string) *txApprovalResponse {
	return &txApprovalResponse{
		Status:       sep8StatusActionRequired,
		Message:      message,
		ActionMethod: "POST",
		StatusCode:   http.StatusOK,
		ActionURL:    actionURL,
		ActionFields: actionFields,
	}
}

func NewSuccessTxApprovalResponse(tx, message string) *txApprovalResponse {
	return &txApprovalResponse{
		Status:     sep8StatusSuccess,
		Tx:         tx,
		Message:    message,
		StatusCode: http.StatusOK,
	}
}

func NewPendingTxApprovalResponse(message string) *txApprovalResponse {
	timeout := int64(0)
	return &txApprovalResponse{
		Status:     sep8StatusPending,
		Message:    message,
		StatusCode: http.StatusOK,
		Timeout:    &timeout,
	}
}

type sep8Status string

const (
	sep8StatusRejected       sep8Status = "rejected"
	sep8StatusRevised        sep8Status = "revised"
	sep8StatusActionRequired sep8Status = "action_required"
	sep8StatusSuccess        sep8Status = "success"
	sep8StatusPending        sep8Status = "pending"
)

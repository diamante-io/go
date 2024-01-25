package serve

import (
	"fmt"
	"net/http"

	"go/services/regulated-assets-approval-server/internal/serve/httperror"
	"go/strkey"
	"go/support/errors"
	"go/support/log"
)

type diamcircleTOMLHandler struct {
	assetCode         string
	approvalServer    string
	issuerAddress     string
	networkPassphrase string
	kycThreshold      int64
}

func (h diamcircleTOMLHandler) validate() error {
	if h.networkPassphrase == "" {
		return errors.New("network passphrase cannot be empty")
	}

	if h.assetCode == "" {
		return errors.New("asset code cannot be empty")
	}

	if h.issuerAddress == "" {
		return errors.New("asset issuer address cannot be empty")
	}

	if !strkey.IsValidEd25519PublicKey(h.issuerAddress) {
		return errors.New("asset issuer address is not a valid public key")
	}

	if h.approvalServer == "" {
		return errors.New("approval server cannot be empty")
	}

	if h.kycThreshold <= 0 {
		return errors.New("kyc threshold cannot be less than or equal to zero")
	}

	return nil
}

func (h diamcircleTOMLHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.validate()
	if err != nil {
		log.Ctx(ctx).Error(errors.Wrap(err, "validating tomlHandler"))
		httperror.InternalServer.Render(rw)
		return
	}

	// Convert kycThreshold value to human readable string; from amount package's int64 5000000000 to 500.00.
	kycThreshold, err := convertAmountToReadableString(h.kycThreshold)
	if err != nil {
		log.Ctx(ctx).Error(errors.Wrap(err, "converting kycThreshold value to human readable string"))
		httperror.InternalServer.Render(rw)
		return
	}

	// Generate toml content.
	fmt.Fprintf(rw, "NETWORK_PASSPHRASE=%q\n", h.networkPassphrase)
	fmt.Fprintf(rw, "[[CURRENCIES]]\n")
	fmt.Fprintf(rw, "code=%q\n", h.assetCode)
	fmt.Fprintf(rw, "issuer=%q\n", h.issuerAddress)
	fmt.Fprintf(rw, "regulated=true\n")
	fmt.Fprintf(rw, "approval_server=%q\n", h.approvalServer)
	fmt.Fprintf(rw, "approval_criteria=\"The approval server currently only accepts payments. The transaction must have exactly one operation of type payment. If the payment amount exceeds %s %s it will need KYC approval if the account hasn’t been previously approved.\"", kycThreshold, h.assetCode)
}

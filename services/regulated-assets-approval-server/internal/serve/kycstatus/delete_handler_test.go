package kycstatus

import (
	"context"
	"net/http"
	"testing"

	"go/keypair"
	"go/services/regulated-assets-approval-server/internal/db/dbtest"
	"go/services/regulated-assets-approval-server/internal/serve/httperror"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler_validate(t *testing.T) {
	// database is nil
	h := DeleteHandler{}
	err := h.validate()
	require.EqualError(t, err, "database cannot be nil")

	// success
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	h = DeleteHandler{DB: conn}
	err = h.validate()
	require.NoError(t, err)
}

func TestDeleteHandler_handle_errors(t *testing.T) {
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	ctx := context.Background()

	h := DeleteHandler{DB: conn}

	// returns "400 - Missing diamcircle address." if no diamcircle address is provided
	in := deleteRequest{}
	err := h.handle(ctx, in)
	require.Equal(t, httperror.NewHTTPError(http.StatusBadRequest, "Missing diamcircle address."), err)

	// returns "404 - Not found." if the provided address could not be found
	accountKP := keypair.MustRandom()
	in = deleteRequest{diamcircleAddress: accountKP.Address()}
	err = h.handle(ctx, in)
	require.Equal(t, httperror.NewHTTPError(http.StatusNotFound, "Not found."), err)
}

func TestDeleteHandler_handle_success(t *testing.T) {
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	ctx := context.Background()

	h := DeleteHandler{DB: conn}

	// tests if the delete handler is really deleting a row from the database
	q := `
		INSERT INTO accounts_kyc_status (diamcircle_address, callback_id, email_address, kyc_submitted_at, approved_at, rejected_at, pending_at)
		VALUES ($1, $2, $3, NOW(), NOW(), NULL, NULL)
	`
	accountKP := keypair.MustRandom()
	callbackID := uuid.New().String()
	emailAddress := "email@approved.com"
	_, err := h.DB.ExecContext(ctx, q, accountKP.Address(), callbackID, emailAddress)
	require.NoError(t, err)

	in := deleteRequest{diamcircleAddress: accountKP.Address()}
	err = h.handle(ctx, in)
	require.NoError(t, err)

	q = `
		SELECT EXISTS(
			SELECT diamcircle_address
			FROM accounts_kyc_status
			WHERE diamcircle_address = $1
		)
	`
	var exists bool
	err = h.DB.QueryRowContext(ctx, q, accountKP.Address()).Scan(&exists)
	require.NoError(t, err)
	require.False(t, exists)
}

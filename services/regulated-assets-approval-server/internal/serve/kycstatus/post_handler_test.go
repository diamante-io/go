package kycstatus

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"github.com/diamcircle/go/services/regulated-assets-approval-server/internal/db/dbtest"
	"github.com/diamcircle/go/services/regulated-assets-approval-server/internal/serve/httperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostHandler_validate(t *testing.T) {
	// database is nil
	h := PostHandler{}
	err := h.validate()
	require.EqualError(t, err, "database cannot be nil")

	// success
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	h = PostHandler{
		DB: conn,
	}
	err = h.validate()
	require.NoError(t, err)
}

func TestIsKYCRejected(t *testing.T) {
	in := kycPostRequest{
		EmailAddress: "test@email.com",
	}
	isRejected := in.isKYCRejected()
	assert.False(t, isRejected)

	// emails starting with "x" should be rejected
	in = kycPostRequest{
		EmailAddress: "xtest@email.com",
	}
	isRejected = in.isKYCRejected()
	assert.True(t, isRejected)
}

func TestIsKYCPending(t *testing.T) {
	in := kycPostRequest{
		EmailAddress: "test@email.com",
	}
	isPending := in.isKYCPending()
	assert.False(t, isPending)

	// emails starting with "y" should be marked as pending
	in = kycPostRequest{
		EmailAddress: "ytest@email.com",
	}
	isPending = in.isKYCPending()
	assert.True(t, isPending)
}

func TestBuildUpdateKYCQuery(t *testing.T) {
	// test rejected query
	in := kycPostRequest{
		CallbackID:   "9999999999-9999",
		EmailAddress: "xtest@email.com",
	}
	query, args := in.buildUpdateKYCQuery()
	expectedQuery := "WITH updated_row AS (UPDATE accounts_kyc_status SET kyc_submitted_at = NOW(), email_address = $1, rejected_at = NOW(), pending_at = NULL, approved_at = NULL WHERE callback_id = $2 RETURNING * )\n\t\tSELECT EXISTS(\n\t\t\tSELECT * FROM updated_row\n\t\t)\n\t"
	expectedArgs := []interface{}{in.EmailAddress, in.CallbackID}
	require.Equal(t, expectedQuery, query)
	require.Equal(t, expectedArgs, args)

	// test pending query
	in = kycPostRequest{
		CallbackID:   "1234567890-12345",
		EmailAddress: "ytest@email.com",
	}
	query, args = in.buildUpdateKYCQuery()
	expectedQuery = "WITH updated_row AS (UPDATE accounts_kyc_status SET kyc_submitted_at = NOW(), email_address = $1, rejected_at = NULL, pending_at = NOW(), approved_at = NULL WHERE callback_id = $2 RETURNING * )\n\t\tSELECT EXISTS(\n\t\t\tSELECT * FROM updated_row\n\t\t)\n\t"
	expectedArgs = []interface{}{in.EmailAddress, in.CallbackID}
	require.Equal(t, expectedQuery, query)
	require.Equal(t, expectedArgs, args)

	// test approved query
	in = kycPostRequest{
		CallbackID:   "1234567890-12345",
		EmailAddress: "test@email.com",
	}
	query, args = in.buildUpdateKYCQuery()
	expectedQuery = "WITH updated_row AS (UPDATE accounts_kyc_status SET kyc_submitted_at = NOW(), email_address = $1, rejected_at = NULL, pending_at = NULL, approved_at = NOW() WHERE callback_id = $2 RETURNING * )\n\t\tSELECT EXISTS(\n\t\t\tSELECT * FROM updated_row\n\t\t)\n\t"
	expectedArgs = []interface{}{in.EmailAddress, in.CallbackID}
	require.Equal(t, expectedQuery, query)
	require.Equal(t, expectedArgs, args)
}

func TestPostHandler_handle_error(t *testing.T) {
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	ctx := context.Background()

	handler := PostHandler{DB: conn}

	// missing callbackID
	in := kycPostRequest{}
	kycPostResp, err := handler.handle(ctx, in)
	require.Nil(t, kycPostResp)
	require.Equal(t, httperror.NewHTTPError(http.StatusBadRequest, "Missing callbackID."), err)

	// missing email_address
	in = kycPostRequest{
		CallbackID: "random-callback-id",
	}
	kycPostResp, err = handler.handle(ctx, in)
	require.Nil(t, kycPostResp)
	require.Equal(t, httperror.NewHTTPError(http.StatusBadRequest, "Missing email_address."), err)

	// invalid email_address
	in = kycPostRequest{
		CallbackID:   "random-callback-id",
		EmailAddress: "invalidemail",
	}
	kycPostResp, err = handler.handle(ctx, in)
	require.Nil(t, kycPostResp)
	require.Equal(t, httperror.NewHTTPError(http.StatusBadRequest, "The provided email_address is invalid."), err)

	// no entry found for the given callbackID
	in = kycPostRequest{
		CallbackID:   "random-callback-id",
		EmailAddress: "email@test.com",
	}
	kycPostResp, err = handler.handle(ctx, in)
	require.Nil(t, kycPostResp)
	require.Equal(t, httperror.NewHTTPError(http.StatusNotFound, "Not found."), err)
}

func TestPostHandler_handle_success(t *testing.T) {
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	ctx := context.Background()

	handler := PostHandler{DB: conn}

	rejectedCallbackID := "rejected-callback-id"
	pendingCallbackID := "pending-callback-id"
	approvedCallbackID := "approved-callback-id"
	q := `
		INSERT INTO accounts_kyc_status (diamcircle_address, callback_id)
		VALUES 
			('rejected-address', $1),
			('pending-address', $2),
			('approved-address', $3)
	`

	_, err := conn.DB.ExecContext(ctx, q, rejectedCallbackID, pendingCallbackID, approvedCallbackID)
	require.NoError(t, err)

	// should be rejected as email starts with "x"
	in := kycPostRequest{
		CallbackID:   rejectedCallbackID,
		EmailAddress: "xemail@test.com",
	}
	kycPostResp, err := handler.handle(ctx, in)
	assert.NoError(t, err)
	require.Equal(t, NewKYCStatusPostResponse(), kycPostResp)

	var rejectedAt, pendingAt, approvedAt sql.NullTime
	q = `
		SELECT rejected_at, pending_at, approved_at
		FROM accounts_kyc_status
		WHERE callback_id = $1
	`
	err = conn.DB.QueryRowContext(ctx, q, rejectedCallbackID).Scan(&rejectedAt, &pendingAt, &approvedAt)
	require.NoError(t, err)

	assert.False(t, pendingAt.Valid)
	assert.False(t, approvedAt.Valid)
	require.True(t, rejectedAt.Valid)

	// should be marked as pending as email starts with "y"
	in = kycPostRequest{
		CallbackID:   pendingCallbackID,
		EmailAddress: "yemail@test.com",
	}
	kycPostResp, err = handler.handle(ctx, in)
	assert.NoError(t, err)
	require.Equal(t, NewKYCStatusPostResponse(), kycPostResp)

	err = conn.DB.QueryRowContext(ctx, q, pendingCallbackID).Scan(&rejectedAt, &pendingAt, &approvedAt)
	require.NoError(t, err)

	assert.False(t, rejectedAt.Valid)
	assert.False(t, approvedAt.Valid)
	require.True(t, pendingAt.Valid)

	// should be approved as email doesn't start with "x" nor "y"
	in = kycPostRequest{
		CallbackID:   pendingCallbackID,
		EmailAddress: "email@test.com",
	}
	kycPostResp, err = handler.handle(ctx, in)
	assert.NoError(t, err)
	require.Equal(t, NewKYCStatusPostResponse(), kycPostResp)

	err = conn.DB.QueryRowContext(ctx, q, pendingCallbackID).Scan(&rejectedAt, &pendingAt, &approvedAt)
	require.NoError(t, err)

	assert.False(t, rejectedAt.Valid)
	assert.False(t, pendingAt.Valid)
	require.True(t, approvedAt.Valid)
}

func TestRxEmail(t *testing.T) {
	// Test empty email string.
	assert.NotRegexp(t, RxEmail, "")

	// Test empty prefix.
	assert.NotRegexp(t, RxEmail, "email.com")

	// Test only domain given.
	assert.NotRegexp(t, RxEmail, "@email.com")

	// Test correct email.
	assert.Regexp(t, RxEmail, "t@email.com")
}

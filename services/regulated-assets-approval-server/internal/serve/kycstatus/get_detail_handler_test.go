package kycstatus

import (
	"context"
	"database/sql"
	"net/http"
	"testing"
	"time"

	"go/services/regulated-assets-approval-server/internal/db/dbtest"
	"go/services/regulated-assets-approval-server/internal/serve/httperror"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDetailHandler_validate(t *testing.T) {
	// database is nil
	h := GetDetailHandler{}
	err := h.validate()
	require.EqualError(t, err, "database cannot be nil")

	// success
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	h = GetDetailHandler{DB: conn}
	err = h.validate()
	require.NoError(t, err)
}

func TestTimePointerIfValid(t *testing.T) {
	// invalid sql.NullTime should return nil
	sqlNullTime := sql.NullTime{}
	timePointer := timePointerIfValid(sqlNullTime)
	require.Nil(t, timePointer)

	// a valid sql.NullTime should return a time.Time pointer
	desiredTime := time.Now()
	sqlNullTime = sql.NullTime{
		Valid: true,
		Time:  desiredTime,
	}
	timePointer = timePointerIfValid(sqlNullTime)
	require.Equal(t, &desiredTime, timePointer)
}

func TestGetDetailHandler_handle_error(t *testing.T) {
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	ctx := context.Background()

	handler := GetDetailHandler{DB: conn}

	// empty parameter will trigger a "400 - Missing diamcircle address or callbackID."
	in := getDetailRequest{}
	kycGetResp, err := handler.handle(ctx, in)
	assert.Nil(t, kycGetResp)
	require.Equal(t, httperror.NewHTTPError(http.StatusBadRequest, "Missing diamcircle address or callbackID."), err)

	// nonexistent row will trigger a "404 - Not found.".
	in = getDetailRequest{diamcircleAddressOrCallbackID: "nonexistent"}
	kycGetResp, err = handler.handle(ctx, in)
	assert.Nil(t, kycGetResp)
	require.Equal(t, httperror.NewHTTPError(http.StatusNotFound, "Not found."), err)
}

func TestGetDetailHandler_handle_success(t *testing.T) {
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()
	ctx := context.Background()

	handler := GetDetailHandler{DB: conn}

	// step 1: insert test data into database
	const q = `
		INSERT INTO accounts_kyc_status (diamcircle_address, callback_id, email_address, kyc_submitted_at, approved_at, pending_at, rejected_at, created_at)
		VALUES
			('rejected-address', 'rejected-callback-id', 'xrejected@test.com', $1::timestamptz, NULL, NULL, $1::timestamptz, $4::timestamptz),
			('pending-address', 'pending-callback-id', 'ypending@test.com', $2::timestamptz, NULL, $2::timestamptz, NULL, $4::timestamptz),
			('approved-address', 'approved-callback-id', 'approved@test.com', $3::timestamptz, $3::timestamptz, NULL, NULL, $4::timestamptz)
	`
	rejectedAt := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	pendingAt := time.Now().Add(-1 * time.Hour).UTC().Truncate(time.Second)
	approvedAt := time.Now().UTC().Truncate(time.Second)
	createdAt := time.Now().UTC().Truncate(time.Second)
	_, err := handler.DB.ExecContext(ctx, q, rejectedAt.Format(time.RFC3339), pendingAt.Format(time.RFC3339), approvedAt.Format(time.RFC3339), createdAt.Format(time.RFC3339))
	require.NoError(t, err)

	// step 2.1: retrieve "rejected" entry with diamcircle address
	in := getDetailRequest{diamcircleAddressOrCallbackID: "rejected-address"}
	kycGetResp, err := handler.handle(ctx, in)
	require.NoError(t, err)
	wantKYCGetResponse := kycGetResponse{
		diamcircleAddress: "rejected-address",
		CallbackID:        "rejected-callback-id",
		EmailAddress:      "xrejected@test.com",
		CreatedAt:         &createdAt,
		KYCSubmittedAt:    &rejectedAt,
		RejectedAt:        &rejectedAt,
		PendingAt:         nil,
		ApprovedAt:        nil,
	}
	assert.Equal(t, &wantKYCGetResponse, kycGetResp)

	// step 2.2: retrieve "rejected" entry with callbackID
	in = getDetailRequest{diamcircleAddressOrCallbackID: "rejected-callback-id"}
	kycGetResp, err = handler.handle(ctx, in)
	require.NoError(t, err)
	assert.Equal(t, &wantKYCGetResponse, kycGetResp)

	// step 3.1: retrieve "pending" entry with diamcircle address
	in = getDetailRequest{diamcircleAddressOrCallbackID: "pending-address"}
	kycGetResp, err = handler.handle(ctx, in)
	require.NoError(t, err)
	wantKYCGetResponse = kycGetResponse{
		diamcircleAddress: "pending-address",
		CallbackID:        "pending-callback-id",
		EmailAddress:      "ypending@test.com",
		CreatedAt:         &createdAt,
		KYCSubmittedAt:    &pendingAt,
		RejectedAt:        nil,
		PendingAt:         &pendingAt,
		ApprovedAt:        nil,
	}
	assert.Equal(t, &wantKYCGetResponse, kycGetResp)

	// step 3.2: retrieve "pending" entry with callbackID
	in = getDetailRequest{diamcircleAddressOrCallbackID: "pending-callback-id"}
	kycGetResp, err = handler.handle(ctx, in)
	require.NoError(t, err)
	assert.Equal(t, &wantKYCGetResponse, kycGetResp)

	// step 4.1: retrieve "approved" entry with diamcircle address
	in = getDetailRequest{diamcircleAddressOrCallbackID: "approved-address"}
	kycGetResp, err = handler.handle(ctx, in)
	require.NoError(t, err)
	wantKYCGetResponse = kycGetResponse{
		diamcircleAddress: "approved-address",
		CallbackID:        "approved-callback-id",
		EmailAddress:      "approved@test.com",
		CreatedAt:         &createdAt,
		KYCSubmittedAt:    &approvedAt,
		RejectedAt:        nil,
		PendingAt:         nil,
		ApprovedAt:        &approvedAt,
	}
	assert.Equal(t, &wantKYCGetResponse, kycGetResp)

	// step 4.2: retrieve "approved" entry with callbackID
	in = getDetailRequest{diamcircleAddressOrCallbackID: "approved-callback-id"}
	kycGetResp, err = handler.handle(ctx, in)
	require.NoError(t, err)
	assert.Equal(t, &wantKYCGetResponse, kycGetResp)
}

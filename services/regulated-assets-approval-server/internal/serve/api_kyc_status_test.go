package serve

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go/keypair"
	"go/services/regulated-assets-approval-server/internal/db/dbtest"
	"go/services/regulated-assets-approval-server/internal/serve/kycstatus"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_postKYCStatus(t *testing.T) {
	ctx := context.Background()
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()

	handler := kycstatus.PostHandler{DB: conn}
	m := chi.NewMux()
	m.Post("/kyc-status/{callback_id}", handler.ServeHTTP)

	q := `
		INSERT INTO accounts_kyc_status (diamcircle_address, callback_id)
		VALUES ($1, $2)
	`
	clientKP := keypair.MustRandom()
	callbackID := uuid.New().String()
	_, err := handler.DB.ExecContext(ctx, q, clientKP.Address(), callbackID)
	require.NoError(t, err)

	r := httptest.NewRequest("POST", "/kyc-status/"+callbackID, strings.NewReader(`{"email_address": "email@test.com"}`))
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	type kycStatusPOSTResponse struct {
		Result string `json:"result"`
	}
	var kycStatusPOSTResponseApprove kycStatusPOSTResponse
	err = json.Unmarshal(body, &kycStatusPOSTResponseApprove)
	require.NoError(t, err)
	wantPostResponse := kycStatusPOSTResponse{
		Result: "no_further_action_required",
	}
	require.Equal(t, wantPostResponse, kycStatusPOSTResponseApprove)

	q = `
		SELECT rejected_at, pending_at, approved_at
		FROM accounts_kyc_status
		WHERE diamcircle_address = $1 AND callback_id = $2
	`
	var rejectedAt, pendingAt, approvedAt sql.NullTime
	err = handler.DB.QueryRowContext(ctx, q, clientKP.Address(), callbackID).Scan(&rejectedAt, &pendingAt, &approvedAt)
	require.NoError(t, err)

	assert.True(t, approvedAt.Valid)
	assert.False(t, rejectedAt.Valid)
	assert.False(t, pendingAt.Valid)
}

func TestAPI_getKYCStatus(t *testing.T) {
	ctx := context.Background()
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()

	handler := kycstatus.GetDetailHandler{DB: conn}
	m := chi.NewMux()
	m.Get("/kyc-status/{diamcircle_address_or_callback_id}", handler.ServeHTTP)

	// step 1: insert data into database
	const q = `
		INSERT INTO accounts_kyc_status (diamcircle_address, callback_id, email_address, created_at, kyc_submitted_at, rejected_at, pending_at, approved_at)
		VALUES
			('rejected-diamcircle-address', 'rejected-callback-id', 'xrejected@test.com', $1::timestamptz, $2::timestamptz, $2::timestamptz, NULL, NULL),
			('pending-diamcircle-address', 'pending-callback-id', 'ypending@test.com', $1::timestamptz, $3::timestamptz, NULL, $3::timestamptz, NULL),
			('approved-diamcircle-address', 'approved-callback-id', 'approved@test.com', $1::timestamptz, $4::timestamptz, NULL, NULL, $4::timestamptz)
	`
	rejectedAt := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second).Format(time.RFC3339)
	pendingAt := time.Now().Add(-1 * time.Hour).UTC().Truncate(time.Second).Format(time.RFC3339)
	approvedAt := time.Now().UTC().Truncate(time.Second).Format(time.RFC3339)
	createdAt := time.Now().Add(-1 * time.Hour).UTC().Truncate(time.Second).Format(time.RFC3339)
	_, err := handler.DB.ExecContext(ctx, q, createdAt, rejectedAt, pendingAt, approvedAt)
	require.NoError(t, err)

	// step 2: GET "rejected" response
	r := httptest.NewRequest("GET", "/kyc-status/rejected-diamcircle-address", nil)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody := fmt.Sprintf(`{
		"diamcircle_address": "rejected-diamcircle-address",
		"callback_id": "rejected-callback-id",
		"email_address": "xrejected@test.com",
		"created_at": "%s",
		"kyc_submitted_at": "%s",
		"rejected_at": "%s"
	}`, createdAt, rejectedAt, rejectedAt)
	require.JSONEq(t, wantBody, string(body))

	// step 2: GET "pending" response
	r = httptest.NewRequest("GET", "/kyc-status/pending-diamcircle-address", nil)
	r = r.WithContext(ctx)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody = fmt.Sprintf(`{
		"diamcircle_address": "pending-diamcircle-address",
		"callback_id": "pending-callback-id",
		"email_address": "ypending@test.com",
		"created_at": "%s",
		"kyc_submitted_at": "%s",
		"pending_at": "%s"
	}`, createdAt, pendingAt, pendingAt)
	require.JSONEq(t, wantBody, string(body))

	// step 3: GET "approved" response
	r = httptest.NewRequest("GET", "/kyc-status/approved-diamcircle-address", nil)
	r = r.WithContext(ctx)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody = fmt.Sprintf(`{
		"diamcircle_address": "approved-diamcircle-address",
		"callback_id": "approved-callback-id",
		"email_address": "approved@test.com",
		"created_at": "%s",
		"kyc_submitted_at": "%s",
		"approved_at": "%s"
	}`, createdAt, approvedAt, approvedAt)
	require.JSONEq(t, wantBody, string(body))
}

func TestAPI_deleteKYCStatus(t *testing.T) {
	ctx := context.Background()
	db := dbtest.Open(t)
	defer db.Close()
	conn := db.Open()
	defer conn.Close()

	handler := kycstatus.DeleteHandler{DB: conn}
	m := chi.NewMux()
	m.Delete("/kyc-status/{diamcircle_address}", handler.ServeHTTP)

	q := `
		INSERT INTO accounts_kyc_status (diamcircle_address, callback_id, email_address, kyc_submitted_at, approved_at, rejected_at, pending_at)
		VALUES ($1, $2, $3, NOW(), NOW(), NULL, NULL)
	`
	approveKP := keypair.MustRandom()
	approveCallbackID := uuid.New().String()
	approveEmailAddress := "email@test.com"
	_, err := handler.DB.ExecContext(ctx, q, approveKP.Address(), approveCallbackID, approveEmailAddress)
	require.NoError(t, err)

	r := httptest.NewRequest("DELETE", fmt.Sprintf("/kyc-status/%s", approveKP.Address()), nil)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	wantBody := `{
		"message":"ok"
	}`
	require.JSONEq(t, wantBody, string(body))

	q = `
		SELECT EXISTS(
			SELECT diamcircle_address
			FROM accounts_kyc_status
			WHERE diamcircle_address = $1
		)
	`
	var exists bool
	err = handler.DB.QueryRowContext(ctx, q, approveKP.Address()).Scan(&exists)
	require.NoError(t, err)
	require.False(t, exists)
}

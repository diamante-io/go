package keystore

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go/support/errors"
	"go/support/render/httpjson"
)

func TestPutKeysAPI(t *testing.T) {
	db := openKeystoreDB(t)
	defer db.Close() // drop test db

	conn := db.Open()
	defer conn.Close() // close db connection

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"userID":"test-user"}`)
	}))
	defer ts.Close()

	h := ServeMux(&Service{
		db: conn.DB,
		authenticator: &Authenticator{
			URL:     ts.URL,
			APIType: REST,
		},
	})

	blob := `[{
		"id": "test-id",
		"salt": "test-salt",
		"encrypterName": "test-encrypter-name",
		"encryptedBlob": "test-encryptedblob"
	}]`
	keysBlob := base64.RawURLEncoding.EncodeToString([]byte(blob))
	body, err := json.Marshal(putKeysRequest{KeysBlob: keysBlob})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("PUT", "/keys", bytes.NewReader([]byte(body)))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("PUT %s responded with %s, want %s", req.URL, http.StatusText(rr.Code), http.StatusText(http.StatusOK))
	}
	got := &encryptedKeysData{}
	json.Unmarshal(rr.Body.Bytes(), &got)
	if got == nil {
		t.Error("Expected to receive an encryptedKeysData response but did not")
	}

	verifyKeysBlob(t, got.KeysBlob, keysBlob)

	if got.CreatedAt.Before(time.Now().Add(-time.Hour)) {
		t.Errorf("got CreatedAt=%s, want CreatedAt within the last hour", got.CreatedAt)
	}
}

func TestGetKeysAPI(t *testing.T) {
	db := openKeystoreDB(t)
	defer db.Close() // drop test db

	conn := db.Open()
	defer conn.Close() // close db connection

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"userID":"test-user"}`)
	}))
	defer ts.Close()

	ctx := withUserID(context.Background(), "test-user")
	s := &Service{
		db: conn.DB,
		authenticator: &Authenticator{
			URL:     ts.URL,
			APIType: REST,
		},
	}
	h := ServeMux(s)

	blob := `[{
		"id": "test-id",
		"salt": "test-salt",
		"encrypterName": "test-encrypter-name",
		"encryptedBlob": "test-encryptedblob"
	}]`
	keysBlob := base64.RawURLEncoding.EncodeToString([]byte(blob))
	_, err := json.Marshal(putKeysRequest{KeysBlob: keysBlob})
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.putKeys(ctx, putKeysRequest{KeysBlob: keysBlob})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/keys", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("GET %s responded with %s, want %s", req.URL, http.StatusText(rr.Code), http.StatusText(http.StatusOK))
	}
	got := &encryptedKeysData{}
	json.Unmarshal(rr.Body.Bytes(), &got)
	if got == nil {
		t.Error("Expected to receive an encryptedKeysData response but did not")
	}

	verifyKeysBlob(t, got.KeysBlob, keysBlob)

	if got.CreatedAt.Before(time.Now().Add(-time.Hour)) {
		t.Errorf("got CreatedAt=%s, want CreatedAt within the last hour", got.CreatedAt)
	}

	err = s.deleteKeys(ctx)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("GET %s responded with %s, want %s", req.URL, http.StatusText(rr.Code), http.StatusText(http.StatusNotFound))
	}
}

func TestDeleteKeysAPI(t *testing.T) {
	db := openKeystoreDB(t)
	defer db.Close() // drop test db

	conn := db.Open()
	defer conn.Close() // close db connection

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"userID":"test-user"}`)
	}))
	defer ts.Close()

	s := &Service{
		db: conn.DB,
		authenticator: &Authenticator{
			URL:     ts.URL,
			APIType: REST,
		},
	}
	h := ServeMux(s)

	blob := `[{
		"id": "test-id",
		"salt": "test-salt",
		"encrypterName": "test-encrypter-name",
		"encryptedBlob": "test-encryptedblob"
	}]`
	keysBlob := base64.RawURLEncoding.EncodeToString([]byte(blob))
	_, err := json.Marshal(putKeysRequest{
		KeysBlob: keysBlob,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := withUserID(context.Background(), "test-user")
	_, err = s.putKeys(ctx, putKeysRequest{KeysBlob: keysBlob})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("DELETE", "/keys", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("DELETE %s responded with %s, want %s", req.URL, http.StatusText(rr.Code), http.StatusText(http.StatusOK))
	}

	got := rr.Body.Bytes()
	dr, _ := json.MarshalIndent(httpjson.DefaultResponse, "", "  ")
	if !bytes.Equal(got, dr) {
		t.Errorf("got: %s, expected: %s", got, dr)
	}

	_, err = s.getKeys(ctx)
	if errors.Cause(err) != sql.ErrNoRows {
		t.Errorf("expect the keys blob of the user %s to be deleted", userID(ctx))
	}
}

package keystore

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/render/problem"
)

type encryptedKeysData struct {
	KeysBlob   string     `json:"keysBlob"`
	CreatedAt  time.Time  `json:"createdAt"`
	ModifiedAt *time.Time `json:"modifiedAt,omitempty"`
}

type encryptedKeyData struct {
	ID            string `json:"id"`
	Salt          string `json:"salt"`
	EncrypterName string `json:"encrypterName"`
	EncryptedBlob string `json:"encryptedBlob"`
}

type putKeysRequest struct {
	KeysBlob string `json:"keysBlob"`
}

func (s *Service) putKeys(ctx context.Context, in putKeysRequest) (*encryptedKeysData, error) {
	userID := userID(ctx)
	if userID == "" {
		return nil, probNotAuthorized
	}

	if in.KeysBlob == "" {
		return nil, problem.MakeInvalidFieldProblem("keysBlob", errRequiredField)
	}

	keysData, err := base64.RawURLEncoding.DecodeString(in.KeysBlob)
	if err != nil {
		// TODO: we need to implement a helper function in the
		// support/error package for keeping the stack trace from err
		// and substitude the root error for the one we want for better
		// debugging experience.
		// Thowing away the original err makes it harder for debugging.
		return nil, probInvalidKeysBlob
	}

	var encryptedKeys []encryptedKeyData
	err = json.Unmarshal(keysData, &encryptedKeys)
	if err != nil {
		return nil, probInvalidKeysBlob
	}

	for _, ek := range encryptedKeys {
		if ek.Salt == "" {
			return nil, problem.MakeInvalidFieldProblem("keysBlob", errors.New("salt is required for all the encrypted key data"))
		}
		if ek.EncrypterName == "" {
			return nil, problem.MakeInvalidFieldProblem("keysBlob", errors.New("encrypterName is required for all the encrypted key data"))
		}
		if ek.EncryptedBlob == "" {
			return nil, problem.MakeInvalidFieldProblem("keysBlob", errors.New("encryptedBlob is required for all the encrypted key data"))
		}
		if ek.ID == "" {
			return nil, problem.MakeInvalidFieldProblem("keysBlob", errors.New("id is required for all the encrypted key data"))
		}
	}

	q := `
		INSERT INTO encrypted_keys (user_id, encrypted_keys_data)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET encrypted_keys_data = excluded.encrypted_keys_data, modified_at = NOW()
		RETURNING encrypted_keys_data, created_at, modified_at
	`
	var (
		keysBlob   []byte
		out        encryptedKeysData
		modifiedAt pq.NullTime
	)
	err = s.db.QueryRowContext(ctx, q, userID, keysData).Scan(&keysBlob, &out.CreatedAt, &modifiedAt)
	if err != nil {
		return nil, errors.Wrap(err, "storing keys blob")
	}

	out.KeysBlob = base64.RawURLEncoding.EncodeToString(keysBlob)
	if modifiedAt.Valid {
		out.ModifiedAt = &modifiedAt.Time
	}
	return &out, nil
}

func (s *Service) getKeys(ctx context.Context) (*encryptedKeysData, error) {
	userID := userID(ctx)
	if userID == "" {
		return nil, probNotAuthorized
	}

	q := `
		SELECT encrypted_keys_data, created_at, modified_at
		FROM encrypted_keys
		WHERE user_id = $1
	`
	var (
		keysBlob   []byte
		out        encryptedKeysData
		modifiedAt pq.NullTime
	)
	err := s.db.QueryRowContext(ctx, q, userID).Scan(&keysBlob, &out.CreatedAt, &modifiedAt)
	if err != nil {
		return nil, errors.Wrap(err, "getting keys blob")
	}

	out.KeysBlob = base64.RawURLEncoding.EncodeToString(keysBlob)
	if modifiedAt.Valid {
		out.ModifiedAt = &modifiedAt.Time
	}
	return &out, nil
}

func (s *Service) deleteKeys(ctx context.Context) error {
	userID := userID(ctx)
	if userID == "" {
		return probNotAuthorized
	}

	q := `
		DELETE FROM encrypted_keys
		WHERE user_id = $1
	`
	_, err := s.db.ExecContext(ctx, q, userID)
	return errors.Wrap(err, "deleting keys blob")
}

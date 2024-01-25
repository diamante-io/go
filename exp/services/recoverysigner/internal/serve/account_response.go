package serve

import "time"

type accountResponse struct {
	Address    string                    `json:"address"`
	Identities []accountResponseIdentity `json:"identities"`
	Signers    []accountResponseSigner   `json:"signers"`
}

type accountResponseIdentity struct {
	Role          string `json:"role"`
	Authenticated bool   `json:"authenticated,omitempty"`
}

type accountResponseSigner struct {
	Key     string    `json:"key"`
	AddedAt time.Time `json:"added_at"`
}

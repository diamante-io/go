-- +migrate Up

CREATE TABLE public.accounts_kyc_status (
    diamcircle_address text NOT NULL PRIMARY KEY,
    callback_id text NOT NULL,
    email_address text,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    kyc_submitted_at timestamp with time zone,
    approved_at timestamp with time zone,
    rejected_at timestamp with time zone
);

-- +migrate Down

DROP TABLE public.accounts_kyc_status;
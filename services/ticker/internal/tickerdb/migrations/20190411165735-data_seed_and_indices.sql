
-- +migrate Up
-- Seed Issuer for 'native' assets
INSERT INTO public.issuers (
    public_key,
    name,
    url,
    toml_url,
    federation_server,
    auth_server,
    transfer_server,
    web_auth_endpoint,
    deposit_server,
    org_twitter
) VALUES (
    'native',
    'Diamcircle Development Foundation',
    'http://diamcircle.org',
    '',
    '',
    '',
    '',
    '',
    '',
    'https://twitter.com/diamcircleorg'
);

INSERT INTO public.assets (
    code,
    type,
    num_accounts,
    auth_required,
    auth_revocable,
    amount,
    asset_controlled_by_domain,
    anchor_asset_code,
    anchor_asset_type,
    is_valid,
    validation_error,
    last_valid,
    last_checked,
    display_decimals,
    name,
    description,
    conditions,
    is_asset_anchored,
    fixed_number,
    max_number,
    is_unlimited,
    redemption_instructions,
    collateral_addresses,
    collateral_address_signatures,
    countries,
    status,
    issuer_id,
    issuer_account
) VALUES (
    'XLM',
    'native',
    0,
    FALSE,
    FALSE,
    0.0,
    TRUE,
    '',
    '',
    TRUE,
    '',
    now(),
    now(),
    7,
    'Diamcircle Lumens',
    '',
    '',
    FALSE,
    0,
    0,
    FALSE,
    '',
    '',
    '',
    '',
    '',
    (SELECT id FROM public.issuers WHERE public_key = 'native' AND org_twitter = 'https://twitter.com/diamcircleorg'),
    'native'
);

CREATE INDEX trades_ledger_close_time_idx ON public.trades (ledger_close_time DESC);


-- +migrate Down
DROP INDEX trades_ledger_close_time_idx;

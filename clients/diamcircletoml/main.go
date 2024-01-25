package diamcircletoml

import "net/http"

// diamcircleTomlMaxSize is the maximum size of diamcircle.toml file
const diamcircleTomlMaxSize = 100 * 1024

// WellKnownPath represents the url path at which the diamcircle.toml file should
// exist to conform to the federation protocol.
const WellKnownPath = "/.well-known/diamcircle.toml"

// HTTP represents the http client that a stellertoml resolver uses to make http
// requests.
type HTTP interface {
	Get(url string) (*http.Response, error)
}

// Client represents a client that is capable of resolving a diamcircle.toml file
// using the internet.
type Client struct {
	// HTTP is the http client used when resolving a diamcircle.toml file
	HTTP HTTP

	// UseHTTP forces the client to resolve against servers using plain HTTP.
	// Useful for debugging.
	UseHTTP bool
}

type ClientInterface interface {
	GetdiamcircleToml(domain string) (*Response, error)
	GetdiamcircleTomlByAddress(addr string) (*Response, error)
}

// DefaultClient is a default client using the default parameters
var DefaultClient = &Client{HTTP: http.DefaultClient}

type Principal struct {
	Name                  string `toml:"name"`
	Email                 string `toml:"email"`
	Keybase               string `toml:"keybase"`
	Telegram              string `toml:"telegram"`
	Twitter               string `toml:"twitter"`
	Github                string `toml:"github"`
	IdPhotoHash           string `toml:"id_photo_hash"`
	VerificationPhotoHash string `toml:"verification_photo_hash"`
}

type Currency struct {
	Code                        string   `toml:"code"`
	CodeTemplate                string   `toml:"code_template"`
	Issuer                      string   `toml:"issuer"`
	Status                      string   `toml:"status"`
	DisplayDecimals             int      `toml:"display_decimals"`
	Name                        string   `toml:"name"`
	Desc                        string   `toml:"desc"`
	Conditions                  string   `toml:"conditions"`
	Image                       string   `toml:"image"`
	FixedNumber                 int      `toml:"fixed_number"`
	MaxNumber                   int      `toml:"max_number"`
	IsUnlimited                 bool     `toml:"is_unlimited"`
	IsAssetAnchored             bool     `toml:"is_asset_anchored"`
	AnchorAsset                 string   `toml:"anchor_asset"`
	RedemptionInstructions      string   `toml:"redemption_instructions"`
	CollateralAddresses         []string `toml:"collateral_addresses"`
	CollateralAddressMessages   []string `toml:"collateral_address_messages"`
	CollateralAddressSignatures []string `toml:"collateral_address_signatures"`
	Regulated                   string   `toml:"regulated"`
	ApprovalServer              string   `toml:"APPROVAL_SERVER"`
	ApprovalCriteria            string   `toml:"APPROVAL_CRITERIA"`
}

type Validator struct {
	Alias       string `toml:"ALIAS"`
	DisplayName string `toml:"DISPLAY_NAME"`
	PublicKey   string `toml:"PUBLIC_KEY"`
	Host        string `toml:"HOST"`
	History     string `toml:"HISTORY"`
}

// SEP-1 commit
// https://github.com/diamcircle/diamcircle-protocol/blob/f8993e36fa6b5b8bba1254c21c2174d250af4958/ecosystem/sep-0001.md
type Response struct {
	Version                       string      `toml:"VERSION"`
	NetworkPassphrase             string      `toml:"NETWORK_PASSPHRASE"`
	FederationServer              string      `toml:"FEDERATION_SERVER"`
	AuthServer                    string      `toml:"AUTH_SERVER"`
	TransferServer                string      `toml:"TRANSFER_SERVER"`
	TransferServer0024            string      `toml:"TRANSFER_SERVER_0024"`
	KycServer                     string      `toml:"KYC_SERVER"`
	WebAuthEndpoint               string      `toml:"WEB_AUTH_ENDPOINT"`
	SigningKey                    string      `toml:"SIGNING_KEY"`
	auroraUrl                     string      `toml:"aurora_URL"`
	Accounts                      []string    `toml:"ACCOUNTS"`
	UriRequestSigningKey          string      `toml:"URI_REQUEST_SIGNING_KEY"`
	DirectPaymentServer           string      `toml:"DIRECT_PAYMENT_SERVER"`
	OrgName                       string      `toml:"ORG_NAME"`
	OrgDba                        string      `toml:"ORG_DBA"`
	OrgUrl                        string      `toml:"ORG_URL"`
	OrgLogo                       string      `toml:"ORG_LOGO"`
	OrgDescription                string      `toml:"ORG_DESCRIPTION"`
	OrgPhysicalAddress            string      `toml:"ORG_PHYSICAL_ADDRESS"`
	OrgPhysicalAddressAttestation string      `toml:"ORG_PHYSICAL_ADDRESS_ATTESTATION"`
	OrgPhoneNumber                string      `toml:"ORG_PHONE_NUMBER"`
	OrgPhoneNumberAttestation     string      `toml:"ORG_PHONE_NUMBER_ATTESTATION"`
	OrgKeybase                    string      `toml:"ORG_KEYBASE"`
	OrgTwitter                    string      `toml:"ORG_TWITTER"`
	OrgGithub                     string      `toml:"ORG_GITHUB"`
	OrgOfficialEmail              string      `toml:"ORG_OFFICIAL_EMAIL"`
	OrgLicensingAuthority         string      `toml:"ORG_LICENSING_AUTHORITY"`
	OrgLicenseType                string      `toml:"ORG_LICENSE_TYPE"`
	OrgLicenseNumber              string      `toml:"ORG_LICENSE_NUMBER"`
	Principals                    []Principal `toml:"PRINCIPALS"`
	Currencies                    []Currency  `toml:"CURRENCIES"`
	Validators                    []Validator `toml:"VALIDATORS"`
}

// GetdiamcircleToml returns diamcircle.toml file for a given domain
func GetdiamcircleToml(domain string) (*Response, error) {
	return DefaultClient.GetdiamcircleToml(domain)
}

// GetdiamcircleTomlByAddress returns diamcircle.toml file of a domain fetched from a
// given address
func GetdiamcircleTomlByAddress(addr string) (*Response, error) {
	return DefaultClient.GetdiamcircleTomlByAddress(addr)
}

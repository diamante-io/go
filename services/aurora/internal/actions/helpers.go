package actions

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/gorilla/schema"

	"go/services/aurora/internal/assets"
	auroraContext "go/services/aurora/internal/context"
	"go/services/aurora/internal/db2"
	"go/services/aurora/internal/ledger"
	hProblem "go/services/aurora/internal/render/problem"
	"go/services/aurora/internal/toid"
	"go/support/errors"
	"go/support/render/problem"
	"go/xdr"
)

// TODO: move these constants to urlparam.go as we should parse the params with http handlers
// in the upper level package.
const (
	// ParamCursor is a query string param name
	ParamCursor = "cursor"
	// ParamOrder is a query string param name
	ParamOrder = "order"
	// ParamLimit is a query string param name
	ParamLimit = "limit"
	// LastLedgerHeaderName is the header which is set on all endpoints
	LastLedgerHeaderName = "Latest-Ledger"
)

type Opt int

const (
	// DisableCursorValidation disables cursor validation in GetPageQuery
	DisableCursorValidation Opt = iota
)

// HeaderWriter is an interface for setting HTTP response headers
type HeaderWriter interface {
	Header() http.Header
}

// SetLastLedgerHeader sets the Latest-Ledger header
func SetLastLedgerHeader(w HeaderWriter, lastLedger uint32) {
	w.Header().Set(LastLedgerHeaderName, strconv.FormatUint(uint64(lastLedger), 10))
}

// getCursor retrieves a string from either the URLParams, form or query string.
// This method uses the priority (URLParams, Form, Query).
func getCursor(ledgerState *ledger.State, r *http.Request, name string) (string, error) {
	cursor, err := getString(r, name)

	if err != nil {
		return "", err
	}

	if cursor == "now" {
		tid := toid.AfterLedger(ledgerState.CurrentStatus().HistoryLatest)
		cursor = tid.String()
	}

	if lastEventID := r.Header.Get("Last-Event-ID"); lastEventID != "" {
		cursor = lastEventID
	}

	// In case cursor is negative value, return InvalidField error
	cursorInt, err := strconv.Atoi(cursor)
	if err == nil && cursorInt < 0 {
		msg := fmt.Sprintf("the cursor %d is a negative number: ", cursorInt)

		return "", problem.MakeInvalidFieldProblem(
			name,
			errors.New(msg),
		)
	}

	return cursor, nil
}

func checkUTF8(name, value string) error {
	if !utf8.ValidString(value) {
		return problem.MakeInvalidFieldProblem(name, errors.New("invalid value"))
	}
	return nil
}

// getStringFromURLParam retrieves a string from the URLParams.
func getStringFromURLParam(r *http.Request, name string) (string, error) {
	fromURL, ok := getURLParam(r, name)
	if ok {
		ret, err := url.PathUnescape(fromURL)
		if err != nil {
			return "", problem.MakeInvalidFieldProblem(name, err)
		}

		if err := checkUTF8(name, ret); err != nil {
			return "", err
		}
		return ret, nil
	}

	return "", nil
}

// getString retrieves a string from either the URLParams, form or query string.
// This method uses the priority (URLParams, Form, Query).
func getString(r *http.Request, name string) (string, error) {
	fromURL, ok := getURLParam(r, name)
	if ok {
		ret, err := url.PathUnescape(fromURL)
		if err != nil {
			return "", problem.MakeInvalidFieldProblem(name, err)
		}

		if err := checkUTF8(name, ret); err != nil {
			return "", err
		}

		return ret, nil
	}

	fromForm := r.FormValue(name)
	if fromForm != "" {
		if err := checkUTF8(name, fromForm); err != nil {
			return "", err
		}
		return fromForm, nil
	}

	value := r.URL.Query().Get(name)
	if err := checkUTF8(name, value); err != nil {
		return "", err
	}

	return value, nil
}

// getLimit retrieves a uint64 limit from the action parameter of the given
// name. Populates err if the value is not a valid limit.  Uses the provided
// default value if the limit parameter is a blank string.
func getLimit(r *http.Request, name string, def uint64, max uint64) (uint64, error) {
	limit, err := getString(r, name)

	if err != nil {
		return 0, err
	}
	if limit == "" {
		return def, nil
	}

	asI64, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {

		return 0, problem.MakeInvalidFieldProblem(name, errors.New("unparseable value"))
	}

	if asI64 <= 0 {
		err = errors.New("invalid limit: non-positive value provided")
	} else if asI64 > int64(max) {
		err = errors.Errorf("invalid limit: value provided that is over limit max of %d", max)
	}

	if err != nil {
		return 0, problem.MakeInvalidFieldProblem(name, err)
	}

	return uint64(asI64), nil
}

// GetPageQuery is a helper that returns a new db.PageQuery struct initialized
// using the results from a call to GetPagingParams()
func GetPageQuery(ledgerState *ledger.State, r *http.Request, opts ...Opt) (db2.PageQuery, error) {
	disableCursorValidation := false
	for _, opt := range opts {
		if opt == DisableCursorValidation {
			disableCursorValidation = true
		}
	}

	cursor, err := getCursor(ledgerState, r, ParamCursor)
	if err != nil {
		return db2.PageQuery{}, err
	}
	order, err := getString(r, ParamOrder)
	if err != nil {
		return db2.PageQuery{}, err
	}
	limit, err := getLimit(r, ParamLimit, db2.DefaultPageSize, db2.MaxPageSize)
	if err != nil {
		return db2.PageQuery{}, err
	}

	pageQuery, err := db2.NewPageQuery(cursor, !disableCursorValidation, order, limit)
	if err != nil {
		if invalidFieldError, ok := err.(*db2.InvalidFieldError); ok {
			err = problem.MakeInvalidFieldProblem(
				invalidFieldError.Name,
				err,
			)
		} else {
			err = problem.BadRequest
		}

		return db2.PageQuery{}, err
	}

	return pageQuery, nil
}

// GetTransactionID retireves a transaction identifier by attempting to decode an hex-encoded,
// 64-digit lowercase string at the provided name.
func GetTransactionID(r *http.Request, name string) (string, error) {
	value, err := getStringFromURLParam(r, name)
	if err != nil {
		return "", err
	}

	if value != "" {
		if _, err = hex.DecodeString(value); err != nil || len(value) != 64 || strings.ToLower(value) != value {
			return "", problem.MakeInvalidFieldProblem(name, errors.New("invalid hash format"))
		}
	}

	return value, nil
}

// getAccountID retireves an xdr.AccountID by attempting to decode a diamcircle
// address at the provided name.
func getAccountID(r *http.Request, name string) (xdr.AccountId, error) {
	value, err := getString(r, name)
	if err != nil {
		return xdr.AccountId{}, err
	}

	result, err := xdr.AddressToAccountId(value)
	if err != nil {
		return result, problem.MakeInvalidFieldProblem(
			name,
			errors.New("invalid address"),
		)
	}

	return result, nil
}

// getAssetType is a helper that returns a xdr.AssetType by reading a string
func getAssetType(r *http.Request, name string) (xdr.AssetType, error) {
	val, err := getString(r, name)
	if err != nil {
		return xdr.AssetTypeAssetTypeNative, nil
	}

	t, err := assets.Parse(val)
	if err != nil {
		return t, problem.MakeInvalidFieldProblem(
			name,
			err,
		)
	}

	return t, nil
}

// getAsset decodes an asset from the request fields prefixed by `prefix`.  To
// succeed, three prefixed fields must be present: asset_type, asset_code, and
// asset_issuer.
func getAsset(r *http.Request, prefix string) (xdr.Asset, error) {
	var value interface{}
	t, err := getAssetType(r, prefix+"asset_type")
	if err != nil {
		return xdr.Asset{}, err
	}

	switch t {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		a := xdr.AlphaNum4{}
		a.Issuer, err = getAccountID(r, prefix+"asset_issuer")
		if err != nil {
			return xdr.Asset{}, err
		}

		var code string
		code, err = getString(r, prefix+"asset_code")
		if err != nil {
			return xdr.Asset{}, err
		}
		if len(code) > len(a.AssetCode) {
			err := problem.MakeInvalidFieldProblem(
				prefix+"asset_code",
				errors.New("code too long"),
			)
			return xdr.Asset{}, err
		}

		copy(a.AssetCode[:len(code)], []byte(code))
		value = a
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		a := xdr.AlphaNum12{}
		a.Issuer, err = getAccountID(r, prefix+"asset_issuer")
		if err != nil {
			return xdr.Asset{}, err
		}

		var code string
		code, err = getString(r, prefix+"asset_code")
		if err != nil {
			return xdr.Asset{}, err
		}
		if len(code) > len(a.AssetCode) {
			err := problem.MakeInvalidFieldProblem(
				prefix+"asset_code",
				errors.New("code too long"),
			)
			return xdr.Asset{}, err
		}

		copy(a.AssetCode[:len(code)], []byte(code))
		value = a
	}

	result, err := xdr.NewAsset(t, value)
	if err != nil {
		panic(err)
	}

	return result, nil
}

// getURLParam returns the corresponding URL parameter value from the request
// routing context and an additional boolean reflecting whether or not the
// param was found. This is ported from Chi since the Chi version returns ""
// for params not found. This is undesirable since "" also is a valid url param.
// Ref: https://github.com/go-chi/chi/blob/d132b31857e5922a2cc7963f4fcfd8f46b3f2e97/context.go#L69
func getURLParam(r *http.Request, key string) (string, bool) {
	rctx := chi.RouteContext(r.Context())

	if rctx == nil {
		return "", false
	}

	// Return immediately if keys does not match Values
	// This can happen when a named param is not specified.
	// This is a bug in chi: https://github.com/go-chi/chi/issues/426
	if len(rctx.URLParams.Keys) != len(rctx.URLParams.Values) {
		return "", false
	}

	for k := len(rctx.URLParams.Keys) - 1; k >= 0; k-- {
		if rctx.URLParams.Keys[k] == key {
			return rctx.URLParams.Values[k], true
		}
	}

	return "", false
}

// FullURL returns a URL containing the information regarding the original
// request stored in the context.
func FullURL(ctx context.Context) *url.URL {
	url := auroraContext.BaseURL(ctx)
	r := auroraContext.RequestFromContext(ctx)
	if r != nil {
		url.Path = r.URL.Path
		url.RawQuery = r.URL.RawQuery
	}
	return url
}

// Note from chi: it is a good idea to set a Decoder instance as a package
// global, because it caches meta-data about structs, and an instance can be
// shared safely:
var decoder = schema.NewDecoder()

// getParams fills a struct with values read from a request's query parameters.
func getParams(dst interface{}, r *http.Request) error {
	query := r.URL.Query()

	// Merge chi's URLParams with URL Query Params. Given
	// `/accounts/{account_id}/transactions?foo=bar`, chi's URLParams will
	// contain `account_id` and URL Query params will contain `foo`.
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		for _, key := range rctx.URLParams.Keys {
			if key == "*" {
				continue
			}
			param, _ := getURLParam(r, key)
			query.Set(key, param)
		}
	}

	if err := decoder.Decode(dst, query); err != nil {
		for k, e := range err.(schema.MultiError) {
			return problem.NewProblemWithInvalidField(
				problem.BadRequest,
				k,
				getSchemaErrorFieldMessage(k, e),
			)
		}
	}

	if _, err := govalidator.ValidateStruct(dst); err != nil {
		field, message := getErrorFieldMessage(err)
		err = problem.MakeInvalidFieldProblem(
			getSchemaTag(dst, field),
			errors.New(message),
		)

		return err
	}

	if v, ok := dst.(Validateable); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func getSchemaTag(params interface{}, field string) string {
	v := reflect.ValueOf(params).Elem()
	qt := v.Type()
	f, _ := qt.FieldByName(field)
	return f.Tag.Get("schema")
}

// getURIParams returns a list of query parameters for a given query struct
func getURIParams(query interface{}, paginated bool) []string {
	params := getSchemaTags(reflect.ValueOf(query).Elem())
	if paginated {
		pagingParams := []string{
			ParamCursor,
			ParamLimit,
			ParamOrder,
		}
		params = append(params, pagingParams...)
	}
	return params
}

func getURITemplate(query interface{}, basePath string, paginated bool) string {
	return "/" + basePath + "{?" + strings.Join(getURIParams(query, paginated), ",") + "}"
}

func getSchemaTags(v reflect.Value) []string {
	qt := v.Type()
	fields := make([]string, 0, v.NumField())

	for i := 0; i < qt.NumField(); i++ {
		f := qt.Field(i)
		// Query structs can have embedded query structs
		if f.Type.Kind() == reflect.Struct {
			fields = append(fields, getSchemaTags(v.Field(i))...)
		} else {
			tag, ok := f.Tag.Lookup("schema")
			if ok {
				fields = append(fields, tag)
			}
		}
	}

	return fields
}

// validateAssetParams runs multiple checks on an asset query parameter
func validateAssetParams(aType, code, issuer, prefix string) error {
	// If asset type is not present but code or issuer are, then there is a
	// missing parameter and the request is unprocessable.
	if len(aType) == 0 {
		if len(code) > 0 || len(issuer) > 0 {
			return problem.MakeInvalidFieldProblem(
				prefix+"asset_type",
				errors.New("Missing parameter"),
			)
		}

		return nil
	}

	t, err := assets.Parse(aType)
	if err != nil {
		return problem.MakeInvalidFieldProblem(
			prefix+"asset_type",
			err,
		)
	}

	var validLen int
	switch t {
	case xdr.AssetTypeAssetTypeNative:
		// If asset type is native, issuer or code should not be included in the
		// request
		switch {
		case len(code) > 0:
			return problem.MakeInvalidFieldProblem(
				prefix+"asset_code",
				errors.New("native asset does not have a code"),
			)
		case len(issuer) > 0:
			return problem.MakeInvalidFieldProblem(
				prefix+"asset_issuer",
				errors.New("native asset does not have an issuer"),
			)
		}

		return nil
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		validLen = len(xdr.AlphaNum4{}.AssetCode)
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		validLen = len(xdr.AlphaNum12{}.AssetCode)
	}

	codeLen := len(code)
	if codeLen == 0 || codeLen > validLen {
		return problem.MakeInvalidFieldProblem(
			prefix+"asset_code",
			errors.New("Asset code must be 1-12 alphanumeric characters"),
		)
	}

	if len(issuer) == 0 {
		return problem.MakeInvalidFieldProblem(
			prefix+"asset_issuer",
			errors.New("Missing parameter"),
		)
	}

	return nil
}

// validateCursorWithinHistory compares the requested page of data against the
// ledger state of the history database.  In the event that the cursor is
// guaranteed to return no results, we return a 410 GONE http response.
func validateCursorWithinHistory(ledgerState *ledger.State, pq db2.PageQuery) error {
	// an ascending query should never return a gone response:  An ascending query
	// prior to known history should return results at the beginning of history,
	// and an ascending query beyond the end of history should not error out but
	// rather return an empty page (allowing code that tracks the procession of
	// some resource more easily).
	if pq.Order != "desc" {
		return nil
	}

	var cursor int64
	var err error

	// Checking for the presence of "-" to see whether we should use CursorInt64
	// or CursorInt64Pair
	if strings.Contains(pq.Cursor, "-") {
		cursor, _, err = pq.CursorInt64Pair("-")
	} else {
		cursor, err = pq.CursorInt64()
	}

	if err != nil {
		return problem.MakeInvalidFieldProblem("cursor", errors.New("invalid value"))
	}

	elder := toid.New(ledgerState.CurrentStatus().HistoryElder, 0, 0)

	if cursor <= elder.ToInt64() {
		return &hProblem.BeforeHistory
	}

	return nil
}

func countNonEmpty(params ...interface{}) (int, error) {
	count := 0

	for _, param := range params {
		switch param := param.(type) {
		default:
			return 0, errors.Errorf("unexpected type %T", param)
		case int32:
			if param != 0 {
				count++
			}
		case uint32:
			if param != 0 {
				count++
			}
		case int64:
			if param != 0 {
				count++
			}
		case uint64:
			if param != 0 {
				count++
			}
		case string:
			if param != "" {
				count++
			}
		case *xdr.Asset:
			if param != nil {
				count++
			}
		}
	}

	return count, nil
}

func init() {
	decoder.IgnoreUnknownKeys(true)
}

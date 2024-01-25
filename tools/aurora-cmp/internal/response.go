package cmp

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const fileLengthLimit = 100

var findResultMetaXDR = regexp.MustCompile(`"result_meta_xdr":[ ]?"([^"]*)",`)

// removeRegexps contains a list of regular expressions that, when matched,
// will be changed to an empty string. This is done to exclude known
// differences in responses between two Aurora version.
//
// Let's say that next Aurora version adds a new bool field:
// `is_authorized` on account balances list. You want to remove this
// field so it's not reported for each `/accounts/{id}` response.
var removeRegexps = []*regexp.Regexp{
	// regexp.MustCompile(`This (is ){0,1}usually`),
	// Removes joined transaction (join=transactions) added in Aurora 0.19.0.
	// Remove for future versions.
	// regexp.MustCompile(`(?msU)"transaction":\s*{\s*("memo|"_links)[\n\s\S]*][\n\s\S]*}(,\s{9}|,)`),
	// regexp.MustCompile(`\s*"is_authorized": true,`),
	// regexp.MustCompile(`\s*"is_authorized": false,`),
	// regexp.MustCompile(`\s*"successful": true,`),
	// regexp.MustCompile(`\s*"transaction_count": [0-9]+,`),
	// regexp.MustCompile(`\s*"last_modified_ledger": [0-9]+,`),
	// regexp.MustCompile(`\s*"public_key": "G.*",`),
	// regexp.MustCompile(`,\s*"paging_token": ?""`),
	// Removes last_modified_time field, introduced in aurora 1.3.0
	// regexp.MustCompile(`\s*"last_modified_time": ?"[^"]*",`),
}

type replace struct {
	regexp *regexp.Regexp
	repl   string
}

var replaceRegexps = []replace{
	// Offer ID in /offers
	// {regexp.MustCompile(`"id":( ?)([0-9]+)`), `"id":${1}"${2}"`},
	// {regexp.MustCompile(`"timestamp":( ?)([0-9]+)`), `"timestamp":${1}"${2}"`},
	// {regexp.MustCompile(`"trade_count":( ?)([0-9]+)`), `"trade_count":${1}"${2}"`},
	// {regexp.MustCompile(`"type":( ?)"manage_offer",`), `"type":${1}"manage_sell_offer",`},
	// {regexp.MustCompile(`"type":( ?)"path_payment",`), `"type":${1}"path_payment_strict_receive",`},
	// {regexp.MustCompile(`"type":( ?)"create_passive_offer",`), `"type":${1}"create_passive_sell_offer",`},
	// {regexp.MustCompile(
	// 	// Removes paging_token from /accounts/*
	// 	`"data":( ?){([^}]*)},\s*"paging_token":( ?)"([0-9A-Z]*)"`),
	// 	`"data":${1}{${2}},"paging_token":${3}""`,
	// },
	// {regexp.MustCompile(
	// 	// fee_charged is a string since aurora 1.3.0
	// 	`"fee_charged":( ?)([\d]+),`),
	// 	`"fee_charged":${1}"${2}",`,
	// },
	// {regexp.MustCompile(
	// 	// max_fee is a string since aurora 1.3.0
	// 	`"max_fee":( ?)([\d]+),`),
	// 	`"max_fee":${1}"${2}",`,
	// },
	// // Removes trailing SSE data, fixed in aurora 1.7.0
	// {regexp.MustCompile(
	// 	`\nretry:.*\nevent:.*\ndata:.*\n`),
	// 	``,
	// },
	// // Removes clawback, fixed in aurora 2.1.0
	// {regexp.MustCompile(
	// 	`,\s*"auth_clawback_enabled":\s*false`),
	// 	``,
	// },
}

var newAccountDetailsPathWithLastestLedger = regexp.MustCompile(`^/accounts/[A-Z0-9]+/(transactions|operations|payments|effects|trades)/?`)

type Response struct {
	Domain string
	Path   string
	Stream bool

	StatusCode   int
	LatestLedger string
	Body         string
	// NormalizedBody is body without parts that identify a single
	// server (ex. domain) and fields known to be different between
	// instances (ex. `result_meta_xdr`).
	NormalizedBody string
}

func NewResponse(domain, path string, stream bool) *Response {
	response := &Response{
		Domain: domain,
		Path:   path,
		Stream: stream,
	}

	req, err := http.NewRequest("GET", domain+path, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	if stream {
		req.Header.Add("Accept", "text/event-stream")
		// Since requests are made in separate go routines we can
		// set timeout to one minute.
		client.Timeout = time.Minute
	}

	resp, err := client.Do(req)
	if err != nil {
		response.Body = err.Error()
		response.NormalizedBody = err.Error()
		return response
	}

	response.StatusCode = resp.StatusCode

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNotFound &&
		resp.StatusCode != http.StatusNotAcceptable &&
		resp.StatusCode != http.StatusBadRequest &&
		resp.StatusCode != http.StatusGatewayTimeout &&
		resp.StatusCode != http.StatusGone &&
		resp.StatusCode != http.StatusServiceUnavailable {
		panic(resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	// We ignore the error below to timeout streaming requests.
	// net/http: request canceled (Client.Timeout exceeded while reading body)
	if err != nil && !stream {
		response.Body = err.Error()
		response.NormalizedBody = err.Error()
		return response
	}

	if string(body) == "" {
		response.Body = fmt.Sprintf("Empty body [%d]", rand.Uint64())
	}

	response.LatestLedger = resp.Header.Get("Latest-Ledger")
	response.Body = string(body)

	normalizedBody := response.Body
	// `result_meta_xdr` can differ between core instances (confirmed this with core team)
	normalizedBody = findResultMetaXDR.ReplaceAllString(normalizedBody, "")
	// Remove Aurora URL from the _links
	normalizedBody = strings.Replace(normalizedBody, domain, "", -1)

	for _, reg := range removeRegexps {
		normalizedBody = reg.ReplaceAllString(normalizedBody, "")
	}

	for _, reg := range replaceRegexps {
		normalizedBody = reg.regexp.ReplaceAllString(normalizedBody, reg.repl)
	}

	// 1.1.0 - skip Latest-Ledger header in newly incorporated endpoints
	if !(newAccountDetailsPathWithLastestLedger.Match([]byte(path)) ||
		strings.HasPrefix(path, "/ledgers") ||
		strings.HasPrefix(path, "/transactions") ||
		strings.HasPrefix(path, "/operations") ||
		strings.HasPrefix(path, "/payments") ||
		strings.HasPrefix(path, "/effects") ||
		strings.HasPrefix(path, "/transactions") ||
		strings.Contains(path, "/trade")) {
		normalizedBody = fmt.Sprintf("Latest-Ledger: %s\n%s", resp.Header.Get("Latest-Ledger"), normalizedBody)
	}
	response.NormalizedBody = normalizedBody
	return response
}

func (r *Response) Equal(other *Response) bool {
	return r.NormalizedBody == other.NormalizedBody
}

func (r *Response) Size() int {
	return len(r.Body)
}

func (r *Response) SaveDiff(outputDir string, other *Response) {
	if r.Path != other.Path {
		panic("Paths are different")
	}

	fileName := pathToFileName(r.Path, r.Stream)

	if len(fileName) > fileLengthLimit {
		fileName = fileName[0:fileLengthLimit]
	}

	fileA := fmt.Sprintf("%s/%s.old", outputDir, fileName)
	fileB := fmt.Sprintf("%s/%s.new", outputDir, fileName)
	fileDiff := fmt.Sprintf("%s/%s.diff", outputDir, fileName)

	// We compare normalized body to see actual differences in the diff instead
	// of a lot of domain diffs.
	err := ioutil.WriteFile(fileA, []byte(r.Domain+" "+r.Path+"\n\n"+r.NormalizedBody), 0744)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileB, []byte(other.Domain+" "+other.Path+"\n\n"+other.NormalizedBody), 0744)
	if err != nil {
		panic(err)
	}

	// Ignore `err`, user will generate diff manually
	out, _ := exec.Command("diff", fileA, fileB).Output()

	if len(out) != 0 {
		err = ioutil.WriteFile(fileDiff, out, 0744)
		if err != nil {
			panic(err)
		}
	}
}

// GetPaths finds all URLs in the response body and returns paths
// (without domain).
func (r *Response) GetPaths() []string {
	// escapedDomain := strings.Replace(r.Domain, `\`, `\\`, -1)
	var linksRegexp = regexp.MustCompile(`"` + r.Domain + `(.*?)["{]`)
	found := linksRegexp.FindAllSubmatch([]byte(r.Body), -1)
	links := make([]string, 0, len(found))

	for _, link := range found {
		l := strings.Replace(string(link[1]), "\\u0026", "&", -1)
		links = append(links, l)
	}

	return links
}

func pathToFileName(path string, stream bool) string {
	if stream {
		path = "stream_" + path
	}
	path = strings.Replace(path, "/", "_", -1)
	path = strings.Replace(path, "?", "_", -1)
	path = strings.Replace(path, "&", "_", -1)
	path = strings.Replace(path, "=", "_", -1)
	return path
}

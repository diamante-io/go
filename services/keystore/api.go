package keystore

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
	"github.com/diamcircle/go/support/errors"
	"github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/support/render/health"
	"github.com/diamcircle/go/support/render/httpjson"
	"github.com/diamcircle/go/support/render/problem"
)

func init() {
	// register errors
	problem.RegisterError(httpjson.ErrBadRequest, probInvalidRequest)
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)

	// register service host as an empty string
	problem.RegisterHost("")
}

func (s *Service) wrapMiddleware(handler http.Handler) http.Handler {
	handler = authHandler(handler, s.authenticator)
	handler = recoverHandler(handler)
	handler = corsHandler(handler)
	return handler
}

func ServeMux(s *Service) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/keys", s.wrapMiddleware(s.keysHTTPMethodHandler()))
	mux.Handle("/health", s.wrapMiddleware(health.PassHandler{}))
	return mux
}

func (s *Service) keysHTTPMethodHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			jsonHandler(s.getKeys).ServeHTTP(rw, req)

		case http.MethodPut:
			jsonHandler(s.putKeys).ServeHTTP(rw, req)

		case http.MethodDelete:
			jsonHandler(s.deleteKeys).ServeHTTP(rw, req)

		default:
			problem.Render(req.Context(), rw, probMethodNotAllowed)
		}
	})
}

type authResponse struct {
	UserID string `json:"userID"`
}

var forwardHeaders = map[string]struct{}{
	"authorization": struct{}{},
	"cookie":        struct{}{},
}

func authHandler(next http.Handler, authenticator *Authenticator) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if authenticator == nil {
			// to facilitate API testing
			next.ServeHTTP(rw, req.WithContext(withUserID(req.Context(), "test-user")))
			return
		}

		var (
			proxyReq *http.Request
			err      error
			clientIP string
		)
		ctx := req.Context()
		// set a 5-second timeout
		client := http.Client{Timeout: time.Duration(5 * time.Second)}

		switch authenticator.APIType {
		case REST:
			proxyReq, err = http.NewRequest("GET", authenticator.URL, nil)
			if err != nil {
				problem.Render(ctx, rw, errors.Wrap(err, "creating the auth proxy request"))
				return
			}

		case GraphQL:
			// to be implemented later
		default:
			problem.Render(ctx, rw, probNotAuthorized)
			return
		}

		proxyReq.Header = make(http.Header)
		for k, v := range req.Header {
			// http headers are case-insensitive
			// https://www.ietf.org/rfc/rfc2616.txt
			if _, ok := forwardHeaders[strings.ToLower(k)]; ok {
				proxyReq.Header[k] = v
			}
		}

		if clientIP, _, err = net.SplitHostPort(req.RemoteAddr); err == nil {
			proxyReq.Header.Set("X-Forwarded-For", clientIP)
		}
		proxyReq.Header.Set("Accept-Encoding", "identity")

		resp, err := client.Do(proxyReq)
		if err != nil {
			problem.Render(ctx, rw, errors.Wrap(err, "sending the auth proxy request"))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			problem.Render(ctx, rw, probNotAuthorized)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			problem.Render(ctx, rw, errors.Wrap(err, "reading the auth response"))
			return
		}

		var authResp authResponse
		err = json.Unmarshal(body, &authResp)
		if err != nil {
			log.Ctx(ctx).Infof("Response body as a plain string: %s\n. Response body as a hex dump string: %s\n", string(body), hex.Dump(body))
			problem.Render(ctx, rw, errors.Wrap(err, "unmarshaling the auth response"))
			return
		}
		if authResp.UserID == "" {
			problem.Render(ctx, rw, probNotAuthorized)
			return
		}

		next.ServeHTTP(rw, req.WithContext(withUserID(ctx, authResp.UserID)))
	})
}

func jsonHandler(f interface{}) http.Handler {
	h, err := httpjson.ReqBodyHandler(f, httpjson.JSON)
	if err != nil {
		panic(err)
	}
	return h
}

func recoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("panic: %v", r)
			}
			if errors.Cause(err) == http.ErrAbortHandler {
				panic(err)
			}

			ctx := req.Context()
			log.Ctx(ctx).WithStack(err).Error(err)
			problem.Render(ctx, rw, err)
		}()

		next.ServeHTTP(rw, req)
	})
}

func corsHandler(next http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "POST", "PATCH", "DELETE", "HEAD", "OPTIONS"},
	})
	return cors.Handler(next)
}

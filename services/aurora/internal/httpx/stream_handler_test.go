package httpx

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/go-chi/chi"

	"github.com/diamcircle/go/services/aurora/internal/actions"
	auroraContext "github.com/diamcircle/go/services/aurora/internal/context"
	"github.com/diamcircle/go/services/aurora/internal/ledger"
	"github.com/diamcircle/go/services/aurora/internal/render/sse"
	"github.com/diamcircle/go/support/db"
	"github.com/diamcircle/go/support/render/hal"
)

type testingFactory struct {
	ledgerSource ledger.Source
}

func (f *testingFactory) Get() ledger.Source {
	return f.ledgerSource
}

// StreamTest utility struct to wrap SSE related tests.
type StreamTest struct {
	ledgerSource  *ledger.TestingSource
	cancel        context.CancelFunc
	wg            *sync.WaitGroup
	w             *httptest.ResponseRecorder
	checkResponse func(w *httptest.ResponseRecorder)
	ctx           context.Context
}

func newStreamTest(
	handler http.HandlerFunc,
	ledgerSource *ledger.TestingSource,
	request *http.Request,
	checkResponse func(w *httptest.ResponseRecorder),
) *StreamTest {
	s := &StreamTest{
		ledgerSource:  ledgerSource,
		w:             httptest.NewRecorder(),
		checkResponse: checkResponse,
		wg:            &sync.WaitGroup{},
	}
	s.ctx, s.cancel = context.WithCancel(request.Context())

	s.wg.Add(1)
	go func() {
		handler(s.w, request.WithContext(s.ctx))
		s.wg.Done()
		s.cancel()
	}()

	return s
}

// NewStreamablePageTest tests the SSE functionality of a pageAction
func NewStreamablePageTest(
	action *testPageAction,
	currentLedger uint32,
	request *http.Request,
	checkResponse func(w *httptest.ResponseRecorder),
) *StreamTest {
	ledgerSource := ledger.NewTestingSource(currentLedger)
	action.ledgerSource = ledgerSource
	streamHandler := sse.StreamHandler{LedgerSourceFactory: &testingFactory{ledgerSource}}
	handler := streamableStatePageHandler(&ledger.State{}, action, streamHandler)

	return newStreamTest(
		handler.renderStream,
		ledgerSource,
		request,
		checkResponse,
	)
}

// NewStreamableObjectTest tests the SSE functionality of a streamableObjectAction
func NewStreamableObjectTest(
	action *testObjectAction,
	currentLedger uint32,
	request *http.Request,
	limit int,
	checkResponse func(w *httptest.ResponseRecorder),
) *StreamTest {
	ledgerSource := ledger.NewTestingSource(currentLedger)
	action.ledgerSource = ledgerSource
	streamHandler := sse.StreamHandler{LedgerSourceFactory: &testingFactory{ledgerSource}}
	handler := streamableObjectActionHandler{action: action, limit: limit, streamHandler: streamHandler}

	return newStreamTest(
		handler.renderStream,
		ledgerSource,
		request,
		checkResponse,
	)
}

// AddLedger pushes a new ledger to the stream handler. AddLedger() will block until
// the new ledger has been read by the stream handler
func (s *StreamTest) AddLedger(sequence uint32) {
	s.ledgerSource.AddLedger(sequence)
}

// Stop ends the stream request and checks the response
func (s *StreamTest) Stop() {
	s.cancel()
	s.wg.Wait()
	s.checkResponse(s.w)
}

// Wait blocks testing until the stream test has finished running and checks the response
func (s *StreamTest) Wait() {
	s.wg.Wait()
	s.checkResponse(s.w)
}

type testPage struct {
	Value       string `json:"value"`
	pagingToken int
}

func (p testPage) PagingToken() string {
	return fmt.Sprintf("%v", p.pagingToken)
}

type testPageAction struct {
	objects      map[uint32][]string
	ledgerSource ledger.Source
}

func (action *testPageAction) GetResourcePage(
	w actions.HeaderWriter,
	r *http.Request,
) ([]hal.Pageable, error) {
	objects, ok := action.objects[action.ledgerSource.CurrentLedger()]
	if !ok {
		return nil, fmt.Errorf("unexpected ledger")
	}

	cursor := r.Header.Get("Last-Event-ID")
	if cursor == "" {
		cursor = r.URL.Query().Get("cursor")
	}
	if cursor == "" {
		cursor = "0"
	}
	parsedCursor, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, err
	}

	limit := len(objects)
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return nil, err
		}
	}

	if parsedCursor < 0 {
		return nil, fmt.Errorf("cursor cannot be negative")
	}

	if parsedCursor >= len(objects) {
		return []hal.Pageable{}, nil
	}

	response := []hal.Pageable{}
	for i, object := range objects[parsedCursor:] {
		if len(response) >= limit {
			break
		}

		response = append(response, testPage{Value: object, pagingToken: parsedCursor + i + 1})
	}

	return response, nil
}

func streamRequest(t *testing.T, queryParams string) *http.Request {
	request, err := http.NewRequest("GET", "http://localhost?"+queryParams, nil)
	if err != nil {
		t.Fatalf("could not construct request: %v", err)
	}
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, chi.NewRouteContext())
	request = request.WithContext(ctx)

	return request
}

func unmarashalPage(jsonString string) (string, error) {
	var page testPage
	err := json.Unmarshal([]byte(jsonString), &page)
	return page.Value, err
}

func expectResponse(
	t *testing.T,
	unmarshal func(string) (string, error),
	expectedResponse []string,
) func(*httptest.ResponseRecorder) {
	return func(w *httptest.ResponseRecorder) {
		var response []string
		for _, line := range strings.Split(w.Body.String(), "\n") {
			if line == "data: \"hello\"" || line == "data: \"byebye\"" {
				continue
			}

			if strings.HasPrefix(line, "data: ") {
				jsonString := line[len("data: "):]
				value, err := unmarshal(jsonString)
				if err != nil {
					t.Fatalf("could not parse json %v", err)
				}
				response = append(response, value)
			}
		}

		if len(expectedResponse) != len(response) {
			t.Fatalf("expected %v but got %v", expectedResponse, response)
		}

		for i, entry := range expectedResponse {
			if entry != response[i] {
				t.Fatalf("expected %v but got %v", expectedResponse, response)
			}
		}
	}
}

func TestPageStream(t *testing.T) {
	t.Run("without offset", func(t *testing.T) {
		request := streamRequest(t, "")
		action := &testPageAction{
			objects: map[uint32][]string{
				3: {"a", "b", "c"},
				4: {"a", "b", "c", "d", "e"},
				6: {"a", "b", "c", "d", "e", "f"},
				7: {"a", "b", "c", "d", "e", "f"},
			},
		}
		st := NewStreamablePageTest(
			action,
			3,
			request,
			expectResponse(t, unmarashalPage, []string{"a", "b", "c", "d", "e", "f"}),
		)

		st.AddLedger(4)
		st.AddLedger(6)
		st.AddLedger(7)

		st.Stop()
	})

	t.Run("with offset", func(t *testing.T) {
		request := streamRequest(t, "cursor=1")
		action := &testPageAction{
			objects: map[uint32][]string{
				3: {"a", "b", "c"},
				4: {"a", "b", "c", "d", "e"},
				6: {"a", "b", "c", "d", "e", "f"},
				7: {"a", "b", "c", "d", "e", "f"},
			},
		}
		st := NewStreamablePageTest(
			action,
			3,
			request,
			expectResponse(t, unmarashalPage, []string{"b", "c", "d", "e", "f"}),
		)

		st.AddLedger(4)
		st.AddLedger(6)
		st.AddLedger(7)

		st.Stop()
	})

	t.Run("with limit", func(t *testing.T) {
		request := streamRequest(t, "limit=2")
		action := &testPageAction{
			objects: map[uint32][]string{
				3: {"a", "b", "c"},
			},
		}
		st := NewStreamablePageTest(
			action,
			3,
			request,
			expectResponse(t, unmarashalPage, []string{"a", "b"}),
		)

		st.Wait()
	})

	t.Run("with limit and offset", func(t *testing.T) {
		request := streamRequest(t, "limit=2&cursor=1")
		action := &testPageAction{
			objects: map[uint32][]string{
				3: {"a", "b", "c", "d", "e"},
			},
		}
		st := NewStreamablePageTest(
			action,
			3,
			request,
			expectResponse(t, unmarashalPage, []string{"b", "c"}),
		)

		st.Wait()
	})

	t.Run("reach limit after multiple iterations", func(t *testing.T) {
		request := streamRequest(t, "limit=3&cursor=1")
		action := &testPageAction{
			objects: map[uint32][]string{
				3: {"a"},
				4: {"a", "b"},
				5: {"a", "b", "c", "d", "e", "f", "g"},
			},
		}
		st := NewStreamablePageTest(
			action,
			3,
			request,
			expectResponse(t, unmarashalPage, []string{"b", "c", "d"}),
		)

		st.AddLedger(4)
		st.AddLedger(5)

		st.Wait()
	})
}

type stringObject string

func (s stringObject) Equals(other actions.StreamableObjectResponse) bool {
	otherString, ok := other.(stringObject)
	if !ok {
		return false
	}
	return s == otherString
}

func unmarashalString(jsonString string) (string, error) {
	var object stringObject
	err := json.Unmarshal([]byte(jsonString), &object)
	return string(object), err
}

type testObjectAction struct {
	objects      map[uint32]stringObject
	ledgerSource ledger.Source
}

func (action *testObjectAction) GetResource(
	w actions.HeaderWriter,
	r *http.Request,
) (actions.StreamableObjectResponse, error) {
	ledger := action.ledgerSource.CurrentLedger()
	object, ok := action.objects[ledger]
	if !ok {
		return nil, fmt.Errorf("unexpected ledger: %v", ledger)
	}

	return object, nil
}

func TestObjectStream(t *testing.T) {
	t.Run("without interior duplicates", func(t *testing.T) {
		request := streamRequest(t, "")
		action := &testObjectAction{
			objects: map[uint32]stringObject{
				3: "a",
				4: "b",
				5: "c",
				6: "c",
			},
		}

		st := NewStreamableObjectTest(
			action,
			3,
			request,
			10,
			expectResponse(t, unmarashalString, []string{"a", "b", "c"}),
		)

		st.AddLedger(4)
		st.AddLedger(5)
		st.AddLedger(6)
		st.Stop()
	})

	t.Run("with interior duplicates", func(t *testing.T) {
		request := streamRequest(t, "")
		action := &testObjectAction{
			objects: map[uint32]stringObject{
				3: "a",
				4: "b",
				5: "b",
				6: "c",
				7: "c",
			},
		}

		st := NewStreamableObjectTest(
			action,
			3,
			request,
			10,
			expectResponse(t, unmarashalString, []string{"a", "b", "c"}),
		)

		st.AddLedger(4)
		st.AddLedger(5)
		st.AddLedger(6)
		st.AddLedger(7)

		st.Stop()
	})

	t.Run("limit reached", func(t *testing.T) {
		request := streamRequest(t, "")
		action := &testObjectAction{
			objects: map[uint32]stringObject{
				1: "a",
				2: "b",
				3: "b",
				4: "c",
				5: "d",
			},
		}

		st := NewStreamableObjectTest(
			action,
			1,
			request,
			4,
			expectResponse(
				t,
				unmarashalString,
				[]string{
					"a", "b", "c", "d",
				},
			),
		)

		st.AddLedger(2)
		st.AddLedger(3)
		st.AddLedger(4)
		st.AddLedger(5)

		st.Wait()
	})
}

func TestRepeatableReadStream(t *testing.T) {
	t.Run("page stream creates repeatable read tx", func(t *testing.T) {
		action := &testPageAction{
			objects: map[uint32][]string{
				3: {"a"},
				4: {"a", "b"},
			},
		}

		session := &db.MockSession{}
		session.On("BeginTx", &sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
			ReadOnly:  true,
		}).Return(nil).Once()
		session.On("Rollback").Return(nil).Once()

		session.On("BeginTx", &sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
			ReadOnly:  true,
		}).Return(nil).Once()
		session.On("Rollback").Return(nil).Once()

		request := streamRequest(t, "limit=2")
		request = request.WithContext(context.WithValue(
			request.Context(),
			&auroraContext.SessionContextKey,
			session,
		))

		st := NewStreamablePageTest(
			action,
			3,
			request,
			expectResponse(t, unmarashalPage, []string{"a", "b"}),
		)
		st.AddLedger(4)
		st.Wait()
		session.AssertExpectations(t)
	})

	t.Run("object stream creates repeatable read tx", func(t *testing.T) {
		action := &testObjectAction{
			objects: map[uint32]stringObject{
				3: "a",
				4: "b",
			},
		}

		session := &db.MockSession{}
		session.On("BeginTx", &sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
			ReadOnly:  true,
		}).Return(nil).Once()
		session.On("Rollback").Return(nil).Once()

		session.On("BeginTx", &sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
			ReadOnly:  true,
		}).Return(nil).Once()
		session.On("Rollback").Return(nil).Once()

		request := streamRequest(t, "")
		request = request.WithContext(context.WithValue(
			request.Context(),
			&auroraContext.SessionContextKey,
			session,
		))

		st := NewStreamableObjectTest(
			action,
			3,
			request,
			2,
			expectResponse(t, unmarashalString, []string{"a", "b"}),
		)
		st.AddLedger(4)
		st.Wait()
		session.AssertExpectations(t)
	})
}

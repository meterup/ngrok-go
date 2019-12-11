package ngrok

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/kevinburke/rest"
)

// Use this when initializing a new test - first hit the real ngrok API
// with DEBUG_HTTP_TRAFFIC=true in your environment, save the result in
// responses_test, and then replay that response over and over again.
var envClient = New(BaseURL, os.Getenv("NGROK_API_TOKEN"))
var _ = envClient

// this is all taken from kevinburke/twilio-go, MIT licensed

type Server struct {
	s *httptest.Server
	// copied from httptest.Server
	URL string
	// URLs of incoming requests, in order
	URLs []*url.URL
	mu   sync.Mutex
}

func newServer(response []byte, code int) *Server {
	serv := &Server{URLs: make([]*url.URL, 0)}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serv.mu.Lock()
		serv.URLs = append(serv.URLs, r.URL)
		serv.mu.Unlock()
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(code)
		if _, err := w.Write(response); err != nil {
			panic(err)
		}
	}))
	serv.s = s
	serv.URL = s.URL
	return serv
}

var badContentType = `
{
    "details": null,
    "error_code": "ERR_NGROK_210",
    "msg": "The content type you specified 'application/x-www-form-urlencoded; blah' is not supported by the API. Please check your API client implementation and see the list of supported content types: https://ngrok.com/docs/ngrok-link#service-api-content-type",
    "status_code": 415
}
`

func TestErrorParse(t *testing.T) {
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(415)
		io.WriteString(w, badContentType)
	}))
	defer errorServer.Close()
	client := New(errorServer.URL, os.Getenv("NGROK_API_TOKEN"))
	_, err := client.Creds.Create(context.Background(), url.Values{})
	if err == nil {
		t.Fatalf("got nil error, want non-nil")
	}
	rerr, ok := err.(*rest.Error)
	if !ok {
		t.Fatalf("want a rest.Error, got %v", err)
	}
	if rerr.ID != "ERR_NGROK_210" {
		t.Errorf("bad ID: got %q want ERR_NGROK_210", rerr.ID)
	}
}

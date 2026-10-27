package corellium

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockAPIServer is a local stand-in for the Corellium API used by unit
// tests. It serves TLS (self-signed) because the generated API client
// hardcodes the "https" scheme and only the host is overridable via the
// provider's "host" configuration attribute. The API client falls back to
// http.DefaultClient when no HTTPClient is configured, so newMockAPIServer
// swaps http.DefaultClient for one that trusts the test certificate,
// avoiding any production code changes.
//
// All unit tests must configure the provider with the same literal token
// value (e.g. token = ""), matching the existing acceptance test
// convention. corellium/pkg/api.SetAccessToken is a process-wide
// set-once singleton, so the first test to configure the provider in a
// given `go test` binary run freezes the token for every other test in
// that binary.
type mockAPIServer struct {
	*httptest.Server
	mux *http.ServeMux
}

// newMockAPIServer starts a local TLS server pre-registered with a
// GET /api/v3/ready route (used by the provider's Configure readiness
// check and by the corellium_v1ready data source) and returns it ready
// for additional routes to be registered via handle. All routes must
// include the "/api" prefix baked into the generated client's server
// URL template.
func newMockAPIServer(t *testing.T) *mockAPIServer {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v3/ready", jsonHandler(http.StatusOK, map[string]string{"status": "ok"}))

	server := httptest.NewTLSServer(mux)
	t.Cleanup(server.Close)

	// The provider runs in-process and builds its API client with
	// http.DefaultClient, so point that at a client trusting this server's
	// self-signed certificate for the duration of the test.
	origClient := http.DefaultClient
	http.DefaultClient = server.Client()
	t.Cleanup(func() { http.DefaultClient = origClient })

	return &mockAPIServer{Server: server, mux: mux}
}

// handle registers an additional route on the mock server, e.g.
// handle("POST /api/v3/users", jsonHandler(http.StatusOK, someBody)).
func (m *mockAPIServer) handle(pattern string, fn http.HandlerFunc) {
	m.mux.HandleFunc(pattern, fn)
}

// hostAttr returns the server's address without the scheme, suitable for
// interpolating into the provider's "host" configuration attribute.
func (m *mockAPIServer) hostAttr() string {
	return strings.TrimPrefix(m.URL, "https://")
}

// jsonHandler returns an http.HandlerFunc that writes the given status
// code and JSON-encodes body as the response.
func jsonHandler(status int, body any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}
}

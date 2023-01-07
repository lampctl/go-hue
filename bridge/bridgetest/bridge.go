package bridgetest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"

	"github.com/lampctl/go-hue/bridge"
)

// Bridge represents a "fake" Hue bridge that simulates the functions of an
// actual bridge as closely as possible.
type Bridge struct {

	// URL represents the URL of the server.
	URL string

	mutex         sync.Mutex
	mux           *http.ServeMux
	server        *httptest.Server
	resources     map[string]*bridge.Resource
	buttonPressed bool
}

// New creates a new bridge.
func New() (*Bridge, error) {
	s := httptest.NewTLSServer(nil)
	u, err := url.Parse(s.URL)
	if err != nil {
		return nil, err
	}
	b := &Bridge{
		URL:       u.Host,
		mux:       http.NewServeMux(),
		server:    s,
		resources: make(map[string]*bridge.Resource),
	}
	b.mux.HandleFunc("/api", b.handleApi)
	b.mux.HandleFunc("/clip/v2/resource", b.requireAuth(b.handleResource))
	return b, nil
}

// AddResource adds the provided resource to the bridge.
func (b *Bridge) AddResource(r *bridge.Resource) {
	defer b.mutex.Unlock()
	b.mutex.Lock()
	b.resources[r.ID] = r
}

// PushButton simulates a user pressing the button on the bridge.
func (b *Bridge) PushButton() {
	defer b.mutex.Unlock()
	b.mutex.Lock()
	b.buttonPressed = true
}

// Close shuts down the bridge.
func (b *Bridge) Close() {
	b.server.Close()
}

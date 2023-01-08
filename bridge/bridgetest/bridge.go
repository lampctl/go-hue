package bridgetest

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"

	"github.com/lampctl/go-hue/bridge"
)

var errResourceNotFound = errors.New("resource not found")

// Bridge represents a "fake" Hue bridge that simulates the functions of an
// actual bridge as closely as possible.
type Bridge struct {

	// URL represents the URL of the server.
	URL string

	mutex         sync.Mutex
	mux           *http.ServeMux
	server        *httptest.Server
	resources     map[string]*bridge.Resource
	apiRequested  bool
	buttonPressed bool
}

// New creates a new bridge.
func New() (*Bridge, error) {
	var (
		m = http.NewServeMux()
		s = httptest.NewTLSServer(m)
	)
	u, err := url.Parse(s.URL)
	if err != nil {
		return nil, err
	}
	b := &Bridge{
		URL:       u.Host,
		mux:       m,
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
	b.mux.HandleFunc(
		fmt.Sprintf("/clip/v2/resource/%s/%s", r.Type, r.ID),
		b.requireAuth(b.handleResourceByID(r.ID)),
	)
}

// GetResource attempts to retrieve the specified resource by its ID.
func (b *Bridge) GetResource(id string) (*bridge.Resource, error) {
	defer b.mutex.Unlock()
	b.mutex.Lock()
	r, ok := b.resources[id]
	if !ok {
		return nil, errResourceNotFound
	}
	return r, nil
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

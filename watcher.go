package hue

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/lampctl/go-hue/bridge"
	"github.com/lampctl/go-sse"
)

// Watcher receives events from the bridge and maintains a list of resources
// and their current states.
type Watcher struct {
	mutex      sync.Mutex
	client     *Client
	sseClient  *sse.Client
	resources  map[string]*bridge.Resource
	closedChan <-chan any
}

func (w *Watcher) loadResources() error {
	resources, err := w.client.Resources()
	if err != nil {
		return err
	}
	for _, r := range resources {
		w.resources[r.ID] = r
	}
	return nil
}

func (w *Watcher) lifecycleLoop(closedChan chan<- any) {
	defer close(closedChan)
	for {
		e, ok := <-w.sseClient.Events
		if !ok {
			return
		}
		var (
			stream = []*bridge.EventStream{}
			reader = strings.NewReader(e.Data)
		)
		if err := json.NewDecoder(reader).Decode(&stream); err != nil {
			continue
		}
		func() {
			defer w.mutex.Unlock()
			w.mutex.Lock()
			for _, v := range stream {
				for _, e := range v.Data {
					r, ok := w.resources[e.ID]
					if ok {
						r.CopyFrom(e)
					}
				}
			}
		}()
	}
}

// NewWatcher creates and initializes a new watcher using the provided client.
func NewWatcher(c *Client) (*Watcher, error) {
	r, err := c.newRequest(http.MethodGet, "/eventstream/clip/v2", nil)
	if err != nil {
		return nil, err
	}
	var (
		closedChan = make(chan any)
		w          = &Watcher{
			client:     c,
			sseClient:  sse.NewClient(r, c.Client),
			resources:  make(map[string]*bridge.Resource),
			closedChan: closedChan,
		}
	)
	if err := w.loadResources(); err != nil {
		return nil, err
	}
	go w.lifecycleLoop(closedChan)
	return w, nil
}

// Resources returns a slice of resources and their current values.
func (w *Watcher) Resources() []*bridge.Resource {
	defer w.mutex.Unlock()
	w.mutex.Lock()
	resources := []*bridge.Resource{}
	for _, r := range w.resources {
		rCopy := &bridge.Resource{}
		rCopy.CopyFrom(r)
		resources = append(resources, rCopy)
	}
	return resources
}

// Close shuts down the watcher.
func (w *Watcher) Close() {
	w.sseClient.Close()
	<-w.closedChan
}

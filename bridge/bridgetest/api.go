package bridgetest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lampctl/go-hue/bridge"
)

const (
	// Username is the application key for the bridge.
	Username = "username"

	buttonNotPressedError = "button has not been pressed"
	invalidUsernameError  = "invalid username supplied"
)

func (b *Bridge) writeJson(w http.ResponseWriter, statusCode int, v any) {
	d, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(d)))
	w.WriteHeader(statusCode)
	w.Write(d)
}

func (b *Bridge) writeError(w http.ResponseWriter, statusCode int, desc string) {
	b.writeJson(w, statusCode, &bridge.Response{
		Errors: []*bridge.Error{{Description: desc}},
	})
}

func (b *Bridge) writeHTTPError(w http.ResponseWriter, statusCode int) {
	b.writeError(w, statusCode, http.StatusText(statusCode))
}

func (b *Bridge) writeData(w http.ResponseWriter, v any) {
	d, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	b.writeJson(w, http.StatusOK, &bridge.Response{Data: json.RawMessage(d)})
}

func (b *Bridge) requireAuth(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("hue-application-key") != Username {
			http.Error(w, invalidUsernameError, http.StatusForbidden)
			return
		}
		fn(w, r)
	})
}

func (b *Bridge) handleApi(w http.ResponseWriter, r *http.Request) {
	defer b.mutex.Unlock()
	b.mutex.Lock()
	response := &bridge.RegistrationResponse{}
	if b.apiRequested && b.buttonPressed {
		response.Success = &bridge.RegistrationResponseSuccess{
			Username: Username,
		}
	} else {
		b.apiRequested = true
		response.Error = &bridge.RegistrationResponseError{
			Description: buttonNotPressedError,
		}
	}
	b.writeJson(w, http.StatusOK, []*bridge.RegistrationResponse{response})
}

func (b *Bridge) handleResource(w http.ResponseWriter, r *http.Request) {
	defer b.mutex.Unlock()
	b.mutex.Lock()
	resources := []*bridge.Resource{}
	for _, r := range b.resources {
		resources = append(resources, r)
	}
	b.writeData(w, resources)
}

func (b *Bridge) handleResourceByID(id string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer b.mutex.Unlock()
		b.mutex.Lock()
		if r.Method != http.MethodPut {
			b.writeHTTPError(w, http.StatusMethodNotAllowed)
			return
		}
		src := &bridge.Resource{}
		if err := json.NewDecoder(r.Body).Decode(src); err != nil {
			b.writeHTTPError(w, http.StatusBadRequest)
			return
		}
		dest, ok := b.resources[id]
		if !ok {
			b.writeHTTPError(w, http.StatusNotFound)
			return
		}
		dest.CopyFrom(src)
	})
}

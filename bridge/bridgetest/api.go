package bridgetest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lampctl/go-hue/bridge"
)

const (
	username              = "username"
	buttonNotPressedError = "button has not been pressed"
	invalidUsernameError  = "invalid username supplied"
)

func (b *Bridge) writeJson(w http.ResponseWriter, v any) {
	d, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(d)))
	w.Write(d)
}

func (b *Bridge) requireAuth(fn http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Response.Header.Get("hue-application-key") != username {
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
	if b.buttonPressed {
		response.Success = &bridge.RegistrationResponseSuccess{
			Username: username,
		}
	} else {
		response.Error = &bridge.RegistrationResponseError{
			Description: buttonNotPressedError,
		}
	}
	b.writeJson(w, []*bridge.RegistrationResponse{response})
}

func (b *Bridge) handleResource(w http.ResponseWriter, r *http.Request) {
	defer b.mutex.Unlock()
	b.mutex.Lock()
	resources := []*bridge.Resource{}
	for _, r := range b.resources {
		resources = append(resources, r)
	}
	b.writeJson(w, resources)
}

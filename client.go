package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

var (
	errInvalidResponse = errors.New("invalid response received from bridge")

	// TODO: there must be a better way to do this

	// DefaultClient provides an HTTP client implementation that uses the
	// InsecureSkipVerify flag to bypass errors caused by the bridge's
	// certificate - it does not contain a SAN for its IP address.
	DefaultClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

type bridgeResponse struct {
	Errors []any           `json:"errors"`
	Data   json.RawMessage `json:"data"`
}

type bridgeRegistrationRequest struct {
	DeviceType        string `json:"devicetype"`
	GenerateClientKey bool   `json:"generateclientkey"`
}

type bridgeRegistrationResponse struct {
	Error *struct {
		Description string `json:"description"`
	} `json:"error"`
	Success *struct {
		Username string `json:"username"`
	} `json:"success"`
}

// Client represents a connection to a Hue bridge. Hostname must be set to the
// address of the bridge and Username needs to be set to the application key to
// use authenticated methods. Use the Register() method to obtain an
// application key.
type Client struct {

	// Host represents the host address of the bridge.
	Host string

	// Username represents the application key used for authentication.
	Username string

	// Client represents the http.Client used for the connections. If this
	// value is nil, DefaultClient is used. Note that the client must be
	// willing to accept the bridge's certificate, which is not signed by a CA
	// recognized in most operating systems.
	Client *http.Client
}

func (c *Client) newRequest(method, path string, v any) (*http.Request, error) {
	var reader io.Reader
	if v != nil {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewBuffer(b)
	}
	u := &url.URL{
		Scheme: "https",
		Host:   c.Host,
		Path:   path,
	}
	r, err := http.NewRequest(method, u.String(), reader)
	if err != nil {
		return nil, err
	}
	if len(c.Username) != 0 {
		r.Header.Add("hue-application-key", c.Username)
	}
	return r, nil
}

func (c *Client) doRequest(req *http.Request) (*bridgeResponse, error) {
	client := c.Client
	if client == nil {
		client = DefaultClient
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	response := &bridgeResponse{}
	if err := json.NewDecoder(r.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, err
}

func (c *Client) do(method, path string, v any) (*bridgeResponse, error) {
	r, err := c.newRequest(method, path, v)
	if err != nil {
		return nil, err
	}
	return c.doRequest(r)
}

// Register attempts to register with the bridge and obtain an application key.
// This will likely fail the first time with an error instructing the user to
// push the physical button on their bridge - an identical call to this method
// a second time should then succeed. Upon success, the application key will be
// stored in the Username field of the Client.
func (c *Client) Register(appName string) error {
	r, err := c.newRequest(http.MethodPost, "/api", &bridgeRegistrationRequest{
		DeviceType:        appName,
		GenerateClientKey: true,
	})
	if err != nil {
		return err
	}
	defer r.Body.Close()
	responses := []*bridgeRegistrationResponse{}
	if err := json.NewDecoder(r.Body).Decode(&responses); err != nil {
		return err
	}
	if len(responses) < 1 {
		return errInvalidResponse
	}
	if responses[0].Error != nil {
		return errors.New(responses[0].Error.Description)
	}
	if responses[0].Success == nil {
		return errInvalidResponse
	}
	c.Username = responses[0].Success.Username
	return nil
}

// Resources retrieves all of the resources on the bridge.
func (c *Client) Resources() ([]*Resource, error) {
	r, err := c.do(http.MethodGet, "/clip/v2/resource", nil)
	if err != nil {
		return nil, err
	}
	resources := []*Resource{}
	if err := json.Unmarshal(r.Data, &resources); err != nil {
		return nil, err
	}
	return resources, nil
}

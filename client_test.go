package hue

import (
	"testing"

	"github.com/lampctl/go-hue/bridge/bridgetest"
)

func TestClient(t *testing.T) {
	for _, v := range []struct {
		Name string
		Fn   func(c *Client, s *bridgetest.Bridge) error
		Err  error
	}{
		{
			Name: "unauthenticated request",
			Fn: func(c *Client, s *bridgetest.Bridge) error {
				_, err := c.Resources()
				return err
			},
			Err: errForbiddenResponse,
		},
		{
			Name: "authenticated request",
			Fn: func(c *Client, s *bridgetest.Bridge) error {
				c.Username = bridgetest.Username
				_, err := c.Resources()
				return err
			},
			Err: nil,
		},
	} {
		s, err := bridgetest.New()
		if err != nil {
			t.Fatalf("%s: %s", v.Name, err)
		}
		defer s.Close()
		c := &Client{Host: s.URL}
		if err := v.Fn(c, s); err != v.Err {
			t.Fatalf("%s: %+v != %+v", v.Name, err, v.Err)
		}
	}
}

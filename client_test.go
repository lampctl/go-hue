package hue

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lampctl/go-hue/bridge/bridgetest"
)

func TestClient(t *testing.T) {
	for _, v := range []struct {
		Name        string
		Fn          func(c *Client, s *bridgetest.Bridge) error
		ReturnError bool
	}{
		{
			Name: "register without button press",
			Fn: func(c *Client, s *bridgetest.Bridge) error {
				return c.Register("")
			},
			ReturnError: true,
		},
		{
			Name: "register with only button press",
			Fn: func(c *Client, s *bridgetest.Bridge) error {
				s.PushButton()
				return c.Register("")
			},
			ReturnError: true,
		},
		{
			Name: "register with call and button press",
			Fn: func(c *Client, s *bridgetest.Bridge) error {
				if err := c.Register(""); err == nil {
					return errors.New("error expected without button press")
				}
				s.PushButton()
				if err := c.Register(""); err != nil {
					return err
				}
				if c.Username != bridgetest.Username {
					return fmt.Errorf("%+v != %+v", c.Username, bridgetest.Username)
				}
				return nil
			},
			ReturnError: false,
		},
		{
			Name: "unauthenticated request",
			Fn: func(c *Client, s *bridgetest.Bridge) error {
				_, err := c.Resources()
				return err
			},
			ReturnError: true,
		},
		{
			Name: "authenticated request",
			Fn: func(c *Client, s *bridgetest.Bridge) error {
				c.Username = bridgetest.Username
				_, err := c.Resources()
				return err
			},
			ReturnError: false,
		},
	} {
		s, err := bridgetest.New()
		if err != nil {
			t.Fatalf("%s: %s", v.Name, err)
		}
		defer s.Close()
		c := &Client{Host: s.URL}
		var (
			fnErr       = v.Fn(c, s)
			returnError = fnErr != nil
		)
		if returnError != v.ReturnError {
			if v.ReturnError {
				t.Fatalf("%s: %s", v.Name, "error expected")
			} else {
				t.Fatalf("%s: %s", v.Name, fnErr)
			}
		}
	}
}

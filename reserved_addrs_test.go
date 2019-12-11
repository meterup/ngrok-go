package ngrok

import (
	"context"
	"net/url"
	"testing"
)

func TestReservedAddrsCreate(t *testing.T) {
	t.Parallel()
	s := newServer(addrsCreateResponse, 201)
	client := New(s.URL, "")
	data := url.Values{}
	data.Set("region", string(RegionUS))
	data.Set("description", "go:test")
	data.Set("metadata", "go:test")
	addr, err := client.ReservedAddrs.Create(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if addr.Description != "go:test" {
		t.Errorf("bad description, should be go.test, got %q", addr.Description)
	}
}

func TestAddrsDelete(t *testing.T) {
	t.Parallel()
	s := newServer([]byte{}, 204)
	client := New(s.URL, "")
	err := client.ReservedAddrs.Delete(context.Background(), "ra_3GCBoX98pv2bBjWqS9AfL")
	if err != nil {
		t.Fatal(err)
	}
}

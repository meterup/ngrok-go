package ngrok

import (
	"context"
	"net/url"
	"testing"
)

func TestCredsCreate(t *testing.T) {
	t.Parallel()
	s := newServer(credsCreateResponse, 200)
	client := New(s.URL, "")
	data := url.Values{}
	data.Set("description", "go:test")
	data.Set("metadata", "go:test")
	cred, err := client.Creds.Create(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	if cred.Description != "go:test" {
		t.Errorf("bad description, should be go.test, got %q", cred.Description)
	}
}

func TestCredsDelete(t *testing.T) {
	t.Parallel()
	s := newServer([]byte{}, 204)
	client := New(s.URL, "")
	err := client.Creds.Delete(context.Background(), "cr_1UoBXvrhYlYA9um85RKLjAGabBJ")
	if err != nil {
		t.Fatal(err)
	}
}

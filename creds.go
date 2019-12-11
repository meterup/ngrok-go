package ngrok

import (
	"context"
	"net/url"
	"time"
)

type CredService struct {
	client *Client
}

// https://ngrok.com/docs/ngrok-link#list-credentials
type Cred struct {
	// Only present on the response from a Create request, otherwise empty
	Token       string    `json:"token"`
	Description string    `json:"description"`
	ACL         []string  `json:"acl"`
	CreatedAt   time.Time `json:"created_at"`
	ID          string    `json:"id"`
	Metadata    string    `json:"metadata"`
	URI         string    `json:"uri"`
}

const credPathPart = "/credentials"

func (c *CredService) Create(ctx context.Context, data url.Values) (*Cred, error) {
	cred := new(Cred)
	err := c.client.CreateResource(ctx, credPathPart, data, cred)
	return cred, err
}

func (c *CredService) Delete(ctx context.Context, id string) error {
	return c.client.DeleteResource(ctx, credPathPart, url.PathEscape(id))
}

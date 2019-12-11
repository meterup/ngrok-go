package ngrok

import (
	"context"
	"net/url"
	"time"
)

type ReservedAddrService struct {
	client *Client
}

// https://ngrok.com/docs/ngrok-link#list-reserved-addrs
type ReservedAddr struct {
	Addr        string    `json:"addr"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	ID          string    `json:"id"`
	Metadata    string    `json:"metadata"`
	Region      string    `json:"region"`
	URI         string    `json:"uri"`

	EndpointConfiguration map[string]string `json:"endpoint_configuration"`
}

const reservedAddrsPathPart = "/reserved_addrs"

// https://ngrok.com/docs/ngrok-link#reserve-addr
func (c *ReservedAddrService) Create(ctx context.Context, data url.Values) (*ReservedAddr, error) {
	addr := new(ReservedAddr)
	err := c.client.CreateResource(ctx, reservedAddrsPathPart, data, addr)
	return addr, err
}

// https://ngrok.com/docs/ngrok-link#reserve-addr
func (c *ReservedAddrService) Delete(ctx context.Context, id string) error {
	return c.client.DeleteResource(ctx, reservedAddrsPathPart, url.PathEscape(id))
}

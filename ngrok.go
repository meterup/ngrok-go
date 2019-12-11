// Package ngrok simplifies interaction with the ngrok API.
//
// This client is not endpoint compatible with the API - many endpoints are
// missing. Still, it should be trivial to add new *Service objects and
// endpoints, following the existing pattern.
//
// For more information on the ngrok API, see the documentation:
//
// https://ngrok.com/docs/ngrok-link#service-api
package ngrok

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/kevinburke/rest"
)

var defaultUserAgent string

const BaseURL = "https://api.ngrok.com"

const Version = "0.1"

func init() {
	gv := strings.Replace(runtime.Version(), "go", "", 1)
	defaultUserAgent = fmt.Sprintf("ngrok-go/%s go/%s (%s/%s)", Version, gv,
		runtime.GOOS, runtime.GOARCH)
}

type Client struct {
	*rest.Client
	userAgent string
	token     string

	Creds         *CredService
	ReservedAddrs *ReservedAddrService
}

const accept = "application/json"

func New(baseURL string, token string) *Client {
	restclient := rest.NewClient("", "", baseURL)
	restclient.ErrorParser = errorParser
	c := &Client{
		Client:    restclient,
		token:     token,
		userAgent: defaultUserAgent,
	}
	c.Creds = &CredService{c}
	c.ReservedAddrs = &ReservedAddrService{c}
	return c
}

type Region string

const RegionUS Region = "us"
const RegionEU Region = "eu"
const RegionAP Region = "ap"
const RegionAU Region = "au"

type ngrokError struct {
	ErrorCode  string          `json:"error_code"`
	Message    string          `json:"msg"`
	StatusCode int             `json:"status_code"`
	Details    json.RawMessage `json:"details"`
}

// errorParser tries to return a *rest.Error
func errorParser(resp *http.Response) error {
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	rerr := new(ngrokError)
	err = json.Unmarshal(resBody, rerr)
	if err != nil {
		return fmt.Errorf("invalid response body: %q", string(resBody))
	}
	if rerr.Message == "" {
		return fmt.Errorf("invalid response body: %q", string(resBody))
	}
	restError := &rest.Error{
		Title:  rerr.Message,
		Status: rerr.StatusCode,
		ID:     rerr.ErrorCode,
		Detail: string(rerr.Details),
	}
	return restError
}

// NewRequest creates a new signed request. The base URL will be prepended to
// the path, and the user-agent will also be attached.
func (c *Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.Client.Base+path, body)
	if err != nil {
		return nil, err
	}
	if ua := req.Header.Get("User-Agent"); ua == "" {
		req.Header.Set("User-Agent", c.userAgent)
	} else {
		req.Header.Set("User-Agent", c.userAgent+" "+ua)
	}
	req.Header.Set("Accept", accept)
	if body != nil {
		// ngrok doesn't allow charset
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-Ngrok-Version", "1")
	return req, nil
}

// from kevinburke/twilio-go/http.go

// GetResource retrieves an instance resource with the given path part (e.g.
// "/Messages") and sid (e.g. "MM123").
func (c *Client) GetResource(ctx context.Context, pathPart string, sid string, v interface{}) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	return c.MakeRequest(ctx, "GET", sidPart, nil, v)
}

// CreateResource makes a POST request to the given resource.
func (c *Client) CreateResource(ctx context.Context, pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest(ctx, "POST", pathPart, data, v)
}

func (c *Client) UpdateResource(ctx context.Context, pathPart string, sid string, data url.Values, v interface{}) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	return c.MakeRequest(ctx, "POST", sidPart, data, v)
}

func (c *Client) DeleteResource(ctx context.Context, pathPart string, sid string) error {
	sidPart := strings.Join([]string{pathPart, sid}, "/")
	err := c.MakeRequest(ctx, "DELETE", sidPart, nil, nil)
	if err == nil {
		return nil
	}
	rerr, ok := err.(*rest.Error)
	if ok && rerr.Status == http.StatusNotFound {
		return nil
	}
	return err
}

func (c *Client) ListResource(ctx context.Context, pathPart string, data url.Values, v interface{}) error {
	return c.MakeRequest(ctx, "GET", pathPart, data, v)
}

// Make a request to the ngrok API.
func (c *Client) MakeRequest(ctx context.Context, method string, pathPart string, data url.Values, v interface{}) error {
	rb := new(strings.Reader)
	if data != nil && (method == "POST" || method == "PUT") {
		rb = strings.NewReader(data.Encode())
	}
	if method == "GET" && data != nil {
		pathPart = pathPart + "?" + data.Encode()
	}
	req, err := c.NewRequest(method, pathPart, rb)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	return c.Do(req, &v)
}

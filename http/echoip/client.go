package echoip

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	echoip "github.com/mpolden/echoip/http"
)

// Options configure the Client's lookup behaviour.
type Options struct {
	// If IP is not nil, information will be looked up for this IP address
	// instead of the origin address of the request.
	IP net.IP
}

// Client can lookup IP information from an echoip service (see:
// https://github.com/mpolden/echoip).
type Client struct {
	*http.Client

	BaseURL string
}

// NewClient creates a new *Client for the echoip service reachable at baseURL.
func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

// Lookup performs an IP lookup against the echoip service with given opts. The
// context can be used to cancel the request.
func (c *Client) Lookup(ctx context.Context, opts *Options) (*echoip.Response, error) {
	url := fmt.Sprintf("%s/json", strings.TrimRight(c.BaseURL, "/"))
	if opts != nil && opts.IP != nil {
		url = fmt.Sprintf("%s?ip=%s", url, opts.IP)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse

		err = json.Unmarshal(body, &errResp)
		if err != nil {
			return nil, err
		}

		return nil, &errResp
	}

	var v echoip.Response

	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}

	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json")

	return client.Do(req)
}

// ErrorResponse is the error returned by *Client.Lookup if the service returns
// a non-200 response.
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

// Error implements the error interface.
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("echoip service returned HTTP status %d: %s", e.Status, e.Message)
}

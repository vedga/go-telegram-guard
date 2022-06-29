package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// MethodGET is HTTP request method GET
	MethodGET = "GET"
	// MethodPOST is HTTP request method POST
	MethodPOST = "POST"
	// MethodPUT is HTTP request method PUT
	MethodPUT = "PUT"
	// MethodDELETE is HTTP request method DELETE
	MethodDELETE      = "DELETE"
	headerContentType = "Content-Type"
	mediaTypeJSON     = "application/json"
	fmtTelegramAPI    = `https://api.telegram.org/bot%s/%s`
)

// Client is implementation HTTP(s) client
type Client struct {
	httpClient *http.Client
}

var (
	httpClient *Client
)

func init() {
	httpClient = NewClient()
}

// NewClient create control client with specified base URL and transport
func NewClient() *Client {
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}

	return &Client{
		httpClient: &http.Client{
			Transport: transport,
		},
	}
}

// Call execute ARI request and decode response if necessary
func (c *Client) Call(ctx context.Context, method, botToken, url string, resp interface{}, req interface{}) (err error) {
	var reqBody io.Reader
	if req != nil {
		reqBody, err = structToRequestBody(req)
		if nil != err {
			return err
		}
	}

	var r *http.Request

	if r, err = http.NewRequest(method, fmt.Sprintf(fmtTelegramAPI, botToken, url), reqBody); nil != err {
		return err
	}

	r.Header.Set(headerContentType, mediaTypeJSON)

	r = r.WithContext(ctx)

	return c.do(ctx, r, resp)
}

func structToRequestBody(req interface{}) (io.Reader, error) {
	buf := new(bytes.Buffer)

	if req != nil {
		if err := json.NewEncoder(buf).Encode(req); err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)

	res, err := c.httpClient.Do(req)
	if err != nil {
		// Request can't be sent or response not received
		return err
	}

	if err := decodeJSON(res.Header, res.Body, v); err != nil {
		return err
	}

	return nil

}

func decodeJSON(header http.Header, res io.ReadCloser, v interface{}) (err error) {
	var b bytes.Buffer

	_, err = io.Copy(&b, res)
	if err != nil {
		return
	}

	defer func() {
		if bodyErr := res.Close(); bodyErr != nil {
			err = bodyErr
		}
	}()

	// If v is nil it means we are not interested in decoding response.
	// It might be because of empty body or non-JSON content type.
	if v == nil {
		return nil
	}

	// Decode response only for Content-Type: application/json
	contentType := strings.ToLower(header.Get(headerContentType))
	if !strings.HasPrefix(contentType, mediaTypeJSON) {
		return fmt.Errorf("got %q response, expected "+mediaTypeJSON+": %w", contentType, err)
	}

	if jsonErr := json.NewDecoder(&b).Decode(v); jsonErr != nil {
		return err
	}

	return nil
}

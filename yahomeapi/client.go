package yahomeapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultHost = "https://api.iot.yandex.net"
)

type APIClient interface {
	GetDeviceInfo(ctx context.Context, deviceID string) (*GetDeviceInfoResponse, error)
	SendActionsForDevice(ctx context.Context, actions SendActionsRequest) (*SendActionsResponse, error)
}

type Config struct {
	APIHost        string
	BearerToken    string
	RequestTimeout time.Duration
}

type Client struct {
	APIHost     string
	BearerToken string
	httpClient  *http.Client
}

type rawResponse struct {
	httpCode int
	respBody []byte
}

var _ APIClient = &Client{}

func NewClient(config *Config) *Client {
	client := &Client{
		APIHost:     defaultHost,
		BearerToken: config.BearerToken,
		httpClient: &http.Client{
			Timeout: config.RequestTimeout,
		},
	}
	if config.APIHost != "" {
		client.APIHost = config.APIHost
	}
	return client
}

func (c *Client) get(ctx context.Context, URI string) (*rawResponse, error) {
	return c.do(ctx, URI, http.MethodGet, nil)
}

func (c *Client) post(ctx context.Context, URI string, body interface{}) (*rawResponse, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, URI, http.MethodPost, jsonBytes)
}

func (c *Client) do(ctx context.Context, URI, method string, reqBody []byte) (rawResp *rawResponse, err error) {
	var reqReader io.Reader = nil
	if reqBody != nil {
		reqReader = bytes.NewBuffer(reqBody)
	}

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, method, c.APIHost+URI, reqReader); err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	if resp, err = c.httpClient.Do(req); err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	rawResp = &rawResponse{
		httpCode: resp.StatusCode,
	}
	if rawResp.respBody, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("response reading: %w", err)
	}

	return rawResp, nil
}

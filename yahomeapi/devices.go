package yahomeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetDeviceInfo returns all available device info
func (c *Client) GetDeviceInfo(ctx context.Context, deviceID string) (*GetDeviceInfoResponse, error) {
	resp, err := c.get(ctx, "/v1.0/devices/"+deviceID)
	if err != nil {
		return nil, fmt.Errorf("yahomeapi: GetDevice failed %w", err)
	}

	var result GetDeviceInfoResponse
	if err = json.Unmarshal(resp.respBody, &result); err != nil {
		return nil, fmt.Errorf("yahomeapi: GetDevice response unmarshal failed %w", err)
	}

	if resp.httpCode != http.StatusOK {
		return &result, fmt.Errorf("yahomeapi: Incorrect response code %d", resp.httpCode)
	}

	return &result, nil
}

// SendActionsForDevice execute state changes for device list
func (c *Client) SendActionsForDevice(ctx context.Context, actions SendActionsRequest) (*SendActionsResponse, error) {
	body, err := c.post(ctx, "/v1.0/devices/actions", actions)
	if err != nil {
		return nil, fmt.Errorf("yahomeapi: SendActionsForDevice failed %w", err)
	}

	var result SendActionsResponse
	if err = json.Unmarshal(body.respBody, &result); err != nil {
		return nil, fmt.Errorf("yahomeapi: SendActionsForDevice response unmarshal failed %w", err)
	}

	return &result, nil
}

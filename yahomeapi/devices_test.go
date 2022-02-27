package yahomeapi

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClient_GetDeviceInfo(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.String() {
		case "/v1.0/devices/correct":
			_, _ = rw.Write([]byte(getDeviceCorrectResponse))
			return
		case "/v1.0/devices/notfound":
			_, _ = rw.Write([]byte(getDeviceNotFoundResponse))
			return
		case "/v1.0/devices/incorrect":
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte("{}"))
			return
		default:
			t.Errorf("unexpected URL: %s", req.URL.String())
		}
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	client := NewClient(&Config{
		APIHost:     server.URL,
		BearerToken: "BearerToken",
	})
	ctx := context.Background()

	type args struct {
		deviceId string
	}
	getInfoTests := []struct {
		name     string
		args     args
		wantResp *GetDeviceInfoResponse
		wantErr  error
	}{
		{
			name: "incorrect device request",
			args: args{
				deviceId: "incorrect",
			},
			wantResp: &GetDeviceInfoResponse{},
			wantErr:  errors.New("yahomeapi: Incorrect response code 500"),
		},
		{
			name: "correct device request",
			args: args{
				deviceId: "correct",
			},
			wantResp: &GetDeviceInfoResponse{
				PlatformStatusResponse: PlatformStatusResponse{
					RequestID: "request_id",
					Status:    "ok",
					Message:   "",
				},
				DeviceResponse: DeviceResponse{
					Id:         "id",
					Name:       "name",
					Aliases:    []string{},
					Type:       "devices.types.light",
					State:      "online",
					Groups:     []string{},
					Room:       "room",
					ExternalId: "external_id",
					SkillId:    "skill_id",
					Capabilities: []Capability{
						{
							Retrievable: true,
							Type:        "devices.capabilities.on_off",
							Parameters: map[string]interface{}{
								"split": false,
							},
							State: map[string]interface{}{
								"instance": "on",
								"value":    false,
							},
							LastUpdated: 1645961915.3663018,
						},
						{
							Retrievable: true,
							Type:        "devices.capabilities.range",
							Parameters: map[string]interface{}{
								"instance":      "brightness",
								"unit":          "unit.percent",
								"random_access": true,
								"looped":        false,
								"range": map[string]interface{}{
									"min":       float64(0),
									"max":       float64(100),
									"precision": float64(1),
								},
							},
							State: map[string]interface{}{
								"instance": "brightness",
								"value":    float64(100),
							},
							LastUpdated: 1645961915.3663018,
						},
						{
							Retrievable: true,
							Type:        "devices.capabilities.color_setting",
							Parameters: map[string]interface{}{
								"color_model": "rgb",
								"temperature_k": map[string]interface{}{
									"min": float64(4500),
									"max": float64(4500),
								},
							},
							State: map[string]interface{}{
								"instance": "rgb",
								"value":    float64(65280),
							},
							LastUpdated: 1645961915.3663018,
						},
					},
					Properties: []Property{},
				},
			},
		},
		{
			name: "not found device request",
			args: args{
				deviceId: "notfound",
			},
			wantResp: &GetDeviceInfoResponse{
				PlatformStatusResponse: PlatformStatusResponse{
					RequestID: "request_id",
					Status:    "error",
					Message:   "device not found",
				},
				DeviceResponse: DeviceResponse{},
			},
		},
	}
	for _, tt := range getInfoTests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.GetDeviceInfo(ctx, tt.args.deviceId)
			if diff := cmp.Diff(tt.wantResp, resp); diff != "" {
				t.Errorf("GetDeviceInfo() mismatch (-want +got):\n%s", diff)
			}
			if err != tt.wantErr && err.Error() != tt.wantErr.Error() {
				t.Errorf("GetDeviceInfo() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SendActionsForDevice(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Errorf("body read error %s", err)
			return
		}
		originRequest := &SendActionsRequest{}
		err = json.Unmarshal(body, originRequest)
		if err != nil {
			t.Errorf("unmarshal origin request error %s", err)
			return
		}
		if len(originRequest.Devices) < 1 {
			t.Errorf("empty device list in origin request")
			return
		}
		var bodyCheck string
		var resp string
		switch originRequest.Devices[0].Id {
		case "correct":
			bodyCheck = sendActionPayload
			resp = sendActionResponse
		}
		if bodyCheck != string(body) {
			t.Errorf("SendActionsForDevice() unexpected payload, got = %v, want %v", string(body), bodyCheck)
		}
		_, _ = rw.Write([]byte(resp))
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server
	client := NewClient(&Config{
		APIHost:     server.URL,
		BearerToken: "BearerToken",
	})
	ctx := context.Background()

	type args struct {
		request SendActionsRequest
	}
	sendActionsTests := []struct {
		name        string
		args        args
		wantResp    *SendActionsResponse
		wantPayload []byte
		wantErr     error
	}{
		// {
		// 	name: "incorrect device request",
		// 	args: args{
		// 		deviceId: "incorrect",
		// 	},
		// 	wantResp: &SendActionsForDeviceResponse{},
		// 	wantErr:  errors.New("yahomeapi: Incorrect response code 500"),
		// },
		{
			name: "correct device request",
			args: args{
				request: SendActionsRequest{
					Devices: []DeviceActions{{
						Id: "correct",
						Actions: []DeviceAction{
							{
								Type: "devices.capabilities.range",
								State: CapabilityStateIntObject{
									Instance: "brightness",
									Value:    100,
								},
							},
						},
					}},
				},
			},
			wantResp: &SendActionsResponse{
				PlatformStatusResponse: PlatformStatusResponse{
					RequestID: "request_id",
					Status:    "ok",
				},
				Devices: []Device{
					{
						Id: "correct",
						Capabilities: []CapabilityResponse{
							{
								Type: "devices.capabilities.range",
								State: StateChangeResult{
									Instance: "brightness",
									ActionResult: ActionResult{
										Status: "DONE",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range sendActionsTests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.SendActionsForDevice(ctx, tt.args.request)
			if diff := cmp.Diff(tt.wantResp, resp); diff != "" {
				t.Errorf("SendActionsForDevice() mismatch (-want +got):\n%s", diff)
			}
			if err != tt.wantErr && (err == nil || tt.wantErr == nil || err.Error() != tt.wantErr.Error()) {
				t.Errorf("SendActionsForDevice() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

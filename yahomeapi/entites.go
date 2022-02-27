package yahomeapi

type SendActionsRequest struct {
	Devices []DeviceActions `json:"devices"`
}

type DeviceActions struct {
	Id      string         `json:"id"`
	Actions []DeviceAction `json:"actions"`
}

type DeviceAction struct {
	Type  string          `json:"type"`
	State CapabilityState `json:"state"`
}

type CapabilityState interface {
}

type CapabilityStateBoolObject struct {
	Instance string `json:"instance"`
	Value    bool   `json:"value"`
}
type CapabilityStateIntObject struct {
	Instance string `json:"instance"`
	Value    int32  `json:"value"`
}

type SendActionsResponse struct {
	PlatformStatusResponse
	SendActionsDeviceResponse
}

type CapabilityResponse struct {
	Type  string          `json:"type"`
	State CapabilityState `json:"state"`
}

type Device struct {
	Id           string               `json:"id"`
	Capabilities []CapabilityResponse `json:"capabilities"`
}

type SendActionsDeviceResponse struct {
	Devices []Device `json:"devices"`
}

type GetDeviceInfoResponse struct {
	PlatformStatusResponse
	DeviceResponse
}

type PlatformStatusResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type DeviceResponse struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Aliases      []string     `json:"aliases"`
	Type         string       `json:"type"`
	State        string       `json:"state"`
	Groups       []string     `json:"groups"`
	Room         string       `json:"room"`
	ExternalId   string       `json:"external_id"`
	SkillId      string       `json:"skill_id"`
	Capabilities []Capability `json:"capabilities"`
	Properties   []Property   `json:"properties"`
}

type Capability struct {
	Retrievable bool        `json:"retrievable"`
	Type        string      `json:"type"`
	Parameters  interface{} `json:"parameters"`
	State       interface{} `json:"state"`
	LastUpdated float32     `json:"last_updated"`
}

type Property struct {
	Retrievable bool        `json:"retrievable"`
	Type        string      `json:"type"`
	Parameters  interface{} `json:"parameters"`
	State       interface{} `json:"state"`
	LastUpdated float32     `json:"last_updated"`
}

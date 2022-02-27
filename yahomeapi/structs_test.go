package yahomeapi

var (
	getDeviceCorrectResponse = `{
    "status": "ok",
    "request_id": "request_id",
    "id": "id",
    "name": "name",
    "aliases": [],
    "type": "devices.types.light",
    "external_id": "external_id",
    "skill_id": "skill_id",
    "state": "online",
    "groups": [],
    "room": "room",
    "capabilities": [
        {
            "retrievable": true,
            "type": "devices.capabilities.on_off",
            "parameters": {
                "split": false
            },
            "state": {
                "instance": "on",
                "value": false
            },
            "last_updated": 1645961915.3663018
        },
        {
            "retrievable": true,
            "type": "devices.capabilities.range",
            "parameters": {
                "instance": "brightness",
                "unit": "unit.percent",
                "random_access": true,
                "looped": false,
                "range": {
                    "min": 0,
                    "max": 100,
                    "precision": 1
                }
            },
            "state": {
                "instance": "brightness",
                "value": 100
            },
            "last_updated": 1645961915.3663018
        },
        {
            "retrievable": true,
            "type": "devices.capabilities.color_setting",
            "parameters": {
                "color_model": "rgb",
                "temperature_k": {
                    "min": 4500,
                    "max": 4500
                }
            },
            "state": {
                "instance": "rgb",
                "value": 65280
            },
            "last_updated": 1645961915.3663018
        }
    ],
    "properties": []
}`

	getDeviceNotFoundResponse = `{
    "request_id": "request_id",
    "status": "error",
    "message": "device not found"
}`

	sendActionPayload  = `{"devices":[{"id":"correct","actions":[{"type":"devices.capabilities.range","state":{"instance":"brightness","value":100}}]}]}`
	sendActionResponse = `{"status":"ok","request_id":"request_id","devices":[{"id":"correct","capabilities":[{"type":"devices.capabilities.range","state":{"instance":"brightness","action_result":{"status":"DONE"}}}]}]}`
)

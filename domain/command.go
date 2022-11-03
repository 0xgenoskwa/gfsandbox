package domain

type CommandType int32

const (
	CommandTypeInformation CommandType = iota
	CommandTypeSetup
	CommandTypeOpenUrl
	CommandTypeSendTouchEvent
)

type CommandInformationResponse struct {
	DeviceName string `json:"device_name"`
	MacAddress string `json:"mac_address"`
	WifiSsid   string `json:"wifi_ssid"`
	WifiPsk    string `json:"wifi_psk"`
	MqttUrl    string `json:"mqtt_url"`
	MqttPort   int8   `json:"mqtt_port"`
	FaChannel  string `json:"fa_channel"`
	FdChannel  string `json:"fd_channel"`
}

type CommandSetup struct {
	WifiSsid  string `json:"wifi_ssid"`
	WifiPsk   string `json:"wifi_psk"`
	MqttUrl   string `json:"mqtt_url"`
	MqttPort  int8   `json:"mqtt_port"`
	FaChannel string `json:"fa_channel"`
	FdChannel string `json:"fd_channel"`
}

type CommandSetupResponse struct {
	MacAddress string `json:"mac_address"`
}

type CommandOpenUrl struct {
	Url string `json:"url"`
}

type CommandOpenUrlResponse struct {
	Result bool `json:"result"`
}

type CommandSendTouchEvent struct {
	AltKey         bool    `json:"altKey"`
	ChangedTouches []Touch `json:"changedTouches"`
	CtrlKey        bool    `json:"ctrlKey"`
	MetaKey        bool    `json:"metaKey"`
	ShiftKey       bool    `json:"shiftKey"`
	TargetTouches  []Touch `json:"targetTouches"`
	Touches        []Touch `json:"touches"`
}

type CommandSendTouchEventResponse struct {
	Result bool `json:"result"`
}

type ErrorResponse struct {
	Data  []byte `json:"data"`
	Error string `json:"error"`
}

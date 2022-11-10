package domain

type CommandType int32

const (
	CommandTypeInformation CommandType = iota
	CommandTypeSetup
	CommandTypeOpenUrl
	CommandTypeSendTouchEvent
	CommandTypeSendKeyEvent
	CommandTypeToggleDim
	CommandTypeGetMacAddress
)

type CommandInformationResponse struct {
	DeviceName string `json:"device_name"`
	MacAddress string `json:"mac_address"`
	WifiSsid   string `json:"wifi_ssid"`
	WifiPsk    string `json:"wifi_psk"`
	MqttUrl    string `json:"mqtt_url"`
	MqttPort   int    `json:"mqtt_port"`
	FaChannel  string `json:"fa_channel"`
	FdChannel  string `json:"fd_channel"`

	ScreenWidth  string `json:"screen_width"`
	ScreenHeight string `json:"screen_height"`

	Version string `json:"version"`
	Build   string `json:"build"`
}

type CommandSetup struct {
	WifiSsid  string `json:"wifi_ssid"`
	WifiPsk   string `json:"wifi_psk"`
	MqttUrl   string `json:"mqtt_url"`
	MqttPort  int    `json:"mqtt_port"`
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
	Event string     `json:"event"`
	Data  TouchEvent `json:"data"`
}

type CommandSendTouchEventResponse struct {
	Result bool `json:"result"`
}

type CommandSendKeyEvent struct {
	Key string `json:"event"`
}

type CommandSendKeyEventResponse struct {
	Result bool `json:"result"`
}

type ErrorResponse struct {
	Data  []byte `json:"data"`
	Error string `json:"error"`
}

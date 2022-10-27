package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/0xgenoskwa/gfsandbox/chrome"
	"github.com/0xgenoskwa/gfsandbox/config"
	"github.com/0xgenoskwa/gfsandbox/domain"
	"github.com/0xgenoskwa/gfsandbox/system"
	"github.com/0xgenoskwa/gfsandbox/wifi"
)

type Handler struct {
	Chrome *chrome.Chrome
	Config *config.Config
	Wifi   *wifi.Wifi
}

func ProvideHandler(c *chrome.Chrome, cfg *config.Config, w *wifi.Wifi) *Handler {
	return &Handler{
		Chrome: c,
		Config: cfg,
		Wifi:   w,
	}
}

func byteToInt(bytes []byte) int {
	result := 0
	for i := 0; i < 4; i++ {
		result = result << 8
		result += int(bytes[i])

	}

	return result
}

func (h *Handler) OnData(data []byte) ([]byte, error) {
	dataStr := string(data)
	dataParts := strings.Split(dataStr, ":")
	firstNumber, err := strconv.Atoi(dataParts[0])
	if err != nil {
		return nil, err
	}
	cmdType := domain.CommandType(firstNumber)
	var msg []byte
	if len(dataParts) > 0 {
		msg = []byte(strings.Join(dataParts[1:], ":"))
	}

	switch cmdType {
	case domain.CommandTypeInformation:
		return h.getInformation()
	case domain.CommandTypeSetup:
		return h.setup(msg)
	case domain.CommandTypeOpenUrl:
		return h.openUrl(msg)
	default:
		return nil, errors.New("unknown cmd")
	}
}

func (h *Handler) getInformation() ([]byte, error) {
	fmt.Printf("process cmd information")
	macAddr, err := system.GetWirelessMacAddr()
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	data["mac_address"] = macAddr
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (h *Handler) setup(msg []byte) ([]byte, error) {
	var setupCmd domain.CommandSetup
	err := json.Unmarshal(msg, &setupCmd)
	if err != nil {
		return nil, err
	}
	fmt.Printf("process cmd setup %+v\n", setupCmd)
	// TODO init wifi
	connectOutput, err := h.Wifi.Connect(setupCmd.WifiSsid, setupCmd.WifiPsk)
	if err != nil {
		return nil, errors.New("incorrect ssid/psk")
	}
	fmt.Println("Connected", string(connectOutput))

	// Save config
	h.Config.WifiSsid = setupCmd.WifiSsid
	h.Config.WifiPsk = setupCmd.WifiPsk
	h.Config.MqttUrl = setupCmd.MqttUrl
	h.Config.MqttPort = setupCmd.MqttPort
	h.Config.FaChannel = setupCmd.FaChannel
	h.Config.FdChannel = setupCmd.FdChannel
	if err := h.Config.SaveConfig(); err != nil {
		return nil, err
	}

	data, err := h.getInformation()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (h *Handler) openUrl(msg []byte) ([]byte, error) {
	var ouCmd domain.CommandOpenUrl
	err := json.Unmarshal(msg, &ouCmd)
	if err != nil {
		return nil, err
	}
	fmt.Printf("process cmd open url %+v\n", ouCmd)
	if err := h.Chrome.OpenUrl(ouCmd.Url); err != nil {
		return nil, err
	}

	return nil, nil
}

package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

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

func (h *Handler) OnData(data []byte) (domain.CommandType, []byte, error) {
	parts := bytes.Split(data, []byte(":"))
	firstNumber, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return -1, nil, err
	}
	cmdType := domain.CommandType(firstNumber)
	var msg []byte
	if len(parts) > 1 {
		msg = bytes.Join(parts[1:], []byte(":"))
	}

	switch cmdType {
	case domain.CommandTypeInformation:
		resp, err := h.getInformation()
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeSetup:
		resp, err := h.setup(msg)
		if err != nil {
			h.Chrome.Toast(err.Error())
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeOpenUrl:
		resp, err := h.openUrl(msg)
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	default:
		return -1, nil, errors.New("unknown cmd")
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
	_, err = h.Wifi.Connect(setupCmd.WifiSsid, setupCmd.WifiPsk)
	if err != nil {
		err := errors.New("incorrect ssid/psk")
		h.Chrome.Toast(err.Error())
		return nil, err
	}
	h.Chrome.Toast(fmt.Sprintf("Connect wifi %s success", setupCmd.WifiSsid))
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

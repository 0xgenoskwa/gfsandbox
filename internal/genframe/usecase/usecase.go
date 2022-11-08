package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"go.genframe.xyz/config"
	"go.genframe.xyz/domain"
	"go.genframe.xyz/pkg/chrome"
	"go.genframe.xyz/pkg/system"
	"go.genframe.xyz/pkg/wifi"
)

type Usecase struct {
	Chrome *chrome.Chrome
	Config *config.Config
	Wifi   *wifi.Wifi
}

func ProvideUsecase(c *chrome.Chrome, cfg *config.Config, w *wifi.Wifi) *Usecase {
	return &Usecase{
		Chrome: c,
		Config: cfg,
		Wifi:   w,
	}
}

func (u *Usecase) GetInformation() ([]byte, error) {
	fmt.Printf("process cmd information")
	macAddr, err := system.GetWirelessMacAddr()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("bash", "-c", "xdpyinfo | awk '/dimensions/{print $2}'")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	dParts := strings.Split(string(stdout), "x")

	data := domain.CommandInformationResponse{
		MacAddress:   macAddr,
		DeviceName:   u.Config.DeviceName,
		WifiSsid:     u.Config.WifiSsid,
		WifiPsk:      u.Config.WifiPsk,
		MqttUrl:      u.Config.MqttUrl,
		MqttPort:     u.Config.MqttPort,
		FaChannel:    u.Config.FaChannel,
		FdChannel:    u.Config.FdChannel,
		ScreenWidth:  strings.TrimSpace(dParts[0]),
		ScreenHeight: strings.TrimSpace(dParts[1]),
		Version:      u.Config.Version,
		Build:        u.Config.Build,
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (u *Usecase) Setup(msg []byte) ([]byte, error) {
	var setupCmd domain.CommandSetup
	err := json.Unmarshal(msg, &setupCmd)
	if err != nil {
		return nil, err
	}
	_, err = u.Wifi.Connect(setupCmd.WifiSsid, setupCmd.WifiPsk)
	if err != nil {
		err := errors.New("incorrect ssid/psk")
		u.Chrome.Toast(err.Error())
		return nil, err
	}
	u.Chrome.Toast(fmt.Sprintf("Connect wifi %s success", setupCmd.WifiSsid))
	// Save config
	u.Config.WifiSsid = setupCmd.WifiSsid
	u.Config.WifiPsk = setupCmd.WifiPsk
	u.Config.MqttUrl = setupCmd.MqttUrl
	u.Config.MqttPort = setupCmd.MqttPort
	u.Config.FaChannel = setupCmd.FaChannel
	u.Config.FdChannel = setupCmd.FdChannel
	if err := u.Config.SaveConfig(); err != nil {
		return nil, err
	}

	data, err := u.GetInformation()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *Usecase) OpenUrl(msg []byte) ([]byte, error) {
	var ouCmd domain.CommandOpenUrl
	err := json.Unmarshal(msg, &ouCmd)
	if err != nil {
		return nil, err
	}
	fmt.Printf("process cmd open url %+v\n", ouCmd)
	if err := u.Chrome.OpenUrl(ouCmd.Url); err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *Usecase) SendTouchEvent(msg []byte) ([]byte, error) {
	var stCmd domain.CommandSendTouchEvent
	err := json.Unmarshal(msg, &stCmd)
	if err != nil {
		return nil, err
	}
	fmt.Printf("process cmd send touch event %+v\n", stCmd)
	if err := u.Chrome.SendTouchEvent(stCmd.Event, stCmd.Data); err != nil {
		return nil, err
	}

	return nil, nil
}

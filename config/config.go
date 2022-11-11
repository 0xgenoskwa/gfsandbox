package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"go.genframe.xyz/pkg/system"
)

var Version string
var Build string

type Config struct {
	Path    string `json:"-"`
	Version string `json:"-"`
	Build   string `json:"-"`

	changed chan bool `json:"-"`

	DeviceName   string `json:"device_name"`
	WifiSsid     string `json:"wifi_ssid"`
	WifiPsk      string `json:"wifi_psk"`
	MqttUrl      string `json:"mqtt_url"`
	MqttPort     int    `json:"mqtt_port"`
	MqttUsername string `json:"mqtt_username"`
	MqttPassword string `json:"mqtt_password"`
	FaChannel    string `json:"fa_channel"`
	FdChannel    string `json:"fd_channel"`
	ServingUrl   string `json:"serving_url"`
}

func ProvideConfig() *Config {
	c := Config{
		Path:    "config.json",
		Version: Version,
		Build:   Build,
	}

	return &c
}

func (c *Config) LoadConfig() error {
	fmt.Println("Load config", c.Path)
	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		fmt.Println("Load config is not existed path", c.Path)
		if err := c.SaveConfig(); err != nil {
			return err
		}
	}
	fmt.Println("Load config existed path", c.Path)
	jsonFile, err := os.Open(c.Path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(bytes, c)
	if err != nil {
		return err
	}
	fmt.Println("Load config existed path end")
	return nil
}

func (c *Config) SaveConfig() error {
	if c.DeviceName == "" {
		macAddr, err := system.GetWirelessMacAddr()
		if err != nil {
			panic(err)
		}
		deviceName := fmt.Sprintf("Genframe#%s", strings.Replace(macAddr, ":", "", -1))
		c.DeviceName = deviceName
	}
	file, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(c.Path, file, 0644)
	if err != nil {
		return err
	}

	c.changed <- true

	return nil
}

func (c *Config) HasMqttConfig() bool {
	if c.MqttUrl != "" && c.MqttPort > 0 {
		return true
	}

	return false
}

func (c *Config) Changed() chan bool {
	return c.changed
}

package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var Version string

type Config struct {
	Path    string `json:"-"`
	Version string `json:"-"`

	DeviceName string `json:"device_name"`
	WifiSsid   string `json:"wifi_ssid"`
	WifiPsk    string `json:"wifi_psk"`
	MqttUrl    string `json:"mqtt_url"`
	MqttPort   int8   `json:"mqtt_port"`
	FaChannel  string `json:"fa_channel"`
	FdChannel  string `json:"fd_channel"`
}

func ProvideConfig() *Config {
	return &Config{
		Path:    "config.json",
		Version: Version,
	}
}

func (c *Config) LoadConfig() error {
	jsonFile, err := os.Open(c.Path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	fmt.Println("loaded config.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(bytes, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) SaveConfig() error {
	file, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(c.Path, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

package wifi

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type Wifi struct {
}

func ProvideWifi() *Wifi {
	return &Wifi{}
}

func (w *Wifi) Scan() ([]byte, error) {
	cmd := exec.Command("nmcli", "dev", "wifi", "rescan")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.New(stderr.String())
	}
	oldListSsid := ""
	for {
		listSsid, err := w.List()
		if err != nil {
			return nil, err
		}
		if oldListSsid == "" {
			oldListSsid = string(listSsid)
		} else {
			if oldListSsid != string(listSsid) {
				return listSsid, nil
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (w *Wifi) List() ([]byte, error) {
	cmd := exec.Command("nmcli", "dev", "wifi", "list")
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.New(stderr.String())
	}

	return out.Bytes(), nil
}

func (w *Wifi) Connect(ssid, psk string) ([]byte, error) {
	cmd := exec.Command("nmcli", "dev", "wifi", "connect", ssid, "password", psk)
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if strings.HasPrefix(stderr.String(), "Error: No network with SSID") {
			_, err := w.Scan()
			if err != nil {
				return nil, err
			}
			out.Reset()
			stderr.Reset()
			cmd := exec.Command("nmcli", "dev", "wifi", "connect", ssid, "password", psk)
			var out bytes.Buffer
			var stderr bytes.Buffer

			cmd.Stdout = &out
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return out.Bytes(), nil
}

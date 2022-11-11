package wifi

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Wifi struct {
	signal chan bool

	monitoring chan bool
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

func (w *Wifi) HasInternet() (ok bool) {
	if _, err := http.Get("http://clients3.google.com/generate_204"); err != nil {
		return false
	}
	return true
}

func (w *Wifi) StartWifiMonitoring() {
	fmt.Println("Wifi monitoring start")
	noConnCount := 0
	state := false
	for {
		select {
		case <-w.monitoring:
			fmt.Println("stop monitoring")
			return
		default:
			if w.HasInternet() {
				noConnCount = 0
				if !state {
					state = true
					w.signal <- state
				}
			} else {
				noConnCount = noConnCount + 1
				if noConnCount > 3 {
					state = false
					w.signal <- state
				}
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func (w *Wifi) StopWifiMonitoring() {
	fmt.Println("Wifi monitoring stop")
	go func() {
		w.monitoring <- true
	}()
}

func (w *Wifi) Signal() chan bool {
	return w.signal
}

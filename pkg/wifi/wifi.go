package wifi

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Wifi struct {
	signal chan bool
}

func ProvideWifi() *Wifi {
	return &Wifi{
		signal: make(chan bool),
	}
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
	fmt.Println("HasInternet start")
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	req, err := http.NewRequest("GET", "http://clients3.google.com/generate_204", nil)
	if err != nil {
		return false
	}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("HasInternet return false")
		return false
	}
	fmt.Println("HasInternet return true")
	return true
}

func (w *Wifi) StartWifiMonitoring() {
	fmt.Println("Wifi monitoring start")
	noConnCount := 0
	state := false
	for {
		fmt.Println("start check internet")
		if w.HasInternet() {
			fmt.Println("has internet")
			noConnCount = 0
			if !state {
				fmt.Println("return signal has internet")
				state = true
				w.signal <- state
				fmt.Println("end return signal has internet")
			}
		} else {
			noConnCount = noConnCount + 1
			fmt.Println("dont have internet", noConnCount)
			if noConnCount > 3 {
				fmt.Println("return signal dont have internet")
				state = false
				w.signal <- state
				fmt.Println("end return signal dont have internet")
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func (w *Wifi) Signal() <-chan bool {
	return w.signal
}

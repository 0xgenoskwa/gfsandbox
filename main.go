package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/0xgenoskwa/gfsandbox/chrome"
	"github.com/0xgenoskwa/gfsandbox/config"
	"github.com/0xgenoskwa/gfsandbox/handler"
	"github.com/0xgenoskwa/gfsandbox/protocol/bluetooth"
	"github.com/0xgenoskwa/gfsandbox/system"
	"github.com/0xgenoskwa/gfsandbox/wifi"
)

//go autoupdate https://github.com/IndioInc/go-autoupdate/blob/master/autoupdate/autoupdater.go

func main() {
	errChan := make(chan error)
	ctx := context.Background()

	cfg := config.ProvideConfig()
	if err := cfg.LoadConfig(); err != nil {
		panic(err)
	}
	if cfg.DeviceName == "" {
		macAddr, err := system.GetWirelessMacAddr()
		if err != nil {
			panic(err)
		}
		deviceName := fmt.Sprintf("Genframe#%s", strings.Replace(macAddr[len(macAddr)-5:], ":", "", -1))
		cfg.DeviceName = deviceName
		if err := cfg.SaveConfig(); err != nil {
			panic(err)
		}
	}
	w := wifi.ProvideWifi()
	chrome := chrome.ProvideChrome()
	err, cancel := chrome.Init(ctx)
	if err != nil {
		panic(err)
	}
	defer cancel()
	chrome.OpenHtml()
	handler := handler.ProvideHandler(chrome, cfg, w)

	bluetoothptl := bluetooth.ProvideBluetooth(cfg, handler)
	if err := bluetoothptl.Start(); err != nil {
		panic(err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Started genframe %s - %s \n", cfg.Version, cfg.Build)

	for {
		select {
		case <-stop:
			bluetoothptl.Stop()
			return
		case err := <-errChan:
			panic(err)
		}
	}
}

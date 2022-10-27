package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xgenoskwa/gfsandbox/chrome"
	"github.com/0xgenoskwa/gfsandbox/config"
	"github.com/0xgenoskwa/gfsandbox/handler"
	"github.com/0xgenoskwa/gfsandbox/protocol/bluetooth"
	"github.com/0xgenoskwa/gfsandbox/wifi"
)

func main() {
	errChan := make(chan error)
	println("starting")

	ctx := context.Background()

	cfg := config.ProvideConfig()
	if err := cfg.LoadConfig(); err != nil {
		errChan <- err
	}
	w := wifi.ProvideWifi()
	chrome := chrome.ProvideChrome()
	err, cancel := chrome.Init(ctx)
	if err != nil {
		errChan <- err
	}
	defer cancel()
	chrome.OpenHtml()
	handler := handler.ProvideHandler(chrome, cfg, w)

	bluetoothptl := bluetooth.ProvideBluetooth(cfg, handler)
	if err := bluetoothptl.Start(); err != nil {
		errChan <- err
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

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

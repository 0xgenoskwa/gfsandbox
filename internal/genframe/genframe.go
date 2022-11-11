package genframe

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.genframe.xyz/config"
	"go.genframe.xyz/internal/genframe/controller/bluetooth"
	"go.genframe.xyz/internal/genframe/controller/mqtt"
	"go.genframe.xyz/internal/genframe/usecase"
	"go.genframe.xyz/pkg/chrome"
	"go.genframe.xyz/pkg/wifi"
)

//go autoupdate https://github.com/IndioInc/go-autoupdate/blob/master/autoupdate/autoupdater.go

type Genframe struct {
	Config *config.Config
	Wifi   *wifi.Wifi
	Chrome *chrome.Chrome
	//
	Usecase          *usecase.Usecase
	BluetoothHandler *bluetooth.Bluetooth
	MqttHandler      *mqtt.Mqtt

	notify chan error
}

func ProvideGenframe(c *config.Config, w *wifi.Wifi, chr *chrome.Chrome, u *usecase.Usecase, b *bluetooth.Bluetooth, m *mqtt.Mqtt) *Genframe {
	return &Genframe{
		Config:           c,
		Wifi:             w,
		Chrome:           chr,
		Usecase:          u,
		BluetoothHandler: b,
		MqttHandler:      m,
		notify:           make(chan error, 1),
	}
}

func (g *Genframe) Run() {
	ctx := context.Background()

	cancel, err := g.Chrome.Init(ctx)
	if err != nil {
		panic(err)
	}
	defer cancel()
	// open html
	g.Chrome.OpenHtml()

	// start bluetooth
	if err := g.BluetoothHandler.Start(); err != nil {
		panic(err)
	}
	// start mqtt
	if err := g.MqttHandler.Start(); err != nil {
		panic(err)
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Println("app - Run - signal: " + s.String())
	case err = <-g.BluetoothHandler.Notify():
		fmt.Println(fmt.Errorf("app - Run - g.BluetoothHandler.Notify: %w", err))
	case err = <-g.MqttHandler.Notify():
		fmt.Println(fmt.Errorf("app - Run - g.MqttHandler.Notify: %w", err))
	}

	if err := g.BluetoothHandler.Shutdown(); err != nil {
		fmt.Println(fmt.Errorf("app - Run - g.BluetoothHandler.Shutdown: %w", err))
	}
	if err := g.MqttHandler.Shutdown(); err != nil {
		fmt.Println(fmt.Errorf("app - Run - g.MqttHandler.Shutdown: %w", err))
	}
}

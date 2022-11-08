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
	fmt.Println("genframe Chrome run 1")
	ctx := context.Background()

	cancel, err := g.Chrome.Init(ctx)
	if err != nil {
		panic(err)
	}
	defer cancel()
	fmt.Println("genframe Chrome run 1")
	g.Chrome.OpenHtml()
	fmt.Println("genframe Chrome run 2")

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
}

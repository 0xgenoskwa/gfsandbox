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

	if err := g.Config.LoadConfig(); err != nil {
		panic(err)
	}

	cancel, err := g.Chrome.Init(ctx)
	if err != nil {
		panic(err)
	}
	defer cancel()
	// open html
	g.Chrome.OpenHtml()
	if g.Config.ServingUrl != "" {
		g.Chrome.OpenUrl(g.Config.ServingUrl)
	}

	go g.Wifi.StartWifiMonitoring()
	defer g.Wifi.StopWifiMonitoring()

	// start bluetooth
	if err := g.BluetoothHandler.Start(); err != nil {
		panic(err)
	}

	// start mqtt
	if g.Wifi.HasInternet() && g.Config.HasMqttConfig() {
		if err := g.MqttHandler.Start(); err != nil {
			panic(err)
		}
	}

	// monitoring network logic
	quitMonitoring := make(chan bool)
	go func() {
		select {
		case <-quitMonitoring:
			return
		case status := <-g.Wifi.Signal():
			if status && g.Config.HasMqttConfig() && !g.MqttHandler.Started {
				if err := g.MqttHandler.Start(); err != nil {
					panic(err)
				}
			}
			if !status && g.MqttHandler.Started {
				g.MqttHandler.Shutdown()
			}
		case <-g.Config.Changed():
			networkStatus := g.Wifi.HasInternet()
			if networkStatus && g.Config.HasMqttConfig() {
				if g.MqttHandler.Started {
					g.MqttHandler.Shutdown()
				}
				g.MqttHandler.Start()
			}
		}
	}()

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

	<-quitMonitoring

	if g.BluetoothHandler.Started {
		if err := g.BluetoothHandler.Shutdown(); err != nil {
			fmt.Println(fmt.Errorf("app - Run - g.BluetoothHandler.Shutdown: %w", err))
		}
	}

	if g.MqttHandler.Started {
		if err := g.MqttHandler.Shutdown(); err != nil {
			fmt.Println(fmt.Errorf("app - Run - g.MqttHandler.Shutdown: %w", err))
		}
	}
}

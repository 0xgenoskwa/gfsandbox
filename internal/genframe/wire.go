//go:build wireinject

package genframe

import (
	"github.com/google/wire"
	"go.genframe.xyz/config"
	"go.genframe.xyz/internal/genframe/controller/bluetooth"
	"go.genframe.xyz/internal/genframe/controller/mqtt"
	"go.genframe.xyz/internal/genframe/usecase"
	"go.genframe.xyz/pkg/chrome"
	"go.genframe.xyz/pkg/wifi"
)

func InitializeGenframe() *Genframe {
	panic(wire.Build(
		config.ProviderSet,
		wifi.ProviderSet,
		chrome.ProviderSet,
		bluetooth.ProviderSet,
		mqtt.ProviderSet,
		usecase.ProviderSet,
		ProvideGenframe,
	))
}

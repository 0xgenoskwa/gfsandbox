//go:build wireinject

package bootstrap

import (
	"github.com/google/wire"
	"go.genframe.xyz/config"
)

func InitializeBootstrap() *Bootstrap {
	panic(wire.Build(
		config.ProviderSet,
		ProvideBootstrap,
	))
}

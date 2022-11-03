package wifi

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ProvideWifi,
)

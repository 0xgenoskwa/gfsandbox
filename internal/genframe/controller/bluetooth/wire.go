package bluetooth

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ProvideBluetooth,
)

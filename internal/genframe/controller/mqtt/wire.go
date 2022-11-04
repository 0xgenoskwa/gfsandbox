package mqtt

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ProvideMQTT,
)

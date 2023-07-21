package loader

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewConfigLoader, NewStatusLoader)

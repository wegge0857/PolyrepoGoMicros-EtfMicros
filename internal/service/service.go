package service

import "github.com/google/wire"

// ProviderSet is service providers.
// 加入新的服务，请将服务注册到 ProviderSet 中
var ProviderSet = wire.NewSet(
	NewEtfService,
)

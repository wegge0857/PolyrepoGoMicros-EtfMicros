package server

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry" // 必须导入这个接口包
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
)

// 建议传入你项目生成的 conf 结构体
func NewRegistrar() registry.Registrar {
	cfg := api.DefaultConfig()

	cfg.Address = "127.0.0.1:8500"

	c, err := api.NewClient(cfg)
	if err != nil {
		panic(err) // 注册中心连接不上通常选择初始化时 panic
	}

	// 返回接口类型 registry.Registrar
	return consul.New(c)
}

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewRegistrar)

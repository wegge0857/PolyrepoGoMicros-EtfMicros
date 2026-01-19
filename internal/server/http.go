package server

import (
	"etfMicros/internal/conf"
	"etfMicros/internal/service"

	etfV1 "github.com/wegge0857/PolyrepoGoMicros-ApiLink/etf/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	etf *service.EtfService, //新增：注册Etf 服务--->@todo 这里后期只进行grpc访问，要去掉http的访问方式
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	etfV1.RegisterEtfHTTPServer(srv, etf) //新增：注册你的 Etf 服务
	return srv
}

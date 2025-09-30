package svc

import (
	"api-gateway/internal/config"
	orderpb "microservices-order-proto/order"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	OrderRpc orderpb.OrderServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		OrderRpc: orderpb.NewOrderServiceClient(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: c.OrderRpc.Endpoints,
			Timeout:   c.OrderRpc.Timeout,
		}).Conn()),
	}
}





package svc

import (
	"api-gateway/internal/config"
	goods "api-gateway/proto/goods"

	"github.com/zeromicro/go-zero/zrpc"
)

// ServiceContext 负责管理服务运行期所需的依赖：
// - 配置对象
// - gRPC 客户端（连接到商品 RPC 服务）
type ServiceContext struct {
	Config   config.Config
	GoodsRpc goods.GoodsServiceClient
}

// NewServiceContext 初始化并返回服务上下文对象
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		GoodsRpc: goods.NewGoodsServiceClient(zrpc.MustNewClient(c.GoodsRpc).Conn()),
	}
}

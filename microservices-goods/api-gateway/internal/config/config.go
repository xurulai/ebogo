package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	// HTTP 网关的基础配置（端口、超时、日志等），由 go-zero 提供的 RestConf 定义
	rest.RestConf
	// 商品 RPC 服务的客户端配置（直连或注册中心地址、超时等）
	GoodsRpc zrpc.RpcClientConf `json:",optional"`
}

type ServiceConfig struct {
	// 其他下游 HTTP 服务示例（如果未来接入可使用）
	GoodsApi ServiceEndpoint `json:",optional"`
	UserApi  ServiceEndpoint `json:",optional"`
	OrderApi ServiceEndpoint `json:",optional"`
}

type ServiceEndpoint struct {
	// 服务地址，例如 http://localhost:8888
	Endpoint string `json:",default=http://localhost:8888"`
	// 请求超时时间（毫秒）
	Timeout int `json:",default=5000"`
}

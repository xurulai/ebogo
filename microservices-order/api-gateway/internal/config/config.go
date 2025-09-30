package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	// 订单 RPC 服务配置
	OrderRpc struct {
		Endpoints []string // 订单服务端点列表
		Timeout   int64    // 超时时间（毫秒）
	}
}





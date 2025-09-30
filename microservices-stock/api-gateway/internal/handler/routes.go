package handler

import (
	"net/http"

	"api-gateway/internal/handler/gateway"
	"api-gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			// 库存管理相关接口
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/stock/set",
				Handler: gateway.SetStockHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/v1/stock/get",
				Handler: gateway.GetStockHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/stock/reduce",
				Handler: gateway.ReduceStockHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/stock/rollback",
				Handler: gateway.RollbackStockHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/stock/batch/get",
				Handler: gateway.BatchGetStockHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/stock/batch/reduce",
				Handler: gateway.BatchReduceStockHandler(serverCtx),
			},
		},
	)
}

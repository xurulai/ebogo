package gateway

import (
	"context"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	stockpb "api-gateway/proto/stock"

	"github.com/zeromicro/go-zero/core/logx"
)

type RollbackStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRollbackStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RollbackStockLogic {
	return &RollbackStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RollbackStockLogic) RollbackStock(req *types.RollbackStockRequest) (resp *types.CommonResponse, err error) {
	// 组装 gRPC 请求
	stockReq := &stockpb.RollBackStockInfo{
		GoodsId:     req.GoodsId,
		RollbackNum: req.RollbackNum,
		OrderId:     req.OrderId,
	}

	// 调用 RPC 服务
	stockResp, err := l.svcCtx.StockRpc.RollbackStock(l.ctx, stockReq)
	if err != nil {
		l.Errorf("RollbackStock RPC call failed: %v", err)
		return &types.CommonResponse{
			Success: false,
			Message: "库存回滚失败",
		}, nil
	}

	// 返回响应
	return &types.CommonResponse{
		Success: stockResp.Success,
		Message: stockResp.Message,
	}, nil
}





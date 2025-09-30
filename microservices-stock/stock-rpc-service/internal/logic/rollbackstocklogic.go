package logic

import (
	"context"
	stockpb "stock-rpc-service/proto/stock"
	"stock-rpc-service/internal/biz"
	"stock-rpc-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RollbackStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRollbackStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RollbackStockLogic {
	return &RollbackStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RollbackStock 回滚库存
func (l *RollbackStockLogic) RollbackStock(in *stockpb.RollBackStockInfo) (*stockpb.Response, error) {
	// 参数校验
	if in.GetGoodsId() <= 0 || in.GetRollbackNum() <= 0 || in.GetOrderId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "无效的参数")
	}

	// 创建业务逻辑实例
	stockBiz := biz.NewStockBiz(l.svcCtx)

	// 调用业务逻辑回滚库存
	err := stockBiz.RollbackStock(l.ctx, in.GetGoodsId(), in.GetRollbackNum(), in.GetOrderId())
	if err != nil {
		l.Errorf("RollbackStock failed: %v", err)
		return &stockpb.Response{
			Success: false,
			Message: "库存回滚失败",
		}, nil
	}

	// 返回成功响应
	return &stockpb.Response{
		Success: true,
		Message: "库存回滚成功",
	}, nil
}





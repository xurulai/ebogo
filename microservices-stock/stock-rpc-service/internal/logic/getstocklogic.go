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

type GetStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStockLogic {
	return &GetStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetStock 获取库存
func (l *GetStockLogic) GetStock(in *stockpb.GetStockRequest) (*stockpb.GoodsStockInfo, error) {
	// 参数校验
	if in.GetGoodsId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "无效的商品 ID")
	}

	// 创建业务逻辑实例
	stockBiz := biz.NewStockBiz(l.svcCtx)

	// 调用业务逻辑获取库存
	data, err := stockBiz.GetStockByGoodsId(l.ctx, in.GetGoodsId())
	if err != nil {
		l.Errorf("GetStock failed: %v", err)
		return nil, status.Errorf(codes.Internal, "获取库存失败: %v", err)
	}

	return data, nil
}

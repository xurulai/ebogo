package logic

import (
	"context"

	"goods-rpc-service/internal/biz"
	"goods-rpc-service/internal/svc"
	goodspb "goods-rpc-service/proto/goods"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateGoodsDetailLogic 处理更新商品详情的业务逻辑：
// - 校验参数（商品ID、价格）
// - 调用业务层执行更新
// - 返回成功/失败结果
type UpdateGoodsDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGoodsDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGoodsDetailLogic {
	return &UpdateGoodsDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateGoodsDetailLogic) UpdateGoodsDetail(in *goodspb.UpdateGoodsDetailRequest) (*goodspb.UpdateGoodsDetailResponse, error) {
	// 参数验证
	if in.GoodsId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "goods_id must be greater than 0")
	}
	if in.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be greater than 0")
	}

	// 调用业务逻辑
	bizGoods := biz.NewGoodsBiz(l.svcCtx)
	err := bizGoods.UpdateGoodsDetail(l.ctx, in.GoodsId, in.Price)
	if err != nil {
		l.Logger.Errorf("UpdateGoodsDetail failed: %v", err)
		return &goodspb.UpdateGoodsDetailResponse{
			Success: false,
			Message: "update failed",
		}, nil
	}

	return &goodspb.UpdateGoodsDetailResponse{
		Success: true,
		Message: "商品价格更新成功",
	}, nil
}

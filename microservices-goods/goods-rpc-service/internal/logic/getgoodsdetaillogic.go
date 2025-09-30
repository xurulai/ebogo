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

// GetGoodsDetailLogic 处理获取商品详情的业务逻辑：
// - 校验参数
// - 调用业务层获取详情
// - 统一错误码
type GetGoodsDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsDetailLogic {
	return &GetGoodsDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGoodsDetailLogic) GetGoodsDetail(in *goodspb.GetGoodsDetailRequest) (*goodspb.GetGoodsDetailResponse, error) {
	// 参数验证
	if in.GoodsId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "goods_id must be greater than 0")
	}

	// 调用业务逻辑
	bizGoods := biz.NewGoodsBiz(l.svcCtx)
	result, err := bizGoods.GetGoodsDetail(l.ctx, in.GoodsId)
	if err != nil {
		l.Logger.Errorf("GetGoodsDetail failed: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return result, nil
}

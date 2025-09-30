package server

import (
	"context"

	"goods-rpc-service/internal/logic"
	"goods-rpc-service/internal/svc"
	goodspb "goods-rpc-service/proto/goods"
)

// GoodsServer 实现 gRPC 定义的 GoodsService 接口
type GoodsServer struct {
	svcCtx *svc.ServiceContext
	goodspb.UnimplementedGoodsServiceServer
}

// NewGoodsServer 构造函数，注入服务上下文（数据库、缓存、配置等可从此扩展）
func NewGoodsServer(svcCtx *svc.ServiceContext) *GoodsServer {
	return &GoodsServer{
		svcCtx: svcCtx,
	}
}

// GetGoodsByRoom 获取直播间商品列表
func (s *GoodsServer) GetGoodsByRoom(ctx context.Context, in *goodspb.GetGoodsByRoomRequest) (*goodspb.GetGoodsByRoomResponse, error) {
	l := logic.NewGetGoodsByRoomLogic(ctx, s.svcCtx)
	return l.GetGoodsByRoom(in)
}

// GetGoodsDetail 获取商品详情
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, in *goodspb.GetGoodsDetailRequest) (*goodspb.GetGoodsDetailResponse, error) {
	l := logic.NewGetGoodsDetailLogic(ctx, s.svcCtx)
	return l.GetGoodsDetail(in)
}

// UpdateGoodsDetail 更新商品详情
func (s *GoodsServer) UpdateGoodsDetail(ctx context.Context, in *goodspb.UpdateGoodsDetailRequest) (*goodspb.UpdateGoodsDetailResponse, error) {
	l := logic.NewUpdateGoodsDetailLogic(ctx, s.svcCtx)
	return l.UpdateGoodsDetail(in)
}

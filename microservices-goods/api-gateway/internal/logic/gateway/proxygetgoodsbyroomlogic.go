package gateway

import (
	"context"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	goods "api-gateway/proto/goods"

	"github.com/zeromicro/go-zero/core/logx"
)

// ProxyGetGoodsByRoomLogic 网关层逻辑：代理调用商品 RPC 获取直播间商品列表
// - 负责 HTTP 入参到 RPC 入参的映射
// - 负责 RPC 出参到 HTTP 出参的映射
type ProxyGetGoodsByRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProxyGetGoodsByRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyGetGoodsByRoomLogic {
	return &ProxyGetGoodsByRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyGetGoodsByRoomLogic) ProxyGetGoodsByRoom(req *types.GetGoodsByRoomReq) (resp *types.GoodsListResp, err error) {
	// 组装 RPC 请求
	rpcReq := &goods.GetGoodsByRoomRequest{
		UserId: req.UserId,
		RoomId: req.RoomId,
	}

	// 调用 gRPC 服务
	rpcResp, err := l.svcCtx.GoodsRpc.GetGoodsByRoom(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("Failed to call goods RPC service: %v", err)
		return nil, err
	}

	// 映射 RPC 响应到 HTTP 响应
	var goodsList []*types.GoodsInfo
	for _, item := range rpcResp.Data {
		goodsList = append(goodsList, &types.GoodsInfo{
			GoodsId:     item.GoodsId,
			CategoryId:  item.CategoryId,
			Status:      item.Status,
			Title:       item.Title,
			MarketPrice: item.MarketPrice,
			Price:       item.Price,
			Brief:       item.Brief,
		})
	}

	return &types.GoodsListResp{
		CurrentGoodsId: rpcResp.CurrentGoodsId,
		Data:           goodsList,
	}, nil
}

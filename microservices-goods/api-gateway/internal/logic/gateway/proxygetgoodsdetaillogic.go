package gateway

import (
	"context"
	"time"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	goods "api-gateway/proto/goods"

	"github.com/zeromicro/go-zero/core/logx"
)

// ProxyGetGoodsDetailLogic 网关层逻辑：代理调用商品 RPC 获取商品详情
type ProxyGetGoodsDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProxyGetGoodsDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyGetGoodsDetailLogic {
	return &ProxyGetGoodsDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyGetGoodsDetailLogic) ProxyGetGoodsDetail(req *types.GetGoodsDetailReq) (resp *types.GoodsDetail, err error) {
	// 调用gRPC服务
	rpcReq := &goods.GetGoodsDetailRequest{
		GoodsId: req.GoodsId,
		UserId:  req.UserId,
	}

	start := time.Now()
	l.Infof("[GetGoodsDetail] calling goods-rpc, goods_id=%d user_id=%d", req.GoodsId, req.UserId)
	rpcResp, err := l.svcCtx.GoodsRpc.GetGoodsDetail(l.ctx, rpcReq)
	cost := time.Since(start)
	if err != nil {
		l.Errorf("[GetGoodsDetail] rpc error: %v (cost=%s)", err, cost)
		return nil, err
	}
	l.Infof("[GetGoodsDetail] rpc ok (cost=%s)", cost)

	// 转换响应格式
	return &types.GoodsDetail{
		GoodsId:     rpcResp.Goods.GoodsId,
		CategoryId:  rpcResp.Goods.CategoryId,
		Status:      rpcResp.Goods.Status,
		Title:       rpcResp.Goods.Title,
		Code:        rpcResp.Goods.Code,
		BrandName:   rpcResp.Goods.BrandName,
		MarketPrice: rpcResp.Goods.MarketPrice,
		Price:       rpcResp.Goods.Price,
		Brief:       rpcResp.Goods.Brief,
	}, nil
}

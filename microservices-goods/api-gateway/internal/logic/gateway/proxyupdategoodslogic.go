package gateway

import (
	"context"

	"api-gateway/internal/svc"
	"api-gateway/internal/types"
	goods "api-gateway/proto/goods"

	"github.com/zeromicro/go-zero/core/logx"
)

// ProxyUpdateGoodsLogic 网关层逻辑：代理调用商品 RPC 更新商品详情（价格）
type ProxyUpdateGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProxyUpdateGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyUpdateGoodsLogic {
	return &ProxyUpdateGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyUpdateGoodsLogic) ProxyUpdateGoods(req *types.UpdateGoodsDetailReq) (resp *types.Response, err error) {
	// 调用gRPC服务
	rpcReq := &goods.UpdateGoodsDetailRequest{
		GoodsId: req.GoodsId,
		Price:   req.Price,
	}

	rpcResp, err := l.svcCtx.GoodsRpc.UpdateGoodsDetail(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("Failed to call goods RPC service: %v", err)
		return &types.Response{
			Success: false,
			Message: "internal server error",
		}, nil
	}

	// 转换响应格式
	return &types.Response{
		Success: rpcResp.Success,
		Message: rpcResp.Message,
	}, nil
}

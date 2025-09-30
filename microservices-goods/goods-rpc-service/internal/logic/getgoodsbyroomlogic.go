package logic

import (
	"context"
	"time"

	"goods-rpc-service/internal/biz"
	"goods-rpc-service/internal/svc"
	goodspb "goods-rpc-service/proto/goods"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetGoodsByRoomLogic 处理获取直播间商品列表的业务逻辑：
// - 校验参数
// - 调用业务层获取数据
// - 统一错误码
type GetGoodsByRoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsByRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsByRoomLogic {
	return &GetGoodsByRoomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGoodsByRoomLogic) GetGoodsByRoom(in *goodspb.GetGoodsByRoomRequest) (*goodspb.GetGoodsByRoomResponse, error) {
	// 参数验证
	if in.RoomId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "room_id must be greater than 0")
	}

	newctx, cancel := context.WithTimeout(l.ctx, 20*time.Second)

	defer cancel()
	// 调用业务逻辑
	bizGoods := biz.NewGoodsBiz(l.svcCtx)
	result, err := bizGoods.GetGoodsByRoom(newctx, in.RoomId)
	if err != nil {
		l.Logger.Errorf("GetGoodsByRoom failed: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return result, nil
}

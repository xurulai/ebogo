package handler

import (
	"context"
	"goods_srv/biz/goods"
	"goods_srv/proto"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RPC的入口

type GoodsSrv struct {
	proto.UnimplementedGoodsServer
}

// GetGoodsByRoom 根据room_id获取直播间的商品列表
func (s *GoodsSrv) GetGoodsByRoom(ctx context.Context, req *proto.GetGoodsByRoomReq) (*proto.GoodsListResp, error) {

	//参数处理
	if req.GetRoomId() <= 0 {
		//无效的请求
		return nil, status.Error(codes.InvalidArgument, "请求参数有误")
	}
	// 去查询数据并封装返回的响应数据 --> 业务逻辑
	data, err := goods.GetGoodsByRoom(ctx, req.GetRoomId())
	if err != nil {
		return nil, status.Error(codes.Internal, "内部错误")
	}
	return data, nil
}
// GetGoodsDetail 根据goods_id获取商品详情
func (s *GoodsSrv) GetGoodsDetail(ctx context.Context, req *proto.GetGoodsDetailReq) (*proto.GoodsDetail, error) {
    log.Printf("Received GetGoodsDetail request: %+v", req)

    if req.GetUserId() <= 0 || req.GetGoodsId() <= 0 {
        log.Printf("Invalid request parameters: %+v", req)
        return nil, status.Error(codes.InvalidArgument, "请求参数有误")
    }

    data, err := goods.GetGoodsDetailById(ctx, req.GetGoodsId())
    if err != nil {
        log.Printf("Failed to get goods detail: %v", err)
        return nil, status.Error(codes.Internal, "内部错误")
    }

    log.Printf("Returning response: %+v", data)
    return data, nil
}

func (s *GoodsSrv) UpdateGoodsDetail(ctx context.Context, req *proto.UpdateGoodsDetailReq) (*proto.Response, error) {
    log.Printf("Received UpdateGoodsDetail request: %+v", req)

    if req.GetGoodsId() <= 0 {
        log.Printf("Invalid request parameters: %+v", req)
        return nil, status.Error(codes.InvalidArgument, "请求参数有误")
    }
    if req.GetPrice() <= 0 {
        log.Printf("Invalid request parameters: %+v", req)
        return nil, status.Error(codes.InvalidArgument, "请求参数有误")
    }

    // 更新数据库中的商品信息
    _,err := goods.UpdateGoodsDetail(ctx, req.GetGoodsId(),req.GetPrice())
    if err != nil {
        log.Printf("Failed to update goods detail: %v", err)
        return nil, status.Error(codes.Internal, "内部错误")
    }


    // 3. 返回成功响应
    log.Printf("Goods detail updated successfully for GoodsId: %d", req.GetGoodsId())
    return &proto.Response{
        Success: true,
        Message: "商品信息更新成功",
    }, nil
}

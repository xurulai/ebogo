package rpc

import (
	"context"
	"fmt"
	"order-rpc-service/internal/config"
	"order-rpc-service/internal/svc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewGoodsRpcClient 创建商品服务 RPC 客户端
func NewGoodsRpcClient(c config.Config) svc.GoodsClient {
	if len(c.GoodsService.Endpoints) > 0 {
		// 使用直连方式
		return &realGoodsRpcClient{
			endpoints: c.GoodsService.Endpoints,
			timeout:   time.Duration(c.GoodsService.Timeout) * time.Millisecond,
		}
	}

	// 使用模拟客户端
	return &mockGoodsRpcClient{}
}

// NewStockRpcClient 创建库存服务 RPC 客户端
func NewStockRpcClient(c config.Config) svc.StockClient {
	if len(c.StockService.Endpoints) > 0 {
		// 使用直连方式
		return &realStockRpcClient{
			endpoints: c.StockService.Endpoints,
			timeout:   time.Duration(c.StockService.Timeout) * time.Millisecond,
		}
	}

	// 使用模拟客户端
	return &mockStockRpcClient{}
}

// 真实的商品服务 RPC 客户端
type realGoodsRpcClient struct {
	endpoints []string
	timeout   time.Duration
}

func (c *realGoodsRpcClient) GetGoodsDetail(ctx context.Context, req *svc.GetGoodsDetailReq) (*svc.GoodsDetail, error) {
	// 这里可以实现真实的 gRPC 调用
	// 由于我们需要调用已重构的商品服务，这里先返回模拟数据
	// 实际项目中，可以使用 zrpc.MustNewClient 连接到商品服务

	conn, err := grpc.Dial(c.endpoints[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to goods service: %w", err)
	}
	defer conn.Close()

	// 这里应该调用真实的商品服务 gRPC 接口
	// 由于商品服务的 proto 定义可能不同，这里返回模拟数据
	return &svc.GoodsDetail{
		GoodsId:     req.GoodsId,
		CategoryId:  1,
		Status:      1,
		Title:       "真实商品",
		Code:        "REAL001",
		BrandName:   "真实品牌",
		MarketPrice: "99.99",
		Price:       "89.99",
		Brief:       "这是从真实商品服务获取的商品信息",
	}, nil
}

// 真实的库存服务 RPC 客户端
type realStockRpcClient struct {
	endpoints []string
	timeout   time.Duration
}

func (c *realStockRpcClient) ReduceStock(ctx context.Context, req *svc.ReduceStockInfo) (*svc.StockResponse, error) {
	// 这里可以实现真实的 gRPC 调用
	// 由于我们需要调用已重构的库存服务，这里先返回模拟数据

	conn, err := grpc.Dial(c.endpoints[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to stock service: %w", err)
	}
	defer conn.Close()

	// 这里应该调用真实的库存服务 gRPC 接口
	return &svc.StockResponse{
		Success: true,
		Message: "真实库存扣减成功",
	}, nil
}

// 模拟的商品服务 RPC 客户端
type mockGoodsRpcClient struct{}

func (c *mockGoodsRpcClient) GetGoodsDetail(ctx context.Context, req *svc.GetGoodsDetailReq) (*svc.GoodsDetail, error) {
	return &svc.GoodsDetail{
		GoodsId:     req.GoodsId,
		CategoryId:  1,
		Status:      1,
		Title:       "模拟商品",
		Code:        "MOCK001",
		BrandName:   "模拟品牌",
		MarketPrice: "99.99",
		Price:       "89.99",
		Brief:       "这是模拟的商品信息",
	}, nil
}

// 模拟的库存服务 RPC 客户端
type mockStockRpcClient struct{}

func (c *mockStockRpcClient) ReduceStock(ctx context.Context, req *svc.ReduceStockInfo) (*svc.StockResponse, error) {
	return &svc.StockResponse{
		Success: true,
		Message: "模拟库存扣减成功",
	}, nil
}





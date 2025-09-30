package server

import (
	"context"
	"stock-rpc-service/internal/logic"
	"stock-rpc-service/internal/svc"
	stockpb "stock-rpc-service/proto/stock"
)

type StockServer struct {
	svcCtx *svc.ServiceContext
	stockpb.UnimplementedStockServiceServer
}

func NewStockServer(svcCtx *svc.ServiceContext) *StockServer {
	return &StockServer{
		svcCtx: svcCtx,
	}
}

// SetStock 设置库存
func (s *StockServer) SetStock(ctx context.Context, in *stockpb.GoodsStockInfo) (*stockpb.Response, error) {
	l := logic.NewSetStockLogic(ctx, s.svcCtx)
	return l.SetStock(in)
}

// GetStock 获取库存
func (s *StockServer) GetStock(ctx context.Context, in *stockpb.GetStockRequest) (*stockpb.GoodsStockInfo, error) {
	l := logic.NewGetStockLogic(ctx, s.svcCtx)
	return l.GetStock(in)
}

// ReduceStock 减少库存
func (s *StockServer) ReduceStock(ctx context.Context, in *stockpb.ReduceStockInfo) (*stockpb.Response, error) {
	l := logic.NewReduceStockLogic(ctx, s.svcCtx)
	return l.ReduceStock(in)
}

// RollbackStock 回滚库存
func (s *StockServer) RollbackStock(ctx context.Context, in *stockpb.RollBackStockInfo) (*stockpb.Response, error) {
	l := logic.NewRollbackStockLogic(ctx, s.svcCtx)
	return l.RollbackStock(in)
}

// BatchGetStock 批量获取库存
func (s *StockServer) BatchGetStock(ctx context.Context, in *stockpb.StockInfoList) (*stockpb.StockInfoList, error) {
	l := logic.NewBatchGetStockLogic(ctx, s.svcCtx)
	return l.BatchGetStock(in)
}

// BatchReduceStock 批量减少库存
func (s *StockServer) BatchReduceStock(ctx context.Context, in *stockpb.StockInfoList) (*stockpb.Response, error) {
	l := logic.NewBatchReduceStockLogic(ctx, s.svcCtx)
	return l.BatchReduceStock(in)
}

package model

import "github.com/google/uuid"

type SeckillRequest struct {
	ID       string `json:"id"`       // 唯一标识符，用于幂等性检查
	UserId   int64  `json:"user_id"`  // 用户ID
	GoodsId  int64  `json:"goods_id"` // 商品ID
	Quantity int    `json:"quantity"` // 秒杀数量
}

// NewSeckillRequest 创建一个新的 SeckillRequest 实例，并生成一个唯一的 ID
func NewSeckillRequest(userId int64, goodsId int64, quantity int) *SeckillRequest {
	return &SeckillRequest{
		ID:       uuid.New().String(), // 使用 UUID 生成唯一的请求ID
		UserId:   userId,
		GoodsId:  goodsId,
		Quantity: quantity,
	}
}
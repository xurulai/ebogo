package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	orderpb "microservices-order-proto/order"
	"order-rpc-service/internal/svc"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// OrderBiz 订单业务逻辑
type OrderBiz struct {
	svcCtx *svc.ServiceContext
}

// NewOrderBiz 创建订单业务实例
func NewOrderBiz(svcCtx *svc.ServiceContext) *OrderBiz {
	return &OrderBiz{
		svcCtx: svcCtx,
	}
}

// Order 订单模型
type Order struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	OrderId        int64     `gorm:"column:order_id;uniqueIndex;not null"`
	UserId         int64     `gorm:"column:user_id;not null"`
	PayAmount      int64     `gorm:"column:pay_amount;not null;default:0"`
	Status         string    `gorm:"column:status;not null;default:'pending'"`
	ReceiveAddress string    `gorm:"column:receive_address;not null"`
	ReceiveName    string    `gorm:"column:receive_name;not null"`
	ReceivePhone   string    `gorm:"column:receive_phone;not null"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

// TableName 声明表名
func (Order) TableName() string {
	return "xx_order"
}

// OrderDetail 订单详情模型
type OrderDetail struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	OrderId   int64     `gorm:"column:order_id;not null"`
	UserId    int64     `gorm:"column:user_id;not null"`
	GoodsId   int64     `gorm:"column:goods_id;not null"`
	Num       int64     `gorm:"column:num;not null"`
	Title     string    `gorm:"column:title;not null"`
	Status    string    `gorm:"column:status;not null;default:'pending'"`
	Price     int64     `gorm:"column:price;not null"`
	Brief     string    `gorm:"column:brief"`
	PayAmount int64     `gorm:"column:pay_amount;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName 声明表名
func (OrderDetail) TableName() string {
	return "xx_order_detail"
}

// OrderEntity 事务消息实体，保持原有的事务消息逻辑
type OrderEntity struct {
	OrderId    int64                       // 订单ID
	Param      *orderpb.CreateOrderRequest // 订单请求参数
	Topic      string                      // 事务消息的主题
	RetryCount int64                       // 重试次数
	err        error                       // 本地事务执行过程中可能产生的错误
	svcCtx     *svc.ServiceContext         // 服务上下文
}

// CreateOrder 创建订单 - 保持原有的核心业务逻辑
func (b *OrderBiz) CreateOrder(ctx context.Context, param *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	// 1. 生成订单号
	orderId := b.svcCtx.SnowflakeNode.Generate().Int64()

	// 创建OrderEntity实例，用于事务消息的上下文
	orderEntity := &OrderEntity{
		OrderId: orderId,
		Param:   param,
		Topic:   b.svcCtx.Config.RocketMQ.Topic.CreateOrder, // 默认Topic为创建订单
		svcCtx:  b.svcCtx,
	}

	// 创建事务生产者
	p, err := rocketmq.NewTransactionProducer(
		orderEntity, // 将 OrderEntity 作为事务消息的上下文
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{b.svcCtx.Config.RocketMQ.NameServer})),
		producer.WithRetry(3),
		producer.WithGroupName(b.svcCtx.Config.RocketMQ.GroupName),
	)
	if err != nil {
		zap.L().Error("NewTransactionProducer failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "NewTransactionProducer failed")
	}

	// 启动事务生产者
	p.Start()
	defer p.Shutdown()

	// 构造事务消息的内容
	data := OrderDetail{
		OrderId: orderId,
		GoodsId: param.GoodsId,
		Num:     int64(param.Num),
	}

	msgBody, _ := json.Marshal(data)
	// 构造消息
	msg := &primitive.Message{
		Topic: orderEntity.Topic,
		Body:  msgBody,
	}

	// 发送事务消息
	res, err := p.SendMessageInTransaction(ctx, msg)
	if err != nil {
		zap.L().Error("SendMessageInTransaction failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "create order failed")
	}

	// 根据事务消息的响应状态和Topic判断订单创建是否成功
	if res.State == primitive.CommitMessageState {
		if orderEntity.Topic == b.svcCtx.Config.RocketMQ.Topic.CreateOrderSuccessfully {
			return &orderpb.CreateOrderResponse{
				Success: true,
				Message: "Order created successfully",
				OrderId: orderId,
			}, nil
		} else if orderEntity.Topic == b.svcCtx.Config.RocketMQ.Topic.PayTimeout {
			return nil, status.Error(codes.Internal, "Order creation failed due to timeout")
		} else if orderEntity.Topic == b.svcCtx.Config.RocketMQ.Topic.StockRollback {
			return nil, status.Error(codes.Internal, "Order creation failed")
		}
	}

	if res.State == primitive.RollbackMessageState {
		return nil, status.Error(codes.Internal, "create order failed")
	}

	return nil, status.Error(codes.Internal, "unknown transaction state")
}

// ExecuteLocalTransaction 是 RocketMQ 事务消息的本地事务执行逻辑 - 保持原有逻辑
func (o *OrderEntity) ExecuteLocalTransaction(*primitive.Message) primitive.LocalTransactionState {
	fmt.Println("in ExecuteLocalTransaction...")

	// 参数校验
	if o.Param == nil {
		zap.L().Error("ExecuteLocalTransaction param is nil")
		o.err = status.Error(codes.Internal, "invalid OrderEntity")
		return primitive.RollbackMessageState
	}

	param := o.Param
	ctx := context.Background()

	// 1. 查询商品金额（营销）--> RPC连接 goods_service
	// 调用 goods_service 获取商品详情，包括价格。
	goodsDetail, err := o.svcCtx.GoodsRpc.GetGoodsDetail(ctx, &svc.GetGoodsDetailReq{
		GoodsId: param.GoodsId,
		UserId:  param.UserId,
	})
	if err != nil {
		// 如果查询商品失败，记录日志并返回 Rollback 状态，表示本地事务失败。
		zap.L().Error("GoodsRpc.GetGoodsDetail failed", zap.Error(err))
		o.err = status.Error(codes.Internal, err.Error())
		return primitive.RollbackMessageState
	}

	// 将价格字符串转换为整型（分为单位）
	payAmountStr := goodsDetail.Price
	payAmount, err := strconv.ParseInt(payAmountStr, 10, 64)
	if err != nil {
		zap.L().Error("Failed to parse price", zap.Error(err))
		o.err = status.Error(codes.Internal, "Invalid price format")
		return primitive.RollbackMessageState
	}
	// 转换为分为单位（假设价格是元，需要乘以100）
	payAmount = payAmount * 100

	// 2. 库存校验及扣减  --> RPC连接 stock_service
	// 调用 stock_service 扣减库存。
	_, err = o.svcCtx.StockRpc.ReduceStock(ctx, &svc.ReduceStockInfo{
		GoodsId: o.Param.GoodsId,
		Num:     int64(o.Param.Num),
		OrderId: o.OrderId,
	})

	if err != nil {
		// 如果库存扣减失败，记录日志并返回 Rollback 状态，表示本地事务失败。
		zap.L().Error("StockRpc.ReduceStock failed", zap.Error(err))
		o.err = status.Error(codes.Internal, "ReduceStock failed")
		return primitive.RollbackMessageState
	}

	// 3. 创建订单
	orderData := Order{
		OrderId:        o.OrderId,
		UserId:         param.UserId,
		PayAmount:      payAmount,
		ReceiveAddress: param.Address,
		ReceiveName:    param.Name,
		ReceivePhone:   param.Phone,
		Status:         "pending", // 待支付
	}

	orderDetail := OrderDetail{
		OrderId:   o.OrderId,
		UserId:    param.UserId,
		GoodsId:   param.GoodsId,
		Num:       int64(param.Num),
		Title:     goodsDetail.Title, // 从商品服务获取
		Status:    "pending",
		Price:     payAmount,
		Brief:     goodsDetail.Brief, // 从商品服务获取
		PayAmount: payAmount,
	}

	// 使用事务创建订单和订单详情记录
	err = o.svcCtx.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&orderData).Error; err != nil {
			return err
		}
		if err := tx.Create(&orderDetail).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		// 如果订单创建失败，发送库存回滚消息
		o.Topic = o.svcCtx.Config.RocketMQ.Topic.StockRollback
		msg := primitive.NewMessage(o.Topic, []byte(fmt.Sprintf(`{"orderId":%d,"reason":"order creation failed"}`, o.OrderId)))
		_, errSend := o.svcCtx.MQProducer.SendSync(context.Background(), msg)
		if errSend != nil {
			zap.L().Error("send order_failed msg failed", zap.Error(errSend))
		}
		zap.L().Error("CreateOrderWithTransaction failed", zap.Error(err))
		return primitive.RollbackMessageState
	}

	// 发送延迟消息，用于订单超时处理
	data := OrderDetail{
		OrderId: o.OrderId,
		GoodsId: param.GoodsId,
		Num:     int64(param.Num),
		Status:  "pending",
	}
	timeoutMsgBody, _ := json.Marshal(data)
	o.Topic = o.svcCtx.Config.RocketMQ.Topic.PayTimeout
	msgTimeout := primitive.NewMessage(o.Topic, timeoutMsgBody)
	msgTimeout.WithDelayTimeLevel(3) // 设置延迟级别（例如 10s）

	_, err = o.svcCtx.MQProducer.SendSync(context.Background(), msgTimeout)
	if err != nil {
		o.Topic = o.svcCtx.Config.RocketMQ.Topic.PayTimeout
		msg := primitive.NewMessage(o.Topic, []byte(fmt.Sprintf(`{"orderId":%d,"reason":"timeout message send failed"}`, o.OrderId)))
		_, errSend := o.svcCtx.MQProducer.SendSync(context.Background(), msg)
		if errSend != nil {
			zap.L().Error("send order_failed msg failed", zap.Error(errSend))
		}
		zap.L().Error("send delay msg failed", zap.Error(err))
		return primitive.RollbackMessageState
	}

	// 如果本地事务成功，发送订单创建成功消息
	o.Topic = o.svcCtx.Config.RocketMQ.Topic.CreateOrderSuccessfully
	msgSuccess := primitive.NewMessage(o.Topic, []byte(fmt.Sprintf(`{"orderId":%d,"status":"success"}`, o.OrderId)))
	_, err = o.svcCtx.MQProducer.SendSync(context.Background(), msgSuccess)
	if err != nil {
		zap.L().Error("send order success msg failed", zap.Error(err))
		return primitive.RollbackMessageState
	}

	return primitive.CommitMessageState
}

// CheckLocalTransaction 是 RocketMQ 事务消息的状态回查逻辑 - 保持原有逻辑
func (o *OrderEntity) CheckLocalTransaction(*primitive.MessageExt) primitive.LocalTransactionState {
	// 查询订单是否创建成功
	var order Order
	err := o.svcCtx.DB.Where("order_id = ?", o.OrderId).First(&order).Error
	if err == gorm.ErrRecordNotFound {
		// 如果订单未创建成功，返回 Commit 状态，表示需要回滚库存
		return primitive.CommitMessageState
	}
	// 如果订单已创建成功，返回 Rollback 状态，表示不需要回滚库存
	return primitive.RollbackMessageState
}

// GetOrderList 获取订单列表
func (b *OrderBiz) GetOrderList(ctx context.Context, userId int64, pageNum, pageSize int32) (*orderpb.GetOrderListResponse, error) {
	var orders []Order
	var total int64

	// 查询总数
	err := b.svcCtx.DB.WithContext(ctx).Model(&Order{}).Where("user_id = ?", userId).Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("查询订单总数失败: %w", err)
	}

	// 分页查询订单
	offset := (pageNum - 1) * pageSize
	err = b.svcCtx.DB.WithContext(ctx).
		Where("user_id = ?", userId).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, fmt.Errorf("查询订单列表失败: %w", err)
	}

	// 转换为响应格式
	var orderInfos []*orderpb.OrderInfo
	for _, order := range orders {
		// 查询订单详情
		var orderDetails []OrderDetail
		b.svcCtx.DB.WithContext(ctx).Where("order_id = ?", order.OrderId).Find(&orderDetails)

		var detailInfos []*orderpb.OrderDetailInfo
		for _, detail := range orderDetails {
			detailInfos = append(detailInfos, &orderpb.OrderDetailInfo{
				OrderId:   detail.OrderId,
				GoodsId:   detail.GoodsId,
				Num:       detail.Num,
				Title:     detail.Title,
				Status:    detail.Status,
				Price:     detail.Price,
				Brief:     detail.Brief,
				PayAmount: detail.PayAmount,
			})
		}

		orderInfos = append(orderInfos, &orderpb.OrderInfo{
			OrderId:        order.OrderId,
			UserId:         order.UserId,
			Status:         order.Status,
			PayAmount:      order.PayAmount,
			ReceiveAddress: order.ReceiveAddress,
			ReceiveName:    order.ReceiveName,
			ReceivePhone:   order.ReceivePhone,
			CreatedAt:      order.CreatedAt.Unix(),
			OrderDetails:   detailInfos,
		})
	}

	return &orderpb.GetOrderListResponse{
		Total:  int32(total),
		Orders: orderInfos,
	}, nil
}

// GetOrderDetail 获取订单详情
func (b *OrderBiz) GetOrderDetail(ctx context.Context, orderId, userId int64) (*orderpb.GetOrderDetailResponse, error) {
	var order Order
	err := b.svcCtx.DB.WithContext(ctx).Where("order_id = ? AND user_id = ?", orderId, userId).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	// 查询订单详情
	var orderDetails []OrderDetail
	err = b.svcCtx.DB.WithContext(ctx).Where("order_id = ?", orderId).Find(&orderDetails).Error
	if err != nil {
		return nil, fmt.Errorf("查询订单详情失败: %w", err)
	}

	var detailInfos []*orderpb.OrderDetailInfo
	for _, detail := range orderDetails {
		detailInfos = append(detailInfos, &orderpb.OrderDetailInfo{
			OrderId:   detail.OrderId,
			GoodsId:   detail.GoodsId,
			Num:       detail.Num,
			Title:     detail.Title,
			Status:    detail.Status,
			Price:     detail.Price,
			Brief:     detail.Brief,
			PayAmount: detail.PayAmount,
		})
	}

	orderInfo := &orderpb.OrderInfo{
		OrderId:        order.OrderId,
		UserId:         order.UserId,
		Status:         order.Status,
		PayAmount:      order.PayAmount,
		ReceiveAddress: order.ReceiveAddress,
		ReceiveName:    order.ReceiveName,
		ReceivePhone:   order.ReceivePhone,
		CreatedAt:      order.CreatedAt.Unix(),
		OrderDetails:   detailInfos,
	}

	return &orderpb.GetOrderDetailResponse{
		Order: orderInfo,
	}, nil
}

// UpdateOrderStatus 更新订单状态
func (b *OrderBiz) UpdateOrderStatus(ctx context.Context, orderId int64, status string) error {
	result := b.svcCtx.DB.WithContext(ctx).
		Model(&Order{}).
		Where("order_id = ?", orderId).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("更新订单状态失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("订单不存在")
	}

	// 同时更新订单详情状态
	b.svcCtx.DB.WithContext(ctx).
		Model(&OrderDetail{}).
		Where("order_id = ?", orderId).
		Update("status", status)

	return nil
}

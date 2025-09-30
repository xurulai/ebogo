package main

import (
	"context"
	"fmt"
	"log"
	"time"

	orderpb "microservices-order-proto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial("localhost:9002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := orderpb.NewOrderServiceClient(conn)

	// 测试创建订单
	testCreateOrder(client)
}

func testCreateOrder(client orderpb.OrderServiceClient) {
	fmt.Println("🚀 开始测试订单创建...")

	// 测试用例 1: 正常订单
	fmt.Println("\n📝 测试用例 1: 创建正常订单")
	req1 := &orderpb.CreateOrderRequest{
		GoodsId: 1,
		Num:     1,
		UserId:  1001,
		Address: "北京市朝阳区测试地址123号",
		Name:    "张三",
		Phone:   "13800138000",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp1, err := client.CreateOrder(ctx, req1)
	if err != nil {
		fmt.Printf("❌ 创建订单失败: %v\n", err)
	} else {
		fmt.Printf("✅ 订单创建成功!\n")
		fmt.Printf("   订单ID: %d\n", resp1.OrderId)
		fmt.Printf("   价格: %s\n", resp1.Price)
		fmt.Printf("   消息: %s\n", resp1.Message)
	}

	// 测试用例 2: 多数量订单
	fmt.Println("\n📝 测试用例 2: 创建多数量订单")
	req2 := &orderpb.CreateOrderRequest{
		GoodsId: 2,
		Num:     3,
		UserId:  1002,
		Address: "上海市浦东新区测试路456号",
		Name:    "李四",
		Phone:   "13900139000",
	}

	resp2, err := client.CreateOrder(ctx, req2)
	if err != nil {
		fmt.Printf("❌ 创建订单失败: %v\n", err)
	} else {
		fmt.Printf("✅ 订单创建成功!\n")
		fmt.Printf("   订单ID: %d\n", resp2.OrderId)
		fmt.Printf("   价格: %s\n", resp2.Price)
		fmt.Printf("   消息: %s\n", resp2.Message)
	}

	// 测试用例 3: 参数验证（无效用户ID）
	fmt.Println("\n📝 测试用例 3: 参数验证 (无效用户ID)")
	req3 := &orderpb.CreateOrderRequest{
		GoodsId: 1,
		Num:     1,
		UserId:  0, // 无效用户ID
		Address: "测试地址",
		Name:    "测试用户",
		Phone:   "13800138000",
	}

	resp3, err := client.CreateOrder(ctx, req3)
	if err != nil {
		fmt.Printf("✅ 正确拒绝无效请求: %v\n", err)
	} else {
		fmt.Printf("⚠️ 意外接受了无效请求: %+v\n", resp3)
	}

	// 批量测试
	fmt.Println("\n📝 测试用例 4: 批量创建订单")
	successCount := 0
	totalCount := 5

	for i := 1; i <= totalCount; i++ {
		req := &orderpb.CreateOrderRequest{
			GoodsId: int64(i),
			Num:     int32(i%3 + 1),
			UserId:  int64(2000 + i),
			Address: fmt.Sprintf("批量测试地址%d号", i),
			Name:    fmt.Sprintf("批量测试用户%d", i),
			Phone:   fmt.Sprintf("1380013800%d", i),
		}

		resp, err := client.CreateOrder(ctx, req)
		if err != nil {
			fmt.Printf("❌ 订单 %d 创建失败: %v\n", i, err)
		} else {
			successCount++
			fmt.Printf("✅ 订单 %d 创建成功 (ID: %d)\n", i, resp.OrderId)
		}

		// 避免过快请求
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n🎯 批量测试结果: %d/%d 成功\n", successCount, totalCount)

	fmt.Println("\n🎉 所有测试完成!")
}




package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"user_service/proto" // 导入用户服务的 Protobuf 生成的 Go 包

	"google.golang.org/grpc"                      // 导入 gRPC 包
	"google.golang.org/grpc/credentials/insecure" // 导入用于不安全连接的 gRPC 凭证包
	"google.golang.org/grpc/metadata"
)

var (
	conn   *grpc.ClientConn        // gRPC 客户端连接
	client proto.UserServiceClient // gRPC 客户端对象，用于调用用户服务
)

// 初始化函数，用于建立 gRPC 连接
func init() {
	var err error
	conn, err = grpc.Dial(
		"127.0.0.1:8399", // gRPC 服务地址，根据实际情况修改
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 使用不安全的连接（仅用于测试环境）
	)
	if err != nil {
		panic(err) // 如果连接失败，直接 panic
	}
	client = proto.NewUserServiceClient(conn) // 创建 gRPC 客户端对象
}

// 测试用户注册接口的函数
func TestRegisterUser(wg *sync.WaitGroup, index int, errCount *int32) {
	defer wg.Done() // 在函数返回时通知 WaitGroup 当前协程已完成

	param := &proto.RegisterRequest{
		Username: "agiao3",
		Password: "sz12345678",
		Email:    "2535512844@qq.com",
		Phone:    "1380013800",
	}

	// 创建带有终端信息的上下文
	ctx :=metadata.AppendToOutgoingContext(context.Background(), "terminal", "pc")
	start := time.Now()                                       // 记录调用开始时间
	resp, err := client.Register(ctx, param) // 调用 gRPC 服务的用户注册接口
	duration := time.Since(start)                             // 计算调用耗时

	if err != nil {
		atomic.AddInt32(errCount, 1) // 如果发生错误，原子操作增加错误计数
		fmt.Printf("协程 %d: 调用用户注册接口失败: %v, 耗时: %v\n", index, err, duration)
	} else {
		fmt.Printf("协程 %d: 调用用户注册接口成功: %+v, 耗时: %v\n", index, resp, duration)
	}
}

// 测试用户登录接口的函数
func TestLoginUser(wg *sync.WaitGroup, index int, errCount *int32) {
	defer wg.Done() // 在函数返回时通知 WaitGroup 当前协程已完成

	// 定义登录请求参数
	param := &proto.LoginRequest{
		Username: "agiao3",
		Password: "sz12345678",
	}

	start := time.Now()                                    // 记录调用开始时间
	resp, err := client.Login(context.Background(), param) // 调用 gRPC 服务的用户登录接口
	duration := time.Since(start)                          // 计算调用耗时

	if err != nil {
		atomic.AddInt32(errCount, 1) // 如果发生错误，原子操作增加错误计数
		fmt.Printf("协程 %d: 调用用户登录接口失败: %v, 耗时: %v\n", index, err, duration)
	} else {
		fmt.Printf("协程 %d: 调用用户登录接口成功: %+v, 耗时: %v\n", index, resp, duration)
	}
}

// 测试发送验证码接口的函数
func TestSendSmsCode(wg *sync.WaitGroup, index int, errCount *int32) {
	defer wg.Done() // 在函数返回时通知 WaitGroup 当前协程已完成

	// 定义发送验证码请求参数
	param := &proto.SendSmsCodeRequest{
		Phone: "1380013800", // 测试手机号
	}

	start := time.Now()                                          // 记录调用开始时间
	resp, err := client.SendSmsCode(context.Background(), param) // 调用 gRPC 服务的发送验证码接口
	duration := time.Since(start)                                // 计算调用耗时

	if err != nil {
		atomic.AddInt32(errCount, 1) // 如果发生错误，原子操作增加错误计数
		fmt.Printf("协程 %d: 调用发送验证码接口失败: %v, 耗时: %v\n", index, err, duration)
	} else {
		fmt.Printf("协程 %d: 调用发送验证码接口成功: %+v, 耗时: %v\n", index, resp, duration)
	}
}

// 测试短信验证码登录接口的函数
func TestLoginBySms(wg *sync.WaitGroup, index int, errCount *int32) {
	defer wg.Done() // 在函数返回时通知 WaitGroup 当前协程已完成

	// 定义短信验证码登录请求参数
	param := &proto.LoginBySmsRequest{
		Phone:   "1380013800", // 测试手机号
		SmsCode: "123456789",     // 测试验证码
	}

	// 创建带有终端信息的上下文
	ctx := metadata.AppendToOutgoingContext(context.Background(), "terminal", "mobile")

	start := time.Now()                                         // 记录调用开始时间
	resp, err := client.LoginBySms(ctx, param) // 调用 gRPC 服务的短信验证码登录接口
	duration := time.Since(start)                               // 计算调用耗时

	if err != nil {
		atomic.AddInt32(errCount, 1) // 如果发生错误，原子操作增加错误计数
		fmt.Printf("协程 %d: 调用短信验证码登录接口失败: %v, 耗时: %v\n", index, err, duration)
	} else {
		fmt.Printf("协程 %d: 调用短信验证码登录接口成功: %+v, 耗时: %v\n", index, resp, duration)
	}
}

func main() {
	defer conn.Close()     // 程序退出时关闭 gRPC 连接
	var wg sync.WaitGroup  // 使用 WaitGroup 等待所有协程完成
	var errCount int32 = 0 // 初始化错误计数器

	// 并发调用测试接口
	for i := 0; i < 5; i++ { // 启动 5 组并发测试
		wg.Add(2) // 每组并发调用 1 个注册接口

		//go TestRegisterUser(&wg, i, &errCount) // 启动协程测试用户注册接口
		go TestLoginUser(&wg,i,&errCount)
		//go TestSendSmsCode(&wg,i,&errCount)
		go TestLoginBySms(&wg, i, &errCount)
	}
	wg.Wait()                                             // 等待所有协程完成
	fmt.Printf("总错误数: %d\n", atomic.LoadInt32(&errCount)) // 输出总错误数
}

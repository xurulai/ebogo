package handler

import (
	"context"
	"fmt"
	"time"
	"user_service/biz/user"
	"user_service/proto"
	"user_service/third_party/jwt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserSrv struct {
	proto.UnimplementedUserServiceServer // 嵌入未实现的 OrderServer 接口，用于兼容 gRPC 的接口实现
}

// 用户注册
// Register 用户注册
func (u *UserSrv) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	fmt.Println("in CreateOrder ... ") // 打印进入方法的日志

	// 参数处理
	if req.GetUsername() == "" || req.GetPassword() == "" || req.GetEmail() == "" || req.GetPhone() == "" {
		// 无效的请求
		return nil, status.Error(codes.InvalidArgument, "请求参数有误") // 返回 gRPC 的 InvalidArgument 错误
	}

	// 业务处理
	response, err := user.Register(ctx, req.Username, req.Password, req.Email,req.Phone)
	if err != nil {
		zap.L().Error("user.Register failed", zap.Error(err)) // 记录错误日志
		return response, status.Error(codes.Internal, "内部错误") // 返回 gRPC 的 Internal 错误
	}

	return response, nil // 返回空响应，表示操作成功
}

// Login 用户登录
func (u *UserSrv) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	fmt.Println("in Login ... ") // 打印进入方法的日志

	// 参数处理
	if req.GetUsername() == "" || req.GetPassword() == "" {
		// 无效的请求
		return nil, status.Error(codes.InvalidArgument, "请求参数有误") // 返回 gRPC 的 InvalidArgument 错误
	}

	// 业务处理
	response, err := user.Login(ctx, req.Username, req.Password)
	if err != nil {
		zap.L().Error("user.Login failed", zap.Error(err))    // 记录错误日志
		return response, status.Error(codes.Internal, "内部错误") // 返回 gRPC 的 Internal 错误
	}

	return response, nil // 返回空响应，表示操作成功
}

// RefreshToken 刷新Token
func (u *UserSrv) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {
	// 验证 RefreshToken 是否有效
	claims, err := jwt.ParseToken(req.RefreshToken)
	if err != nil {
		return &proto.RefreshTokenResponse{
			Success: false,
			Message: "刷新Token无效",
		}, status.Error(codes.Unauthenticated, "刷新Token无效")
	}

	// 根据用户ID重新生成新的 Access Token 和 Refresh Token
	userID := claims.UserID // 直接访问结构体字段
	username := claims.Username
	terminal := claims.Terminal

	newAccessToken, err := jwt.GenToken(int64(userID), username, 30*time.Minute,terminal)
	if err != nil {
		return &proto.RefreshTokenResponse{
			Success: false,
			Message: "生成新的Access Token失败",
		}, status.Error(codes.Internal, "内部错误")
	}

	newRefreshToken, err := jwt.GenToken(int64(userID), username, 7*24*time.Hour,terminal)
	if err != nil {
		return &proto.RefreshTokenResponse{
			Success: false,
			Message: "生成新的Refresh Token失败",
		}, status.Error(codes.Internal, "内部错误")
	}

	return &proto.RefreshTokenResponse{
		Success:      true,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// SendSmsCode 发送短信验证码
// SendSmsCode 发送短信验证码
func (u *UserSrv) SendSmsCode(ctx context.Context, req *proto.SendSmsCodeRequest) (*proto.SendSmsCodeResponse, error) {
    fmt.Println("in SendSmsCode ... ") // 打印进入方法的日志

    // 参数处理
    if req.GetPhone() == "" {
        // 无效的请求
        return nil, status.Error(codes.InvalidArgument, "请求参数有误") // 返回 gRPC 的 InvalidArgument 错误
    }

    // 业务处理
    response, err := user.SendSmsCode(ctx, req.Phone)
    if err != nil {
        zap.L().Error("user.SendSmsCode failed", zap.Error(err)) // 记录错误日志
        return response, status.Error(codes.Internal, "内部错误") // 返回 gRPC 的 Internal 错误
    }

    return response, nil // 返回响应
}

// LoginBySms 短信验证码登录
func (u *UserSrv) LoginBySms(ctx context.Context, req *proto.LoginBySmsRequest) (*proto.LoginResponse, error) {
    fmt.Println("in LoginBySms ... ") // 打印进入方法的日志

    // 参数处理
    if req.GetPhone() == "" || req.GetSmsCode() == "" {
        // 无效的请求
        return nil, status.Error(codes.InvalidArgument, "请求参数有误") // 返回 gRPC 的 InvalidArgument 错误
    }

    // 业务处理
    response, err := user.LoginBySms(ctx, req.Phone, req.SmsCode)
    if err != nil {
        zap.L().Error("user.LoginBySms failed", zap.Error(err)) // 记录错误日志
        return response, status.Error(codes.Internal, "内部错误") // 返回 gRPC 的 Internal 错误
    }

    return response, nil // 返回响应
}
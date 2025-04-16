package user

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"user_service/dao/mysql"
	"user_service/dao/redis"
	"user_service/proto"
	"user_service/third_party/jwt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Register 用户注册
func Register(ctx context.Context, username string, password string, email string, phone string) (*proto.RegisterResponse, error) {
	// 调用 mysql 包中的 RegisterUser 方法
	// 对密码进行加密
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return &proto.RegisterResponse{
			Success: false,
			Message: "密码加密失败",
		}, status.Error(codes.Internal, "内部错误")
	}
	err = mysql.RegisterUser(ctx, username, hashedPassword, email, phone)
	if err != nil {
		// 如果注册失败，返回错误
		return &proto.RegisterResponse{
			Success: false,
			Message: "注册用户失败",
		}, status.Error(codes.Internal, "内部错误")
	}

	// 注册成功，返回成功响应
	return &proto.RegisterResponse{
		Success: true,
		Message: "注册用户成功",
	}, nil
}

// Login 用户登录
func Login(ctx context.Context, username string, password string) (*proto.LoginResponse, error) {
	// 调用 mysql 包中的 GetUserByUsername 方法获取用户信息
	user, err := mysql.GetUserByUsername(ctx, username)
	if err != nil {
		// 如果用户不存在或获取用户信息失败，返回错误
		return &proto.LoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		}, status.Error(codes.Unauthenticated, "用户名或密码错误")
	}

	// 验证密码是否正确
	if err := CheckPassword(user.Password, password); err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		}, status.Error(codes.Unauthenticated, "用户名或密码错误")
	}
	// 获取终端类型
    terminal := getTerminalType(ctx)
	// 密码正确，生成 JWT Token
	tokenString, err := jwt.GenToken(user.ID, user.Username, 30*time.Minute,terminal)
	if err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "生成Token失败",
		}, status.Error(codes.Internal, "内部错误")
	}

	refreshToken, err := jwt.GenToken(user.ID, user.Username, 7*24*time.Hour,terminal)
	if err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "生成刷新Token失败",
		}, status.Error(codes.Internal, "内部错误")
	}

	// 返回登录成功响应
	return &proto.LoginResponse{
		Success:      true,
		Message:      "登录成功",
		Token:        tokenString,
		Refreshtoken: refreshToken,
	}, nil
}

// SendSmsCode 发送短信验证码
// SendSmsCode 发送短信验证码
func SendSmsCode(ctx context.Context, phone string) (*proto.SendSmsCodeResponse, error) {
	// 构造缓存键
	cacheKey := fmt.Sprintf("sms_code%s", phone)

	// 检查缓存中是否存在未过期的验证码
	err := redis.GetClient().Get(ctx, cacheKey).Err()
	if err == nil {
		// 验证码已存在且未过期，计算剩余时间
		ttl, _ := redis.GetClient().TTL(ctx, cacheKey).Result()
		remainingTime := int(ttl.Seconds())

		return &proto.SendSmsCodeResponse{
			Success: false,
			Message: fmt.Sprintf("验证码已发送，请在%d秒后重试", remainingTime),
		}, status.Error(codes.ResourceExhausted, "请求过于频繁")
	}

	// 生成验证码
	smsCode := "123456789" // 测试时可以使用固定验证码

	// 将验证码存储到缓存中，设置有效期（如5分钟）
	if err := redis.GetClient().Set(ctx, cacheKey, smsCode, 30*time.Minute).Err(); err != nil {
		return &proto.SendSmsCodeResponse{
			Success: false,
			Message: "验证码存储失败",
		}, err
	}

	// 调用短信服务发送验证码
	if err := sendSms(phone, smsCode); err != nil {
		return &proto.SendSmsCodeResponse{
			Success: false,
			Message: "短信发送失败",
		}, err
	}

	return &proto.SendSmsCodeResponse{
		Success: true,
		Message: "验证码发送成功",
	}, nil
}

// generateSmsCode 生成短信验证码
func generateSmsCode() string {
	// 生成6位随机数字验证码
	code := rand.NewSource(time.Now().UnixNano())
	return fmt.Sprintf("%06d", code)
}

// sendSms 发送短信
func sendSms(phone, code string) error {
	// 调用第三方短信服务发送验证码
	// 这里只是一个示例，实际实现需要根据具体的短信服务提供商进行调整
	return nil
}

// LoginBySms 短信验证码登录
func LoginBySms(ctx context.Context, phone, smsCode string) (*proto.LoginResponse, error) {
	// 验证短信验证码
	if err := verifySmsCode(ctx, phone, smsCode); err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "验证码错误或已过期",
		}, err
	}

	// 根据手机号获取用户信息
	user, err := mysql.GetUserByPhone(ctx, phone)
	if err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "用户不存在",
		}, err
	}

	// 获取终端类型
    terminal := getTerminalType(ctx)

	// 生成 JWT Token 并返回
	accessToken, err := jwt.GenToken(user.ID, user.Username, 30*time.Minute,terminal)
	if err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "生成Token失败",
		}, err
	}

	refreshToken, err := jwt.GenToken(user.ID, user.Username, 7*24*time.Hour,terminal)
	if err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "生成刷新Token失败",
		}, err
	}

	return &proto.LoginResponse{
		Success:      true,
		Message:      "登录成功",
		Token:        accessToken,
		Refreshtoken: refreshToken,
	}, nil
}

// verifySmsCode 验证短信验证码
func verifySmsCode(ctx context.Context, phone, smsCode string) error {
	// 从缓存中获取存储的验证码
	cacheKey := fmt.Sprintf("sms_code%s", phone)
	storedCode, err := redis.GetClient().Get(ctx, cacheKey).Result()
	if err != nil {
		return fmt.Errorf("验证码不存在或已过期")
	}

	// 验证用户输入的验证码是否正确
	if storedCode != smsCode {
		return fmt.Errorf("验证码错误")
	}

	// 验证通过后，删除缓存中的验证码，防止重复使用

	if err := redis.GetClient().Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("删除验证码失败")
	}

	return nil
}

// getTerminalType 获取终端类型
func getTerminalType(ctx context.Context) string {
    // 从上下文中获取终端类型（示例：通过元数据传递）
    md, _ := metadata.FromIncomingContext(ctx)
    if terminal, ok := md["terminal"]; ok && len(terminal) > 0 {
        return terminal[0]
    }
    return "unknown"
}
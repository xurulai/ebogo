package mysql

import (
	"context"
	"user_service/errno"
	"user_service/model"
	"user_service/third_party/snowflake"
)

// RegisterUser 注册用户
func RegisterUser(ctx context.Context, username, password, email,phone string) error {
	// 使用 GORM 的 WithContext 方法确保操作在指定的上下文中执行。
	// 查询用户是否已存在
	var user model.User
	result := db.WithContext(ctx).
		Where("username = ?", username).
		First(&user)

	// 如果用户已存在，返回错误
	if result.Error == nil {
		return errno.ErrUserAlreadyExists
	}

	// 如果用户不存在，创建新用户
	userId := snowflake.GenID()
	newUser := model.User{
		ID:       userId,
		Username: username,
		Password: password, // 实际项目中应先对密码进行加密处理
		Email:    email,
		Phone:    phone,	
	}

	result = db.WithContext(ctx).Create(&newUser)

	// 如果创建失败，返回错误
	if result.Error != nil {
		return errno.ErrCreateUserFailed
	}

	return nil // 操作成功，返回 nil
}

// GetUserByUsername 根据用户名获取用户信息
func GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, errno.ErrUserNotExisted
	}
	return &user, nil
}

// GetUserByPhone 根据手机号获取用户信息
func GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
    var user model.User
    result := db.WithContext(ctx).Where("phone = ?", phone).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
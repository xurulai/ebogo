package errno

import "errors"

var (
	ErrQueryFailed       = errors.New("query db failed")
	ErrQueryEmpty        = errors.New("query empty") // 查询结果为空
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCreateUserFailed  = errors.New("create user failed")
	ErrUserNotExisted	 = errors.New("user not existed")
)

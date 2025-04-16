package jwt

import (
	"errors"
	"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
)

// 定义过期时间
const TokenExpireDuration = time.Hour * 24

var mySecret = []byte("肥嘟嘟左卫门")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
	Terminal string `json:"terminal"`
}

// GenToken 生成JWT
func GenToken(userID int64, username string, expireTime time.Duration,terminal string) (string, error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireTime).Unix(), // 使用传入的过期时间
			Issuer:    "bluebell",                        // 签发人
			
		},
		Terminal: terminal,   //添加终端信息
		
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 创建一个新的 MyClaims 结构体实例，用于存放解析后的 JWT 声明信息
	var mc = new(MyClaims)

	// 使用 jwt.ParseWithClaims 函数解析 tokenString，并将解析结果存入 mc
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		// 返回用于校验的密钥 mySecret
		return mySecret, nil
	})

	// 如果在解析过程中出现错误，返回 nil 和错误信息
	if err != nil {
		return nil, err
	}

	// 校验 token 的有效性
	if token.Valid { // token.Valid 是一个布尔值，表示 token 是否有效
		// 如果有效，返回 MyClaims 实例和 nil 错误
		return mc, nil
	}

	// 如果 token 无效，返回 nil 和一个新的错误信息 "invalid token"
	return nil, errors.New("invalid token")
}

package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MyClaims 自定义声明结构体并内嵌 jwt.RegisteredClaims
// jwt包自带的 jwt.RegisteredClaims 只包含 RFC 7519 标准字段
// 我们这里需要额外记录 username、user_id，所以自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

// JWT 封装令牌的签发与解析，secret 与过期时间通过构造函数注入
type JWT struct {
	secret []byte
	expire time.Duration
}

// NewJWT 创建 JWT 实例，secret 为签名密钥，expireHours 为令牌有效期（小时）
func NewJWT(secret string, expireHours int) *JWT {
	return &JWT{
		secret: []byte(secret),
		expire: time.Duration(expireHours) * time.Hour,
	}
}

// GenToken 生成JWT
func (j *JWT) GenToken(userID int64, username string) (string, error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expire)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // 签发时间
			Issuer:    "li",                                         // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(j.secret)
}

// ParseToken 解析 JWT
func (j *JWT) ParseToken(tokenString string) (*MyClaims, error) {
	// 解析Token
	var mc = new(MyClaims)
	_, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (any, error) {
		// 仅允许 HS256，防止算法混淆攻击（algorithm confusion）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	// token 合法时 mc 已被填充；校验失败会在上面返回 err
	return mc, nil
}

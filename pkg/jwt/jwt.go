package jwt

import (
	"encoding/json"
	"time"

	"moocss.com/tiga/pkg/base64"
	"moocss.com/tiga/pkg/conf"
	"moocss.com/tiga/pkg/log"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var (
	// token过期
	TokenExpired error =  errors.New("Token is expired")

	// token未激活
	TokenNotValidYet error = errors.New("Token not active yet")

	// 不是token
	TokenMalformed error = errors.New("That's not even a token")

	// 无法处理token
	TokenInvalid error = errors.New("Couldn't handle this token:")

	// token生成错误
	TokenFailure error = errors.New("Token generate failure")

	_ JWT = (*jwt)(nil)
)

// 定义载荷
// Structured version of Claims Section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
// See examples for how to use this with your own claim types
type CustomClaims struct {
	Data json.RawMessage // 自由数据
	jwtgo.StandardClaims
}

// JWT is a jwt interface.
type JWT interface {
	// 创建Token
	CreateToken(claims CustomClaims) (string, error)
	// 解析Token
	ParseToken(tokenString string) (*CustomClaims, error)
	// 更新Token
	RefreshToken(tokenString string) (string, error)

	// 验证Token
	TokenValid(tokenString string) error
	// 验证Token是否过期
	TokenExpires(tokenString string) (bool, error)

	GetSigningKey() string
}

// jwt基本数据结构
type jwt struct {
	opts *options
	log   *log.Helper
}

// New new a jwt with options.
func New(opts ...Option) JWT {
	options := DefaultOptions()
	for _, o := range opts {
		o(options)
	}

	return &jwt{
		opts: options,
		log:  log.NewHelper("jwt", options.logger),
	}
}

// 获取base64编码过的 signing key
func (j *jwt) GetSigningKey() string {
	return base64.Base64UrlEncode(j.opts.signingKey)
}

// 创建Token(自定义Data, 例如: 基于用户信息claims)
// 使用根据配置选择算法进行token生成(如果配置文件没有配置, 默认是HS256)
// 使用用户基本信息claims以及签名key(signkey)生成token
func (j *jwt) CreateToken(claims CustomClaims) (string, error) {
	return j.createTokenString(claims)
}

func (j *jwt) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := j.verifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 将token中的claims信息解析出来和用户原始数据进行校验
	// 做以下类型断言，将token.Claims转换成具体用户自定义的Claims结构体
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// 控制台打印信息
		j.log.Infof("claims: %v", claims)

		return claims, nil
	}
	return nil, TokenInvalid
}

func (j *jwt) RefreshToken(tokenString string) (string, error) {
	// TimeFunc为一个默认值是time.Now的当前时间变量,用来解析token后进行过期时间验证
	// 可以使用其他的时间值来覆盖
	jwtgo.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	token, err := j.verifyToken(tokenString)
	if err != nil {
		return "", err
	}

	// 校验token当前还有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwtgo.TimeFunc = time.Now

		// 生成新的token
		return j.createTokenString(*claims)
	}
	return "", TokenInvalid
}

// TokenValid 验证token
func (j *jwt) TokenValid(tokenString string) error {
	token, err := j.verifyToken(tokenString)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return nil
	}
	return TokenInvalid
}

// TokenExpires 验证Token是否过期, 三种状态：未过期、过期、过期但可以Refresh
func (j *jwt) TokenExpires(tokenString string) (bool, error) {
	token, err := j.verifyToken(tokenString)
	if err != nil {
		return false, err
	}

	// 将token中的claims信息解析出来和用户原始数据进行校验
	// 做以下类型断言，将token.Claims转换成具体用户自定义的Claims结构体
	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		// ExpiresAt
		if claims.ExpiresAt > time.Now().Unix() {
			// pp.Printf("ExpiresAt:%v\n", claims.ExpiresAt)
			return false, nil
		}

		// IssuedAt
		if claims.IssuedAt > time.Now().Unix() {
			return true, TokenExpired
		}

		return false, nil
	}

	return false, nil
}

func (j *jwt) createTokenString(claims CustomClaims) (string, error) {
	var method jwtgo.SigningMethod
	signingMethod := conf.Get("app.jwt.signing_method")

	switch signingMethod {
	case "HS256":
		method = jwtgo.SigningMethodHS256
	case "HS384":
		method = jwtgo.SigningMethodHS384
	case "HS512":
		method = jwtgo.SigningMethodHS512
	default:
		method = jwtgo.SigningMethodHS256
	}

	// 修改Claims的过期时间(int64)
	// https://gowalker.org/github.com/dgrijalva/jwt-go#StandardClaims
	claims.StandardClaims.ExpiresAt = j.opts.expiresTime
	// 返回一个token的结构体指针
	token := jwtgo.NewWithClaims(method, claims)
	tokenString, err := token.SignedString(j.opts.signingKey)
	if err != nil {
		j.log.Infof("生成Token错误: %v", err)
		return "", TokenFailure
	}
	return tokenString, nil
}

func (j *jwt) verifyToken(tokenString string) (*jwtgo.Token, error) {
	if len(tokenString) == 0 {
		return nil, TokenMalformed
	}

	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// 输入用户自定义的Claims结构体对象,token,以及自定义函数来解析token字符串为jwt的Token结构体指针
	// Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwtgo.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return j.opts.signingKey, nil
	})
	if err != nil {
		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		// jwt.ValidationError 是一个无效token的错误结构
		if ve, ok := err.(*jwtgo.ValidationError); ok {
			// ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwtgo.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwtgo.ValidationErrorExpired != 0 { // ValidationErrorExpired表示Token过期
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwtgo.ValidationErrorNotValidYet != 0 { // ValidationErrorNotValidYet表示无效token
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	return token, nil
}

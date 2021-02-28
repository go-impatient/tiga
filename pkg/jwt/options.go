package jwt

import (
	"time"

	"moocss.com/tiga/pkg/base64"
	"moocss.com/tiga/pkg/conf"
	"moocss.com/tiga/pkg/log"
)

var (
	SignKey     string = "Z3d0X3NpZ25fa2V5" // gwt_sign_key
	ExpiresTime int64  = 172800             // 2h
)

// Option is jwt option
type Option func(*options)

type options struct {
	signingKey  []byte // 签名 signkey
	expiresTime int64  // 过期时间

	logger log.Logger
}

// DefaultOptions .
func DefaultOptions() *options {
	signingKey := getSignKey()
	expiresTime := getExpiresTime()
	return &options{
		signingKey:  signingKey,
		expiresTime: expiresTime,
		logger:      log.DefaultLogger,
	}
}

// WithSigningKey .
func WithSigningKey(key []byte) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

// WithExpiresTime .
func WithExpiresTime(t int64) Option {
	return func(o *options) {
		o.expiresTime = t
	}
}

// WithLogger with config loogger.
func WithLogger(l log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

// getSignKey 获取signing key
func getSignKey() []byte {
	signingKey := conf.Get("app.jwt.secret")
	if len(signingKey) == 0 {
		signingKey = signingKey // 默认signing key
	}

	// 恢复到原始的 signing key 值, 例如: gwt_sign_key
	s, err := base64.Base64UrlDecode(signingKey)
	if err != nil {
		return nil
	}
	return s
}

// getExpiresTime 获取过期时间
func getExpiresTime() int64 {
	jTime := conf.GetInt64("app.jwt.expire")
	if jTime == 0 {
		jTime = ExpiresTime
	}

	return time.Now().Add(time.Duration(jTime) * time.Second).Unix()
}

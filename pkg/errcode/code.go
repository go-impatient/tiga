package errcode

var (
	// Common errors
	ServerSuccess             = &Error{Code: 0, Message: "OK"}
	ServerError               = &Error{Code: 10000000, Message: "服务内部错误"}
	InvalidParams             = &Error{Code: 10000001, Message: "入参错误"}
	ServerNotFound            = &Error{Code: 10000002, Message: "找不到"}
	UnauthorizedAuthNotExist  = &Error{Code: 10000003, Message: "鉴权失败，找不到对应的AppKey和AppSecret"}
	UnauthorizedTokenError    = &Error{Code: 10000004, Message: "鉴权失败，Token错误"}
	UnauthorizedTokenTimeout  = &Error{Code: 10000005, Message: "鉴权失败，Token超时"}
	UnauthorizedTokenGenerate = &Error{Code: 10000006, Message: "鉴权失败，Token生成失败"}
	TooManyRequests           = &Error{Code: 10000007, Message: "请求过多"}
	BindError                 = &Error{Code: 10000008, Message: "将请求体绑定到结构体时发生错误"}
	ParamsValidateError       = &Error{Code: 10000009, Message: "参数验证失败"} // 普通验证

	ValidationError = &Error{Code: 20000001, Message: "验证失败"} // 验证插件验证
	DatabaseError   = &Error{Code: 20000002, Message: "数据库错误"}

	// Module errors
	EncryptError           = &Error{Code: 20010001, Message: "加密用户密码时发生错误"}
	UserNotFoundError      = &Error{Code: 20010002, Message: "查找用户失败"}
	TokenInvalidError      = &Error{Code: 20010003, Message: "令牌无效"}
	PasswordIncorrectError = &Error{Code: 20010004, Message: "密码不正确"}
)

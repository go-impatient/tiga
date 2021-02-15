package errcode

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Id      string      `json:"id"`              // 唯一标识
	Code    int32       `json:"code"`            // 状态码
	Data    interface{} `json:"data,omitempty"`  // 数据
	Message string      `json:"message"`         // 信息
	Extra   interface{} `json:"extra,omitempty"` // 扩展
}

type Option func(*Error)

// New, 生成自定义错误
func New(id string, code int32, opts ...Option) *Error {
	if len(id) <= 0 {
		id = "error:"
	}
	// default
	entity := &Error{
		Id:      id,
		Code:    code,
		Message: "",
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(entity)
	}
	return entity
}

func Msg(msg string) Option {
	return func(e *Error) {
		e.Message = msg
	}
}

func Data(data interface{}) Option {
	return func(e *Error) {
		e.Data = data
	}
}

func Extra(extra interface{}) Option {
	return func(e *Error) {
		e.Extra = extra
	}
}

func (e *Error) Err() error {
	return e
}

func (e *Error) Errorf() string {
	if e.Data == nil && e.Extra == nil {
		return fmt.Sprintf("go-error: code = %d ,message = %s", e.Code, e.Message)
	} else if e.Extra == nil {
		return fmt.Sprintf("go-error: code = %d ,message = %s ,data = %s", e.Code, e.Message, e.Data)
	} else if e.Data == nil {
		return fmt.Sprintf("go-error: code = %d ,message = %s ,extra = %s", e.Code, e.Message, e.Extra)
	} else {
		return fmt.Sprintf("go-error: code = %d ,message = %s ,data = %s ,extra = %s", e.Code, e.Message, e.Data, e.Extra)
	}
}

func (e *Error) Coded() int32 {
	return e.Code
}

func (e *Error) Msg() string {
	return e.Message
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.Message, args...)
}

func (e *Error) Error() string {
	// []byte类型，转化成string类型便于查看
	b, _ := json.Marshal(e)
	return string(b)
}

// StatusCode, 标准的HTTP状态码
func (e *Error) StatusCode() int {
	switch e.Coded() {
	case ServerSuccess.Coded():
		return http.StatusOK
	case ServerError.Coded():
		return http.StatusInternalServerError
	case InvalidParams.Coded():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Coded():
		fallthrough
	case UnauthorizedTokenError.Coded():
		fallthrough
	case UnauthorizedTokenGenerate.Coded():
		fallthrough
	case UnauthorizedTokenTimeout.Coded():
		return http.StatusUnauthorized
	case TooManyRequests.Coded():
		return http.StatusTooManyRequests
	}

	return http.StatusInternalServerError
}

// Parse 试图将一个 JSON 字符串解析为一个Error。
// 如果失败，它将把给定的字符串设置为错误细节。
func Parse(err string) *Error {
	e := new(Error)
	verr := json.Unmarshal([]byte(err), e)
	if verr != nil {
		e.Message = err
	}
	return e
}

// BadRequest generates a 400 error.
func BadRequest(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf(format, a...),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    http.StatusUnauthorized,
		Message: fmt.Sprintf(format, a...),
	}
}

// Forbidden generates a 403 error.
func Forbidden(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    http.StatusForbidden,
		Message: fmt.Sprintf(format, a...),
	}
}

// NotFound generates a 404 error.
func NotFound(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf(format, a...),
	}
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    http.StatusMethodNotAllowed,
		Message: fmt.Sprintf(format, a...),
	}
}

// Timeout generates a 408 error.
func Timeout(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    http.StatusRequestTimeout,
		Message: fmt.Sprintf(format, a...),
	}
}

// Conflict generates a 409 error.
func Conflict(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    http.StatusConflict,
		Message: fmt.Sprintf(format, a...),
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf(format, a...),
	}
}

// NotImplemented generates a 501 error
func NotImplemented(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    501,
		Message: fmt.Sprintf(format, a...),
	}
}

// BadGateway generates a 502 error
func BadGateway(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    502,
		Message: fmt.Sprintf(format, a...),
	}
}

// AppUnavailable generates a 503 error
func AppUnavailable(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    503,
		Message: fmt.Sprintf(format, a...),
	}
}

// GatewayTimeout generates a 504 error
func GatewayTimeout(id, format string, a ...interface{}) error {
	return &Error{
		Id:      id,
		Code:    504,
		Message: fmt.Sprintf(format, a...),
	}
}

// Equal tries to compare errorsx
func Equal(err1 error, err2 error) bool {
	verr1, ok1 := err1.(*Error)
	verr2, ok2 := err2.(*Error)

	if ok1 != ok2 {
		return false
	}

	if !ok1 {
		return err1 == err2
	}

	if verr1.Code != verr2.Code {
		return false
	}

	return true
}

// FromError try to convert go error to *Error
func FromError(err error) *Error {
	verr, ok := err.(*Error)
	if ok && verr != nil {
		return verr
	}

	return Parse(err.Error())
}

// Wrap wraps errorsx
func Wrap(err error, msg string) error {
	return fmt.Errorf(`%s: %s`, msg, err.Error())
}

func Wrapf(err error, format string, args ...interface{}) error {
	return Wrap(err, fmt.Sprintf(format, args...))
}

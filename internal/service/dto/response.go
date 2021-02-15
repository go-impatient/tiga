package dto

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"moocss.com/tiga/pkg/errcode"
)

// Data is the parent model for returning data in this api,
// includes meta for pagination
type Data struct {
	Result     interface{}     `json:"result"`
	Pagination *PaginationData `json:"pagination,omitempty"`
}

// PaginationData ...
type PaginationData struct {
	Page  int `json:"page"`  // 页码, offset (PageSize *(PageNumber -1)
	Limit int `json:"limit"` // 每页数量, PageSize
	Total int `json:"total"` // 总行数
}

// ResponseListBody ... 返回分页列表数据
type ResponseListBody struct {
	Code    int              `json:"code"`
	Data    map[string]*Data `json:"data"`
	Message string           `json:"message"`
	Extra   interface{}      `json:"extra,omitempty"` // 扩展
}

// ResponseBody ...
type ResponseBody struct {
	ErrorResponseBody
	SuccessResponseBody
}

// SuccessResponseBody 定义了标准的 API 接口成功时返回数据模型
type SuccessResponseBody struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Extra   interface{} `json:"extra,omitempty"` // 扩展
}

// ErrorResponseBody 定义了标准的 API 接口错误时返回数据模型
type ErrorResponseBody struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

// BoolResponseBody ...
type BoolResponseBody struct {
	SuccessResponseBody
	Data bool `json:"data"`
}

// StringResponseBody ...
type StringResponseBody struct {
	SuccessResponseBody
	Data string `json:"data"`
}

func HandleSuccess(c *gin.Context, err error) {
	merr := errcode.FromError(err)
	result := errcode.New("gaia.api.error", merr.Code)
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func HandleSuccessData(c *gin.Context, err error, data interface{}) {
	merr := errcode.FromError(err)
	result := errcode.New("gaia.api.error", merr.Code, errcode.Data(data))
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func HandleError(c *gin.Context, statusCode int, err error) {
	merr := errcode.FromError(err)
	c.AbortWithStatusJSON(statusCode, merr)
}

func HandleErrorF(c *gin.Context, statusCode int, err error) {
	debug.PrintStack()
	merr := errcode.FromError(err)
	c.AbortWithStatusJSON(statusCode, merr)
}

package hserver

import (
	"context"
	"fmt"
	"github.com/brian-god/go-pkg/constant"
	"github.com/brian-god/go-pkg/hserver/herrors"
	"github.com/brian-god/go-pkg/token"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"
	"net/http"
	"time"
)

type Router interface {
	ConfigRoutes(r *server.Hertz, t *token.Token) // 配置路由
}

type Option struct {
	RateQPS         int
	CryptoKey       string
	TokenIssuer     string
	TokenSigningKey string
	ReleaseMode     bool
}

// ResponseResult 响应结果
type ResponseResult struct {
	Data  interface{}
	Error *herrors.ServerError
}

func DefaultResponseResult() *ResponseResult {
	return &ResponseResult{}
}
func (r *ResponseResult) WithData(data interface{}) *ResponseResult {
	r.Data = data
	return r
}

func (r *ResponseResult) WithError(err *herrors.ServerError) *ResponseResult {
	r.Error = err
	return r
}

// ServiceFunc 实际提供服务的函数
type ServiceFunc[T any] func(ctx context.Context, par *T) *ResponseResult

// ServiceNotParFunc 实际提供服务的函数(无参数)
type ServiceNotParFunc func(ctx context.Context) *ResponseResult

// Handler 接口处理器
type Handler[T any] struct {
	Context        context.Context
	RequestContext *app.RequestContext
	Param          *T
	Error          error
}

// NewHandler [T any] ， handler 工厂函数
// 参数：
//
//	ctx ： desc
//	c ： desc
//
// 返回值：
//
//	*Handler[T] ：desc
func NewHandler[T any](ctx context.Context, c *app.RequestContext) *Handler[T] {
	return &Handler[T]{Context: ctx, RequestContext: c}
}

// NewHandlerFu [T any] ， handlerFun 工厂函数
// 参数：
//
//	fun ： desc
//
// 返回值：
//
//	app.HandlerFunc ：desc
func NewHandlerFu[T any](fun ServiceFunc[T]) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		NewHandler[T](c, ctx).WithBinder().Do(fun)
	}
}

// NewNotParHandlerFu [T any] ， 无参数的处理器
// 参数：
//
//	fun ： desc
//
// 返回值：
//
//	app.HandlerFunc ：desc
func NewNotParHandlerFu(fun ServiceNotParFunc) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		NewHandler[int](c, ctx).DoNotPar(fun)
	}
}

// WithBinder ， 绑定并验证参数
// 参数：
//
//	param ： desc
//
// 返回值：
//
//	*Handler[T] ：desc
func (h *Handler[T]) WithBinder() *Handler[T] {
	h.Param = new(T)

	if err := h.RequestContext.BindAndValidate(h.Param); err != nil {
		h.Error = err
		msg := hertzI18n.MustGetMessage(context.TODO(), "StatusInvalidParam")
		if msg == "" {
			msg = "Param err"
		}
		h.RequestContext.JSON(http.StatusOK, utils.H{
			constant.RespCode:      constant.StatusInvalidParam,
			constant.RespMsg:       fmt.Sprintf("%s，fail err ：%+v", msg, h.Error),
			constant.RespReason:    "INVALID_PARAM",
			constant.RespTimestamp: time.Now().Unix(),
		})
		h.RequestContext.Abort()
	}
	return h
}

// Do ， 执行 server 函数
// 参数：
//
//	serviceFunc ： desc
//
// 返回值：
func (h *Handler[T]) Do(serviceFunc ServiceFunc[T]) {
	// 错误处理
	if h.Error != nil {
		return
	}
	if h.RequestContext.IsAborted() {
		return
	}
	// 调用服务函数
	res := serviceFunc(h.Context, h.Param)
	if res.Error != nil {
		ResponseFailureErr(h.RequestContext, res.Error)
	} else {
		ResponseSuccess(h.RequestContext, res.Data)
	}
}

// DoNotPar ， 无参数服务函数
// 参数：
//
//	serviceFunc ： desc
//
// 返回值：
func (h *Handler[T]) DoNotPar(serviceFunc ServiceNotParFunc) {
	// 错误处理
	if h.Error != nil {
		return
	}
	if h.RequestContext.IsAborted() {
		return
	}
	// 调用服务函数
	res := serviceFunc(h.Context)
	if res.Error != nil {
		ResponseFailureErr(h.RequestContext, res.Error)
	} else {
		ResponseSuccess(h.RequestContext, res.Data)
	}
}

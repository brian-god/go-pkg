package hserver

import (
	"github.com/brian-god/go-pkg/constant"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzI18n "github.com/hertz-contrib/i18n"
	"net/http"
	"time"
)

/*
ResponseSuccess 返回成功响应数据
调用示例:
1.成功时，不需要返回数据: server.ResponseSuccess(c, nil)
2.成功时，需要返回数据: server.ResponseSuccess(c, gin.H{"name": "xim","age": 18})
*/
func ResponseSuccess(c *app.RequestContext, data interface{}) {
	if data == nil {
		data = utils.H{}
	}
	message := hertzI18n.MustGetMessage(constant.ReasonSuccess)
	if message == "" {
		message = "ok"
	}
	c.JSON(http.StatusOK, utils.H{constant.RespCode: http.StatusOK, constant.RespMsg: message, constant.RespData: data, constant.RespReason: constant.ReasonSuccess})
}

/*
ResponseFailure 返回失败响应数据
调用示例:
*/
func ResponseFailure(c *app.RequestContext, code int, reason, msg string, data interface{}) {
	if data == nil {
		data = utils.H{}
	}
	if reason != "" {
		i18Mag := hertzI18n.MustGetMessage(reason)
		if i18Mag != "" {
			msg = i18Mag
		}
	}
	c.JSON(http.StatusOK, utils.H{constant.RespCode: code, constant.RespMsg: msg, constant.RespData: data, constant.RespReason: reason, constant.RespTimestamp: time.Now().Format("2006-01-02 15:04:05")})
}

/*
ResponseFailureErr 返回失败响应数据
调用示例:
*/
func ResponseFailureErr(c *app.RequestContext, err *herrors.ServerError) {
	code := err.Code
	if code == 0 {
		code = http.StatusInternalServerError
	}
	msg := err.DefMessage
	if err.Reason != "" {
		i18Mag := hertzI18n.MustGetMessage(err.Reason)
		if i18Mag != "" {
			msg = i18Mag
		}
	}
	errMsg := err.DefMessage
	if err.BusinessError != nil {
		errMsg = err.BusinessError.Error()
	}
	c.JSON(http.StatusOK, utils.H{constant.RespCode: code, constant.RespMsg: msg, constant.ErrMsg: errMsg, constant.RespReason: err.Reason, constant.RespTimestamp: time.Now().Format("2006-01-02 15:04:05")})
}

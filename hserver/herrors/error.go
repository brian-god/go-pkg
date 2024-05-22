package herrors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	ReasonStatusInternalServerError = "STATUS_INTERNAL_SERVER_ERROR"
	ReasonParameterError            = "PARAMETER_ERROR"
)

// ServerError 服务端错误定义
type ServerError struct {
	Code          int
	DefMessage    string
	Reason        string
	BusinessError error
}

func New(code int, reason string, msg string) *ServerError {
	return &ServerError{Code: code, Reason: reason, DefMessage: msg}
}

// DefaultError 默认错误
func DefaultError() *ServerError {
	return New(http.StatusInternalServerError, ReasonStatusInternalServerError, "Server Internal Error")
}

// NewErr 根据error创建错误
func NewErr(err error) *ServerError {
	return &ServerError{Code: http.StatusInternalServerError, Reason: ReasonStatusInternalServerError, DefMessage: err.Error(), BusinessError: err}
}

// NewParameterError 参数错误
func NewParameterError(err error) *ServerError {
	return &ServerError{Code: http.StatusBadRequest, Reason: ReasonParameterError, DefMessage: err.Error(), BusinessError: err}
}
func (r *ServerError) WithCode(code int) *ServerError {
	r.Code = code
	return r
}

func (r *ServerError) WithDefMsg(msg string) *ServerError {
	r.DefMessage = msg
	return r
}
func (r *ServerError) WithReason(reason string) *ServerError {
	r.Reason = reason
	return r
}
func (r *ServerError) WithBusinessError(err error) *ServerError {
	r.BusinessError = err
	return r
}
func (r *ServerError) Error() string {
	return fmt.Sprintf("code:%d,reason:%s,message:%s", r.Code, r.Reason, r.DefMessage)
}

func IsHError(err error) bool {
	var serverError *ServerError
	if errors.As(err, &serverError) {
		return true
	}
	return false
}

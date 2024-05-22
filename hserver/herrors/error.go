package herrors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	ReasonStatusInternalHError = "STATUS_INTERNAL_SERVER_ERROR"
	ReasonParameterError       = "PARAMETER_ERROR"
)

type Herr = *HError

// HError 服务端错误定义
type HError struct {
	Code          int
	DefMessage    string
	Reason        string
	BusinessError error
}

func New(code int, reason string, msg string) *HError {
	return &HError{Code: code, Reason: reason, DefMessage: msg}
}

// DefaultError 默认错误
func DefaultError() *HError {
	return New(http.StatusInternalServerError, ReasonStatusInternalHError, "Server Internal Error")
}

// NewErr 根据error创建错误
func NewErr(err error) *HError {
	return &HError{Code: http.StatusInternalServerError, Reason: ReasonStatusInternalHError, DefMessage: err.Error(), BusinessError: err}
}

// NewParameterError 参数错误
func NewParameterError(err error) *HError {
	return &HError{Code: http.StatusBadRequest, Reason: ReasonParameterError, DefMessage: err.Error(), BusinessError: err}
}
func (r *HError) WithCode(code int) *HError {
	r.Code = code
	return r
}

func (r *HError) WithDefMsg(msg string) *HError {
	r.DefMessage = msg
	return r
}
func (r *HError) WithReason(reason string) *HError {
	r.Reason = reason
	return r
}
func (r *HError) WithBusinessError(err error) *HError {
	r.BusinessError = err
	return r
}
func (r *HError) Error() string {
	return fmt.Sprintf("code:%d,reason:%s,message:%s", r.Code, r.Reason, r.DefMessage)
}

func IsHError(err error) bool {
	var HError *HError
	if errors.As(err, &HError) {
		return true
	}
	return false
}
